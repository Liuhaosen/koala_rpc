package prometheus

import (
	"context"

	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/grpc/status"
)

/* //koala服务端采样打点
type Metrics struct {
	serverRequestCounter *prom.CounterVec //请求数量, 计数器
	serverCodeCounter    *prom.CounterVec //错误数量, 计数器
	serverLatencySummary *prom.SummaryVec //请求耗时,分位图
} */

//通用采样打点(服务端和客户端都能用)
type Metrics struct {
	requestCounter *prom.CounterVec //请求数量, 计数器
	codeCounter    *prom.CounterVec //错误数量, 计数器
	latencySummary *prom.SummaryVec //请求耗时,分位图
}

//生成server metrics实例
func NewRpcMetrics() *Metrics {
	return &Metrics{
		//1. 请求数量计数器
		requestCounter: promauto.NewCounterVec(
			prom.CounterOpts{
				Name: "koala_rpc_call_total",
				Help: "Total number of RPCs completed on the server, regardless of success or failure",
			},
			[]string{"service", "method"}, //两个标签, service和method, 这样就知道采样数据属于哪个服务哪个方法
		),
		//2. 错误码计数器
		codeCounter: promauto.NewCounterVec(
			prom.CounterOpts{
				Name: "koala_rpc_code_total",
				Help: "Total number of RPCs error code on the server",
			},
			[]string{"service", "method", "grpc_code"},
		),
		//3.请求耗时
		latencySummary: promauto.NewSummaryVec(
			prom.SummaryOpts{
				Name:       "koala_rpc_cost",
				Help:       "RPC latency distributions.",
				Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
			},
			[]string{"service", "method"},
		),
	}
}

//生成rpc metrics实例
func NewServerMetrics() *Metrics {
	return &Metrics{
		//1. 请求数量计数器
		requestCounter: promauto.NewCounterVec(
			prom.CounterOpts{
				Name: "koala_server_request_total",
				Help: "Total number of RPCs completed on the server, regardless of success or failure",
			},
			[]string{"service", "method"}, //两个标签, service和method, 这样就知道采样数据属于哪个服务哪个方法
		),
		//2. 错误码计数器
		codeCounter: promauto.NewCounterVec(
			prom.CounterOpts{
				Name: "koala_server_haldled_code_total",
				Help: "Total number of RPCs error code on the server",
			},
			[]string{"service", "method", "grpc_code"},
		),
		//3.请求耗时
		latencySummary: promauto.NewSummaryVec(
			prom.SummaryOpts{
				Name:       "koala_server_proc_cost",
				Help:       "RPC latency distributions.",
				Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
			},
			[]string{"service", "method"},
		),
	}
}

//暴露请求计数方法
func (m *Metrics) IncrRequest(ctx context.Context, serviceName string, methodName string) {
	m.requestCounter.WithLabelValues(serviceName, methodName).Inc()
}

//错误码计数方法
func (m *Metrics) IncrCode(ctx context.Context, serviceName string, methodName string, err error) {
	st, _ := status.FromError(err)
	m.codeCounter.WithLabelValues(serviceName, methodName, st.Code().String()).Inc()
}

//请求耗时计数
func (m *Metrics) Latency(ctx context.Context, serviceName string, methodName string, us int64) {
	m.latencySummary.WithLabelValues(serviceName, methodName).Observe(float64(us))
}
