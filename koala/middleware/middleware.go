package middleware

import (
	"context"
	"fmt"
)

//中间件函数
type MiddlewareFunc func(ctx context.Context, req interface{}) (resp interface{}, err error)

//中间件
type Middleware func(MiddlewareFunc) MiddlewareFunc

//自定义中间件
// var userMiddleware []Middleware

//把中间件函数串成一个执行链.
//调用的时候会按顺序执行
func Chain(outer Middleware, others ...Middleware) Middleware {
	return func(next MiddlewareFunc) MiddlewareFunc {
		for i := len(others) - 1; i >= 0; i-- {
			//先把next当做参数传进去, 然后再执行第i个中间件.
			//得到自己的中间件处理函数. 然后下次再传进去
			// fmt.Printf("next before: %#v\n", next)
			next = others[i](next)
			//这里得到的是:
			//middleware2(next)
			//middleware1(middleware2)
			// fmt.Printf("next after: %#v\n", next)
		}
		//上边for 循环得到的是middleware2, middleware1
		//然后返回outer(middleware1)
		//那么这个outer执行的时候就会先执行outer的代码, 然后执行作为参数的中间件middleware1
		//然后由于middleware2在循环中作为next参数传到了middleware1里. 所以会继续执行middleware2. 直到没有函数可执行
		fmt.Println("outer exec")
		return outer(next)
	}
}

/*

这部分代码移到了server/server.go里. 由Server来决定使用哪些中间件. 而middleware包则专注于如何实现middleware

//对外暴露的Use接口. 用来自定义中间件
func Use(m ...Middleware) {
	userMiddleware = append(userMiddleware, m...)
}

//用户自定义中间件
func BuildServerMiddleware(handle MiddlewareFunc) (handleChain MiddlewareFunc) {
	var mids []Middleware
	mids = append(mids, PrometheusServerMiddleware)
	if len(userMiddleware) != 0 {
		mids = append(mids, userMiddleware...)
	}

	if len(mids) > 0 {
		m := Chain(mids[0], mids[1:]...)
		return m(handle)
	}
	return handle
}
*/
