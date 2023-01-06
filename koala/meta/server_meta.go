package meta

import "context"

//元信息结构.
type ServerMeta struct {
	ServiceName string
	Method      string
	Cluster     string
	TraceID     string
	Env         string
	ServerIP    string
	ClientIP    string
	IDC         string
}

type serverMetaContextKey struct{}

//初始化元信息, 存入context里
//参数: context, 服务名, 方法名
//如果要使用存入的元信息, 那么就用GetServerMeta来获取
func InitServerMeta(ctx context.Context, service string, method string) context.Context {
	meta := &ServerMeta{
		Method:      method,
		ServiceName: service,
	} //生成实例
	return context.WithValue(ctx, serverMetaContextKey{}, meta) //设置到context里.
}

//从context里获取元信息
func GetServerMeta(ctx context.Context) *ServerMeta {
	meta, ok := ctx.Value(serverMetaContextKey{}).(*ServerMeta)
	if !ok {
		meta = &ServerMeta{}
	}
	return meta
}
