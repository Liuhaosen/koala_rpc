package server

//koala服务入口
import (
	"context"
	"fmt"
	"log"
	"modtest/gostudy/lesson2/ibinarytree/koala/logs"
	"modtest/gostudy/lesson2/ibinarytree/koala/middleware"
	"modtest/gostudy/lesson2/ibinarytree/koala/registry"
	_ "modtest/gostudy/lesson2/ibinarytree/koala/registry/etcd"
	"modtest/gostudy/lesson2/ibinarytree/koala/util"
	"net"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go/transport/zipkin"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
)

//把服务放到这里来操作
type KoalaServer struct {
	*grpc.Server
	limiter        *rate.Limiter
	userMiddleware []middleware.Middleware
	register       registry.Registry
}

var koalaServer = &KoalaServer{
	Server: grpc.NewServer(),
}

func Use(m ...middleware.Middleware) {
	koalaServer.userMiddleware = append(koalaServer.userMiddleware, m...)
}

//初始化服务
func Init(serviceName string) (err error) {
	//1. 读取配置文件
	err = InitConfig(serviceName)
	if err != nil {
		fmt.Println("读取配置文件失败, 错误:", err)
		return
	}
	fmt.Println("初始服务配置成功")
	//2. 初始化限流器.
	//这里可以更改为任何一个限流器.
	//当前使用的是令牌桶
	if koalaConf.Limit.SwitchOn {
		koalaServer.limiter = rate.NewLimiter(rate.Limit(koalaConf.Limit.QPSLimit), koalaConf.Limit.QPSLimit)
	}
	fmt.Println("初始化限流器成功")
	//3.初始化日志库
	err = initLogger()
	if err != nil {
		fmt.Println("初始化日志库失败, 错误:", err)
		return
	}
	fmt.Println("初始化日志库成功")
	//4. 初始化注册中心
	err = initRegister(serviceName)
	if err != nil {
		fmt.Println("初始化注册中心失败, 错误:", err)
		logs.Error(context.TODO(), "init register failed, err : %v", err)
		return
	}
	fmt.Println("初始化注册中心成功")

	//初始化分布式追踪系统
	err = initTrace(serviceName)
	if err != nil {
		logs.Error(context.TODO(), "init tracing failed, err : %v\n ", err)
		fmt.Println("初始化分布式追踪系统失败, 错误:", err)
		return
	}
	fmt.Println("初始化分布式追踪系统成功")
	return
}

//初始化日志库
func initLogger() (err error) {
	/*
		//这里是文件日志库
		filename := fmt.Sprintf("%s/%s.log", koalaConf.Log.Dir, koalaConf.ServiceName)
		outputer, err := logs.NewFileOutputer(filename)
		if err != nil {
			return
		}
	*/

	//开启控制台日志库
	outputer := logs.NewConsoleOutputer()

	level := logs.GetLogLevel(koalaConf.Log.Level)
	logs.InitLogger(level, koalaConf.Log.ChanSize, koalaConf.ServiceName)
	logs.AddOutputer(outputer)
	return
}

//初始化追踪系统
func initTrace(serviceName string) (err error) {
	if !koalaConf.Trace.SwitchOn {
		return
	}

	transport, err := zipkin.NewHTTPTransport(
		koalaConf.Trace.ReportAddr,
		zipkin.HTTPBatchSize(16),
		zipkin.HTTPLogger(jaeger.StdLogger),
	)

	if err != nil {
		logs.Error(context.TODO(), "error: cannot init zipkin %v\n", err)
		fmt.Println("启动zipkin失败, 错误:", err)
		return
	}

	cfg := &config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  koalaConf.Trace.SampleType,
			Param: koalaConf.Trace.SampleRate,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
		},
	}

	r := jaeger.NewRemoteReporter(transport)

	tracer, closer, err := cfg.New(
		serviceName,
		config.Logger(jaeger.StdLogger), config.Reporter(r),
	)

	if err != nil {
		logs.Error(context.TODO(), "error : cannot init jaeger: %v\n", err)
		fmt.Println("初始化jaeger失败, 错误:", err)
		return
	}

	_ = closer
	opentracing.SetGlobalTracer(tracer)
	return
}

//运行服务
func Run() {
	/* if koalaConf.Prometheus.SwitchOn {
		go func() {
			http.Handle("/metrics", promhttp.Handler())
			addr := fmt.Sprintf("0.0.0.0:%d", koalaConf.Prometheus.Port)
			log.Fatal(http.ListenAndServe(addr, nil))
		}()
	} */

	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", koalaConf.Port))
	if err != nil {
		log.Fatal("failed to listen: ", err)
	}

	koalaServer.Serve(listener)
}

func GRPCServer() *grpc.Server {
	return koalaServer.Server
}

//创建服务的中间件
func BuildServerMiddleware(handle middleware.MiddlewareFunc) middleware.MiddlewareFunc {
	var mids []middleware.Middleware

	//加入Prometheus中间件
	if koalaConf.Prometheus.SwitchOn {
		mids = append(mids, middleware.PrometheusServerMiddleware)
	}

	//加入限流器中间件
	if koalaConf.Limit.SwitchOn {
		mids = append(mids, middleware.NewRateLimitMiddleware(koalaServer.limiter))
	}

	//加入访问日志中间件
	mids = append(mids, middleware.AccessLogMiddleware)

	//添加用户的中间件
	if len(koalaServer.userMiddleware) != 0 {
		mids = append(mids, koalaServer.userMiddleware...)
	}

	//添加分布式追踪中间件
	if koalaConf.Trace.SwitchOn {
		mids = append(mids, middleware.TraceServerMiddleware)
	}

	if len(koalaServer.userMiddleware) != 0 {
		mids = append(mids, koalaServer.userMiddleware...)
	}

	//把所有中间件组织成一个调用链
	// m := middleware.Chain(mids[0], mids[1:]...)
	m := middleware.Chain(middleware.PrepareMiddleware, mids...)
	//返回调用链的入口函数
	return m(handle)
}

//初始化注册中心
func initRegister(serviceName string) (err error) {
	if !koalaConf.Register.SwitchOn {
		return
	}

	ctx := context.TODO()
	//初始化插件

	// etcd.InitEtcd(ctx)//这里只要引入包, 然后使用_来证明引入该插件即可
	// return
	registryInst, err := registry.InitPlugin(ctx, koalaConf.Register.RegisterName,
		registry.WithAddrs([]string{koalaConf.Register.RegisterAddr}),
		registry.WithTimeout(koalaConf.Register.Timeout),
		registry.WithRegistryPath(koalaConf.Register.RegisterPath),
		registry.WithHeartBeat(koalaConf.Register.HeartBeat))
	if err != nil {
		fmt.Println("init plugin failed, err:", err)
		logs.Error(ctx, "init registry failed, err: %v", err)
		return
	}

	koalaServer.register = registryInst
	service := &registry.Service{
		Name: serviceName,
	}
	//获取本地ip
	ip, err := util.GetLocalIP()
	if err != nil {
		fmt.Println("get local ip failed, err:", err)
		return
	}

	service.Nodes = append(service.Nodes, &registry.Node{
		IP:   ip,
		Port: koalaConf.Port,
	})
	//把服务注册进去
	fmt.Printf("服务内容: %#v\n", service.Nodes[0])
	err = registryInst.Register(context.TODO(), service)
	if err != nil {
		fmt.Println("register service failed, err:", err)
	}
	return
}
