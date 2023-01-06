package registry

import (
	"context"
	"fmt"
	"sync"
)

//插件管理

//定义全局变量, 初始化pluginMgr里的map
var (
	pluginMgr = &PluginMgr{
		plugins: make(map[string]Registry),
	}
)

type PluginMgr struct {
	plugins map[string]Registry //使用map管理注册插件, 每个注册插件都要放到这个map里
	lock    sync.Mutex
}

//注册插件.  将写好的插件放入到插件管理的map中
func (p *PluginMgr) registerPlugin(plugin Registry) (err error) {
	p.lock.Lock()
	defer p.lock.Unlock()

	_, ok := p.plugins[plugin.Name()]
	if ok {
		err = fmt.Errorf("duplicate registry plugin")
		return
	}
	p.plugins[plugin.Name()] = plugin
	return
}

//初始化插件.  要使用的插件需要先初始化.
//参数1: context给Registry接口使用.
//参数2: name插件名(etcd/consul), 要根据插件名找到注册在pluginMgr里的插件
//参数3: 一个或多个处理插件初始化选项的函数
func (p *PluginMgr) initPlugin(ctx context.Context, name string, opts ...OptionHandler) (registry Registry, err error) {
	//1. 查找对应的插件是否存在
	p.lock.Lock()
	defer p.lock.Unlock()
	plugin, ok := p.plugins[name]
	if !ok {
		err = fmt.Errorf("plugin %s is not exist", name)
		return
	}

	//2. 初始化插件, 返回后用户就可以使用插件的注册/反注册方法
	registry = plugin
	err = registry.Init(ctx, opts...)
	return
}

//注册插件
//暴露给外部使用, 写好的插件可以调用该函数, 注册到pluginMgr中进行管理
func RegisterPlugin(registry Registry) (err error) {
	return pluginMgr.registerPlugin(registry)
}

//初始化某个注册中心插件, 对外暴露用来初始化注册中心使用.
//将写好的注册中心插件在注册到管理类完成后就可以初始化
func InitPlugin(ctx context.Context, name string, opts ...OptionHandler) (registry Registry, err error) {
	return pluginMgr.initPlugin(ctx, name, opts...)
}
