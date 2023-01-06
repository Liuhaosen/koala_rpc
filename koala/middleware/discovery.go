package middleware

import (
	"context"
	"modtest/gostudy/lesson2/ibinarytree/koala/logs"
	"modtest/gostudy/lesson2/ibinarytree/koala/meta"
	"modtest/gostudy/lesson2/ibinarytree/koala/registry"
)

//服务发现中间件
//传入了registry.Registry实例.
func NewDiscoveryMiddleware(discovery registry.Registry) Middleware {
	return func(next MiddlewareFunc) MiddlewareFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			rpcMeta := meta.GetRpcMeta(ctx)
			if len(rpcMeta.AllNodes) > 0 {
				return next(ctx, req)
			}

			service, err := discovery.GetService(ctx, rpcMeta.ServiceName)
			if err != nil {
				logs.Error(ctx, "discovery service: %s failed, err : %v", err)
				return
			}
			//rpcMeta获得了服务列表
			rpcMeta.AllNodes = service.Nodes
			//执行下个中间件函数
			resp, err = next(ctx, req)
			return
		}
	}
}
