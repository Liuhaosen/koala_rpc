package middleware

import (
	"context"
	"modtest/gostudy/lesson2/ibinarytree/koala/meta"

	"github.com/afex/hystrix-go/hystrix"
)

//熔断器中间件, 最好放在中间件的第一个.
func HystriMiddleware(next MiddlewareFunc) MiddlewareFunc {
	return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
		rpcMeta := meta.GetRpcMeta(ctx)
		hystrixErr := hystrix.Do(rpcMeta.ServiceName, func() (err error) {
			resp, err = next(ctx, req)
			return err
		}, nil)

		if hystrixErr != nil {
			return nil, hystrixErr
		}

		return resp, hystrixErr
	}
}
