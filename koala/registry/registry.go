package registry

import "context"

//服务注册插件的接口
type Registry interface {
	Name() string                                                              //插件的名字
	Init(ctx context.Context, opts ...OptionHandler) (err error)               //初始化
	Register(ctx context.Context, service *Service) (err error)                //服务注册
	UnRegister(ctx context.Context, service *Service) (err error)              //服务反注册(停用)
	GetService(ctx context.Context, name string) (service *Service, err error) //服务发现: 通过服务的名字获取服务的位置信息(ip和port列表)
}
