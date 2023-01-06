package meta

import (
	"context"
	"modtest/gostudy/lesson2/ibinarytree/koala/registry"

	"google.golang.org/grpc"
)

//元信息结构.
type RpcMeta struct {
	CallerName     string           //调用方名字
	ServiceName    string           //服务方名字
	Method         string           //调用的方法
	CallerCluster  string           //调用方集群
	ServiceCluster string           //服务提供方集群
	TraceID        string           //trace_id
	Env            string           //环境
	CallerIDC      string           //调用方IDC(机房)
	ServiceIDC     string           //服务提供方IDC(机房)
	CurNode        *registry.Node   //当前节点
	HistoryNodes   []*registry.Node //历史选择节点
	AllNodes       []*registry.Node //服务提供方的节点列表
	Conn           *grpc.ClientConn //当前请求使用的连接
}

type rpcMetaContextKey struct{}

//初始化元信息, 存入context里
//参数: context, 服务名, 方法名
//如果要使用存入的元信息, 那么就用GetServerMeta来获取
func InitRpcMeta(ctx context.Context, service string, method string, caller string) context.Context {
	meta := &RpcMeta{
		Method:      method,
		ServiceName: service,
		CallerName:  caller,
	} //生成实例
	return context.WithValue(ctx, rpcMetaContextKey{}, meta) //设置到context里.
}

//从context里获取元信息
func GetRpcMeta(ctx context.Context) *RpcMeta {
	meta, ok := ctx.Value(rpcMetaContextKey{}).(*RpcMeta)
	if !ok {
		meta = &RpcMeta{}
	}
	return meta
}
