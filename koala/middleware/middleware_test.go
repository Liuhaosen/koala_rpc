package middleware

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestMiddleware(t *testing.T) {
	middleware1 := func(next MiddlewareFunc) MiddlewareFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			fmt.Printf("中间件1 开始\n")
			rand.Seed(time.Now().Unix()) //要有随机数种子
			num := rand.Intn(2)
			if num%2 == 0 {
				err = fmt.Errorf("请求终止.")
				return
			}
			resp, err = next(ctx, req)
			if err != nil {
				return
			}
			fmt.Printf("中间件1 结束 \n")
			return
		}
	}

	middleware2 := func(next MiddlewareFunc) MiddlewareFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			fmt.Printf("中间件2 开始\n")

			resp, err = next(ctx, req)
			if err != nil {
				return
			}
			fmt.Printf("中间件2 结束\n")
			return
		}
	}

	outer := func(next MiddlewareFunc) MiddlewareFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			fmt.Printf("外层函数 开始\n")
			resp, err = next(ctx, req)
			if err != nil {
				return
			}
			fmt.Printf("外层函数 结束\n")
			return
		}
	}

	proc := func(ctx context.Context, req interface{}) (resp interface{}, err error) {
		fmt.Println("处理函数 开始")
		fmt.Println("处理函数 结束")
		return
	}

	chain := Chain(outer, middleware1, middleware2) //把所有中间件串成一个链表, 那么按Chain的写法, 会在outer之前依次执行middleware1, middleware2
	// chain(proc)(context.Background(), "test")
	chainFunc := chain(proc)                             //执行这个处理函数, 返回又得到一个中间件函数
	resp, err := chainFunc(context.Background(), "test") //执行这个中间件函数, 就会把上边的处理函数和中间件函数全都执行一遍
	fmt.Printf("resp:%#v, err: %v", resp, err)
	//这里打印出来的resp是nil, 因为在Chain里, next最初只是形参. 传进去就是nil, 所以最后出来也是nil.
}
