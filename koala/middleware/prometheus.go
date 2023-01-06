package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"modtest/gostudy/lesson2/ibinarytree/koala/logs"
	"modtest/gostudy/lesson2/ibinarytree/koala/meta"
	"modtest/gostudy/lesson2/ibinarytree/koala/middleware/prometheus"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	DefaultServerMetrics = prometheus.NewServerMetrics()
	DefaultRpcMetrics    = prometheus.NewRpcMetrics()
)

func init() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		addr := fmt.Sprintf("0.0.0.0:%d", 8888)
		err := http.ListenAndServe(addr, nil)
		if err != nil {
			logs.Error(context.TODO(), "监听失败, 错误: ", err)
			return
		}
	}()
}

//普罗米修斯服务中间件
func PrometheusServerMiddleware(next MiddlewareFunc) MiddlewareFunc {
	return func(ctx context.Context, req interface{}) (resp interface{}, err error) {

		serverMeta := meta.GetServerMeta(ctx)
		DefaultServerMetrics.IncrRequest(ctx, serverMeta.ServiceName, serverMeta.Method)
		startTime := time.Now()
		resp, err = next(ctx, req)
		DefaultServerMetrics.IncrCode(ctx, serverMeta.ServiceName, serverMeta.Method, err)
		DefaultServerMetrics.Latency(ctx, serverMeta.ServiceName, serverMeta.Method, time.Since(startTime).Nanoseconds()/1000)
		return
	}
}

//普罗米修斯rpc中间件
func PrometheusRpcMiddleware(next MiddlewareFunc) MiddlewareFunc {
	return func(ctx context.Context, req interface{}) (resp interface{}, err error) {

		rpcMeta := meta.GetRpcMeta(ctx)
		DefaultRpcMetrics.IncrRequest(ctx, rpcMeta.ServiceName, rpcMeta.Method)
		startTime := time.Now()
		resp, err = next(ctx, req)
		DefaultRpcMetrics.IncrCode(ctx, rpcMeta.ServiceName, rpcMeta.Method, err)
		DefaultRpcMetrics.Latency(ctx, rpcMeta.ServiceName, rpcMeta.Method, time.Since(startTime).Nanoseconds()/1000)
		return
	}
}
