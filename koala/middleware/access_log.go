package middleware

import (
	"context"
	"fmt"
	"modtest/gostudy/lesson2/ibinarytree/koala/logs"
	"modtest/gostudy/lesson2/ibinarytree/koala/meta"
	"time"

	"google.golang.org/grpc/status"
)

func AccessLogMiddleware(next MiddlewareFunc) MiddlewareFunc {
	return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
		ctx = logs.WithFieldContext(ctx)
		//TODO 使用调用方传递过来的traceId, 如果调用方没传过来, 再用下边这行代码生成
		ctx = logs.WithTraceId(ctx, logs.GenTraceId())

		startTime := time.Now()
		resp, err = next(ctx, req)
		if err != nil {
			logs.Error(ctx, "next middleware execute failed, err : %v", err)
		}

		serverMeta := meta.GetServerMeta(ctx)
		errStatus, _ := status.FromError(err)

		cost := time.Since(startTime).Nanoseconds() / 1000
		logs.AddField(ctx, "cost_us", cost)
		logs.AddField(ctx, "method", serverMeta.Method)

		logs.AddField(ctx, "cluster", serverMeta.Cluster)
		logs.AddField(ctx, "env", serverMeta.Env)
		logs.AddField(ctx, "server_ip", serverMeta.ServerIP)
		logs.AddField(ctx, "client_ip", serverMeta.ClientIP)
		logs.AddField(ctx, "idc", serverMeta.IDC)
		logs.Access(ctx, "result=%v", errStatus.Code()) //通过access把以上字段一起输出

		return
	}
}

func RpcLogMiddleware(next MiddlewareFunc) MiddlewareFunc {
	return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
		ctx = logs.WithFieldContext(ctx)
		startTime := time.Now()
		resp, err = next(ctx, req)
		rpcMeta := meta.GetRpcMeta(ctx)
		errStatus, _ := status.FromError(err)

		cost := time.Since(startTime).Nanoseconds() / 1000
		logs.AddField(ctx, "cost_us", cost)               //耗时,微秒
		logs.AddField(ctx, "method", rpcMeta.Method)      //方法
		logs.AddField(ctx, "server", rpcMeta.ServiceName) //服务名

		logs.AddField(ctx, "caller_cluster", rpcMeta.CallerCluster)    //服务调用方所在集群
		logs.AddField(ctx, "upstream_cluster", rpcMeta.ServiceCluster) //服务提供者集群
		logs.AddField(ctx, "rpc", 1)                                   //rpc调用
		logs.AddField(ctx, "env", rpcMeta.Env)                         //服务所处环境

		var upStreamInfo string
		for _, node := range rpcMeta.HistoryNodes {
			upStreamInfo += fmt.Sprintf(" %s:%d, ", node.IP, node.Port)
		}

		// if rpcMeta.CurNode != nil {
		// upStreamInfo = fmt.Sprintf("%s:%d", rpcMeta.CurNode.IP, rpcMeta.CurNode.Port)
		// for _, node := range rpcMeta.HistoryNodes {
		// 	upStreamInfo += fmt.Sprintf(", %s:%d", node.IP, node.Port)
		// }
		// }

		logs.AddField(ctx, "upstream", upStreamInfo)

		logs.AddField(ctx, "caller_idc", rpcMeta.CallerIDC)    //服务调用方数据中心
		logs.AddField(ctx, "upstream_idc", rpcMeta.ServiceIDC) //服务提供者数据中心
		logs.Access(ctx, "result=%v", errStatus.Code())        //打印到access日志
		return
	}
}
