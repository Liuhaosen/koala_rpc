package registry

import (
	"time"
)

//注册选项, 使用选项模式
type Options struct {
	Addrs        []string      //注册中心地址
	Timeout      time.Duration //与注册中心交互的超时时间
	RegistryPath string        //在etcd里的注册路径, 例如: xxx_company/app/douyin/service_A/10.192.1.1:8801, service_A可能有多个节点.
	HeartBeat    int64         //续约心跳时间
}

//处理选项方法的结构, 函数式编程
type OptionHandler func(opts *Options)

//处理地址字段
func WithAddrs(addrs []string) OptionHandler {
	//在服务初始化的过程中, 遍历所有处理选项方法, 然后对选项的每个字段进行赋值
	return func(opts *Options) {
		opts.Addrs = addrs
	}
}

//处理超时字段
func WithTimeout(timeout time.Duration) OptionHandler {
	//在服务初始化的过程中, 遍历所有处理选项方法, 然后对选项的每个字段进行赋值
	return func(opts *Options) {
		opts.Timeout = timeout
	}
}

//处理注册路径
func WithRegistryPath(path string) OptionHandler {
	return func(opts *Options) {
		opts.RegistryPath = path
	}
}

//处理心跳连接
func WithHeartBeat(heartBeat int64) OptionHandler {
	return func(opts *Options) {
		opts.HeartBeat = heartBeat
	}
}
