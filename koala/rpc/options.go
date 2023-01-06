package rpc

import "time"

type RpcOptions struct {
	ConnTimeout       time.Duration //连接超时时间
	WriteTimeout      time.Duration //写入超时
	ReadTimeout       time.Duration //读取超时
	ServiceName       string        //服务名
	RegisterName      string        //注册中心名字
	RegisterAddr      string        //注册中心地址
	RegisterPath      string        //注册中心路径
	MaxLimitQps       int           //限流qps
	TraceReportAddr   string        //上报中心的地址
	TraceSampleType   string        //采样类型
	TraceSampleRate   float64       //采样频率
	ClientServiceName string        //客户端服务名
}

type RpcOptionFunc func(opts *RpcOptions)

//处理qps
func WithLimitQPS(qps int) RpcOptionFunc {
	return func(opts *RpcOptions) {
		opts.MaxLimitQps = qps
	}
}

//处理超时时间
func WithConnTimeout(timeout time.Duration) RpcOptionFunc {
	return func(opts *RpcOptions) {
		opts.ConnTimeout = timeout
	}
}

func WithRegistryPath(rpath string) RpcOptionFunc {
	return func(opts *RpcOptions) {
		opts.RegisterPath = rpath
	}
}

//处理追踪客户端服务名
func WithClientServiceName(name string) RpcOptionFunc {
	return func(opts *RpcOptions) {
		opts.ClientServiceName = name
	}
}

//处理采样速率
func WithTraceSampleRate(rate float64) RpcOptionFunc {
	return func(opts *RpcOptions) {
		opts.TraceSampleRate = rate
	}
}

//处理采样类型
func WithTraceSampleType(stype string) RpcOptionFunc {
	return func(opts *RpcOptions) {
		opts.TraceSampleType = stype
	}
}

//处理上报中心地址
func WithTraceReportAddr(addr string) RpcOptionFunc {
	return func(opts *RpcOptions) {
		opts.TraceReportAddr = addr
	}
}
