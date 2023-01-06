package middleware

import (
	"context"
	"fmt"
	"modtest/gostudy/lesson2/ibinarytree/koala/errno"
	"modtest/gostudy/lesson2/ibinarytree/koala/loadbalance"
	"modtest/gostudy/lesson2/ibinarytree/koala/logs"
	"modtest/gostudy/lesson2/ibinarytree/koala/meta"
)

//负载均衡中间件
func NewLoadBalanceMiddleware(balancer loadbalance.LoadBalance) Middleware {
	return func(next MiddlewareFunc) MiddlewareFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			rpcMeta := meta.GetRpcMeta(ctx)
			//如果没有节点, 那就退出
			if len(rpcMeta.AllNodes) == 0 {
				// err = errors.New("没有找到实例")
				err = errno.NotHaveInstance
				logs.Error(ctx, "not have instance")
				return
			}

			//生成loadbalance的上下文, 用来过滤已选择的节点
			ctx = loadbalance.WithBalanceContext(ctx)
			for {
				//选择一台机器
				rpcMeta.CurNode, err = balancer.Select(ctx, rpcMeta.AllNodes)
				if err != nil {
					fmt.Println("选择机器失败, 错误:", err)
					return
				}
				//打印节点
				logs.Debug(ctx, "select node: %#v", rpcMeta.CurNode)

				//把当前节点放到历史选择节点里
				rpcMeta.HistoryNodes = append(rpcMeta.HistoryNodes, rpcMeta.CurNode)

				//执行下个中间件函数
				resp, err = next(ctx, req)
				if err != nil {
					//连接错误的话, 重试
					if errno.ISConnectError(err) {
						continue
					}
					return
				}
				break
			}
			return
		}
	}
}
