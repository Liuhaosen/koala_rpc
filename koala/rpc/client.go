package rpc

import (
	"context"
	"fmt"
	"modtest/gostudy/lesson2/ibinarytree/koala/loadbalance"
	"modtest/gostudy/lesson2/ibinarytree/koala/logs"
	"modtest/gostudy/lesson2/ibinarytree/koala/meta"
	"modtest/gostudy/lesson2/ibinarytree/koala/middleware"
	"modtest/gostudy/lesson2/ibinarytree/koala/registry"
	_ "modtest/gostudy/lesson2/ibinarytree/koala/registry/etcd"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

var initRegistryOnce sync.Once //多个goroutine执行, 只有一个会执行
var globalRegister registry.Registry

//客户端
type KoalaClient struct {
	opts     *RpcOptions
	register registry.Registry
	balance  loadbalance.LoadBalance
	limiter  *rate.Limiter
}

func NewKoalaClient(serviceName string, optfunc ...RpcOptionFunc) *KoalaClient {
	ctx := context.TODO()
	client := &KoalaClient{
		opts: &RpcOptions{
			ConnTimeout:       DefaultConnTimeout,
			WriteTimeout:      DefaultWriteTimeout,
			ReadTimeout:       DefaultReadTimeout,
			ServiceName:       serviceName,
			RegisterName:      "etcd",
			RegisterAddr:      "127.0.0.1:2379",
			RegisterPath:      "/ibinarytree/koala/service/",
			TraceReportAddr:   "http://60.205.218.189:9411/api/v1/spans",
			TraceSampleType:   "const",
			TraceSampleRate:   1,
			ClientServiceName: "default",
		},
		balance: loadbalance.NewRandomBalance(),
	}

	for _, opt := range optfunc {
		opt(client.opts)
	}

	//初始化日志库
	outputer := logs.NewConsoleOutputer()
	logs.InitLogger(logs.LogLevelDebug, 20000, client.opts.ServiceName)
	logs.AddOutputer(outputer)

	//初始化注册中心
	initRegistryOnce.Do(func() {
		var err error
		globalRegister, err = registry.InitPlugin(
			ctx,
			client.opts.RegisterName,
			registry.WithAddrs([]string{client.opts.RegisterAddr}),
			registry.WithTimeout(time.Second),
			registry.WithRegistryPath(client.opts.RegisterPath),
			registry.WithHeartBeat(10),
		)

		if err != nil {
			logs.Error(ctx, "init registry failed, err : %v", err)
			return
		}
	})

	//初始化限流,只有配置了才会生效.
	if client.opts.MaxLimitQps > 0 {
		client.limiter = rate.NewLimiter(rate.Limit(client.opts.MaxLimitQps), client.opts.MaxLimitQps)
	}

	//初始化分布式追踪
	middleware.InitTrace(client.opts.ClientServiceName, client.opts.TraceReportAddr, client.opts.TraceSampleType, client.opts.TraceSampleRate)

	client.register = globalRegister

	return client
}

//获取调用方名称
func (k *KoalaClient) getCaller(ctx context.Context) string {
	//? 这里有可能是rpcMeta
	// rpcMeta := meta.GetRpcMeta(ctx)
	serverMeta := meta.GetServerMeta(ctx)
	if serverMeta == nil {
		return ""
	}
	return serverMeta.ServiceName
}

//构造中间件
func (k *KoalaClient) buildMiddleware(handle middleware.MiddlewareFunc) middleware.MiddlewareFunc {
	var mids []middleware.Middleware

	mids = append(mids, middleware.PrepareMiddleware)       //trace_id
	mids = append(mids, middleware.RpcLogMiddleware)        //rpc日志
	mids = append(mids, middleware.TraceRpcMiddleware)      //分布式追踪
	mids = append(mids, middleware.PrometheusRpcMiddleware) //普罗米修斯监控

	if k.limiter != nil {
		mids = append(mids, middleware.NewRateLimitMiddleware(k.limiter))
	}

	mids = append(mids, middleware.HystriMiddleware)                    //熔断器
	mids = append(mids, middleware.NewDiscoveryMiddleware(k.register))  //服务发现
	mids = append(mids, middleware.NewLoadBalanceMiddleware(k.balance)) //负载均衡
	mids = append(mids, middleware.ShortConnectMiddleware)              //短链接
	m := middleware.Chain(mids[0], mids[1:]...)                         //按顺序串成中间件执行链
	return m(handle)
}

//调用方法的入口
func (k *KoalaClient) Call(ctx context.Context, method string, r interface{}, handle middleware.MiddlewareFunc) (resp interface{}, err error) {
	//获取调用方名称
	caller := k.getCaller(ctx)
	//完成meta数据初始化
	ctx = meta.InitRpcMeta(ctx, k.opts.ServiceName, method, caller)
	//封装中间件, 执行中间件函数
	//handle: rpc调用
	middlewareFunc := k.buildMiddleware(handle)
	resp, err = middlewareFunc(ctx, r)
	if err != nil {
		fmt.Println("middleware func exec failed, err: ", err)
		return nil, err
	}
	return resp, err
}
