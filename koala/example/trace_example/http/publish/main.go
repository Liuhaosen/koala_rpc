package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go/transport/zipkin"
)

//初始化
func Init(service string) (opentracing.Tracer, io.Closer) {
	//1. 初始化一个transport, 用来收集span
	transport, err := zipkin.NewHTTPTransport(
		"http://60.205.218.189:9411/api/v1/spans",
		zipkin.HTTPBatchSize(1),
		zipkin.HTTPLogger(jaeger.StdLogger),
	)
	if err != nil {
		panic(fmt.Sprintf("Error: cannot init zipkin : %v\n", err))
	}

	fmt.Printf("transport:%v\n", transport)

	//2. 初始化配置
	cfg := &config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
		},
	}

	//3. 初始化上报器, 把transport上报到zipkin
	r := jaeger.NewRemoteReporter(transport)
	fmt.Printf("r = %v\n", r)

	//4. 初始化追踪实例
	tracer, closer, err := cfg.New(
		service,
		config.Logger(jaeger.StdLogger),
		config.Reporter(r),
	)
	if err != nil {
		panic(fmt.Sprintf("Error: cannot init jaeger: %v\n", err))
	}
	return tracer, closer
}

func main() {
	tracer, closer := Init("http-publisher")
	defer closer.Close()

	http.HandleFunc("/publish", func(w http.ResponseWriter, r *http.Request) {
		spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		span := tracer.StartSpan("publish", ext.RPCServerOption(spanCtx))
		defer span.Finish()

		helloStr := r.FormValue("helloStr")
		println(helloStr)
	})

	fmt.Println(http.ListenAndServe(":8082", nil))
}
