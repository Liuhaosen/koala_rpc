package main

import (
	"fmt"
	"io"
	"os"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	jaeger "github.com/uber/jaeger-client-go"
	config "github.com/uber/jaeger-client-go/config"
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
		panic("ERROR: Expectiong one argument")
	}

	//1. 初始化追踪实例, 传入当前服务的名字
	tracer, closer := Init("hello-world")
	defer closer.Close() //关闭实例

	//2. 拿到当前进程的命令行参数
	helloTo := os.Args[1]

	//3. 启动一个追踪, span名是say-hello
	span := tracer.StartSpan("say-hello") //日志记录开始时间
	span.SetTag("hello-to", helloTo)      //打上标签

	//4. 格式化字符串
	helloStr := fmt.Sprintf("hello, %s!", helloTo)

	//5. log记录两个值: event, value. 日志记录耗时和开始时间
	span.LogFields(
		log.String("event", "string-format"),
		log.String("value", helloStr),
	)

	println(helloStr)
	//6. 记录一个kv
	span.LogKV("event", "println")

	//7. 生命周期结束
	span.Finish() //结束时间
}
