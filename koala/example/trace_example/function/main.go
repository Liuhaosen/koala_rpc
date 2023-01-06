package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
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
	if len(os.Args) != 2 {
		panic("缺少参数")
	}

	tracer, closer := Init("function-demo")
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	helloTo := os.Args[1]

	//开启span
	span := tracer.StartSpan("say-hello")
	span.SetTag("hello-to", helloTo)
	defer span.Finish()

	//将span放到context里, 在进程内传播.
	ctx := opentracing.ContextWithSpan(context.Background(), span)

	helloStr := formatString(ctx, helloTo)
	printHello(ctx, helloStr)
}

func formatString(ctx context.Context, helloTo string) string {
	//这里是基于say-hello的span新开启的span. 所以是say-hello的子span
	span, _ := opentracing.StartSpanFromContext(ctx, "formatString")
	defer span.Finish()

	helloStr := fmt.Sprintf("hello, %s !", helloTo)
	span.LogFields(
		log.String("event", "string-format"),
		log.String("value", helloStr),
	)
	//这里全都是内存操作, 所以非常快, 几乎可以忽略不计
	return helloStr
}

func printHello(ctx context.Context, helloStr string) {
	span, _ := opentracing.StartSpanFromContext(ctx, "printHello")
	defer span.Finish()
	//这里有println打印到终端, 属于文件操作. 所以比较耗时.
	println(helloStr)
	span.LogKV("event", "println")
}
