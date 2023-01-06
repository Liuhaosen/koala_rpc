package middleware

import (
	"context"
	"modtest/gostudy/lesson2/ibinarytree/koala/logs"
	"modtest/gostudy/lesson2/ibinarytree/koala/util"

	"google.golang.org/grpc/metadata"
)

func PrepareMiddleware(next MiddlewareFunc) MiddlewareFunc {
	return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
		//处理trace_id
		var traceId string
		//从ctx获取grpc的metadata
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			vals, ok := md[util.TraceID]
			if ok && len(vals) > 0 {
				traceId = vals[0]
			}
		}

		if len(traceId) == 0 {
			traceId = logs.GenTraceId()
		}

		ctx = logs.WithFieldContext(ctx)
		ctx = logs.WithTraceId(ctx, traceId)
		resp, err = next(ctx, req)
		return
	}
}
