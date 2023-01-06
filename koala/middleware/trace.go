package middleware

import (
	"context"
	"fmt"
	"modtest/gostudy/lesson2/ibinarytree/koala/logs"
	"modtest/gostudy/lesson2/ibinarytree/koala/meta"
	"modtest/gostudy/lesson2/ibinarytree/koala/util"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"google.golang.org/grpc/metadata"
)

func TraceServerMiddleware(next MiddlewareFunc) MiddlewareFunc {
	return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
		//1. 从ctx获取grpc的metadata
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			//没有的话新建一个
			md = metadata.Pairs()
		}

		//2. 拿到tracer实例
		tracer := opentracing.GlobalTracer()

		//3. 提取Span
		parentSpanContext, err := tracer.Extract(opentracing.HTTPHeaders, metadataTextMap(md))
		if err != nil && err != opentracing.ErrSpanContextNotFound {
			logs.Warn(ctx, "trace extract failed, parsing trace information: %v", err)
			fmt.Println("追踪extract失败, 解析追踪信息: ", err)
		}

		serverMeta := meta.GetServerMeta(ctx)
		//开始追踪该方法
		serverSpan := tracer.StartSpan(
			serverMeta.Method,
			ext.RPCServerOption(parentSpanContext),
		)

		serverSpan.SetTag("trace_id", logs.GetTraceId(ctx))
		ctx = opentracing.ContextWithSpan(ctx, serverSpan)
		resp, err = next(ctx, req)
		//记录错误
		if err != nil {
			ext.Error.Set(serverSpan, true)
			serverSpan.LogFields(log.String("event", "error"), log.String("message", err.Error()))
		}

		serverSpan.Finish()
		return
	}
}

func TraceRpcMiddleware(next MiddlewareFunc) MiddlewareFunc {
	return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
		tracer := opentracing.GlobalTracer()
		var parentSpanContext opentracing.SpanContext

		if parent := opentracing.SpanFromContext(ctx); parent != nil {
			parentSpanContext = parent.Context()
		}

		opts := []opentracing.StartSpanOption{
			opentracing.ChildOf(parentSpanContext),
			ext.SpanKindRPCClient,
			opentracing.Tag{Key: string(ext.Component), Value: "koala_rpc"},
			opentracing.Tag{Key: util.TraceID, Value: logs.GetTraceId(ctx)},
		}

		rpcMeta := meta.GetRpcMeta(ctx)
		clientSpan := tracer.StartSpan(rpcMeta.ServiceName, opts...)

		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.Pairs()
		}

		//通过http头部的方式注入到md里
		if err := tracer.Inject(clientSpan.Context(), opentracing.HTTPHeaders, metadataTextMap(md)); err != nil {
			logs.Debug(ctx, "grpc_opentracing: failed serializing trace information: %v", err)
		}

		ctx = metadata.NewOutgoingContext(ctx, md)
		ctx = metadata.AppendToOutgoingContext(ctx, util.TraceID, logs.GetTraceId(ctx))
		ctx = opentracing.ContextWithSpan(ctx, clientSpan)
		resp, err = next(ctx, req)
		if err != nil {
			ext.Error.Set(clientSpan, true)
			clientSpan.LogFields(log.String("event", "error"), log.String("message", err.Error()))
		}

		clientSpan.Finish()
		return
	}
}
