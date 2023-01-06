package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"modtest/gostudy/lesson2/ibinarytree/koala/logs"
	"modtest/gostudy/lesson2/ibinarytree/koala/registry"

	"path"
	"sync"
	"sync/atomic"
	"time"

	"github.com/coreos/etcd/clientv3"
)

const (
	MaxServiceNum          = 8                //服务最大数量
	MaxSyncServiceInterval = 10 * time.Second //定时任务的间隔时间. 表示隔多久进行一次服务同步到缓存
)

//etcd注册中心结构体
type EtcdRegistry struct {
	options            *registry.Options           //注册选项
	client             *clientv3.Client            //etcd的client
	serviceChan        chan *registry.Service      //要注册的服务先放到这里
	registryServiceMap map[string]*RegisterService //真正要异步去注册/已注册的服务map
	value              atomic.Value                //原子操作, 可以存任何数据. 存取无需加锁. 这里存储AllServiceInfo结构体.
	lock               sync.Mutex                  //value中没有缓存, 那么就到etcd中读取. 这时候要加锁
}

//要注册的服务结构体
type RegisterService struct {
	id            clientv3.LeaseID                        //租约id
	service       *registry.Service                       //要注册的服务
	registered    bool                                    //是否被注册过
	keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse //永久续约的返回
}

//所有服务的信息
//结构体存到EtcdRegistry的value里, 从value里读出来如果已经有缓存就直接返回.
type AllServiceInfo struct {
	serviceMap map[string]*registry.Service //存储所有service
}

//初始化结构体
var (
	etcdRegistry *EtcdRegistry = &EtcdRegistry{
		serviceChan:        make(chan *registry.Service, MaxServiceNum),
		registryServiceMap: make(map[string]*RegisterService, MaxServiceNum),
		options:            &registry.Options{},
	}
)

//启动etcd注册中心
func init() {
	allServiceInfo := &AllServiceInfo{
		serviceMap: make(map[string]*registry.Service, MaxServiceNum),
	}
	etcdRegistry.value.Store(allServiceInfo) //存取无需加锁, 性能好
	registry.RegisterPlugin(etcdRegistry)    //先将etcd插件注册到pluginMgr里
	go etcdRegistry.run()                    //初始化的时候就异步读取serviceChan, 把服务注册到etcd里
}

//插件的名字
func (e *EtcdRegistry) Name() string {
	return "etcd"
}

//初始化etcd
func (e *EtcdRegistry) Init(ctx context.Context, opts ...registry.OptionHandler) (err error) {
	// e.options = &registry.Options{}
	for _, opt := range opts {
		opt(e.options)
	}
	//初始化etcd的client, 连接上etcd
	e.client, err = clientv3.New(clientv3.Config{
		Endpoints:   e.options.Addrs,
		DialTimeout: e.options.Timeout,
	})

	if err != nil {
		return
	}

	return
}

//服务注册接口
//注册流程: 1. 将服务放到serviceChan里, 2. 异步从serviceChan里取出服务, 判断是否已加入map. 如果未加入就放到map中, 已加入map就不管. 3. 从map中遍历服务, 判断是否已注册服务, 已注册就续约, 未注册就注册.
func (e *EtcdRegistry) Register(ctx context.Context, service *registry.Service) (err error) {
	//将服务放到serviceChan里, 真正需要做的是run()函数
	select {
	case e.serviceChan <- service:
	default:
		err = fmt.Errorf("register chan is full")
		return
	}
	return
}

//服务反注册
func (e *EtcdRegistry) UnRegister(ctx context.Context, service *registry.Service) (err error) {

	return
}

//服务异步注册接口.
//这里会从serviceChan里读出service服务, 然后把服务放到map中, 然后轮询map.
func (e *EtcdRegistry) run() {
	ticker := time.NewTicker(MaxSyncServiceInterval)

	//通过clientv3租约的形式注册
	for {
		//从serviceChan里取出service, 放到map里. map里没有就注册到map
		select {
		case service := <-e.serviceChan:
			registryService, ok := e.registryServiceMap[service.Name]
			//1. 如果要注册的服务已经在已注册的map里了, 那就break
			if ok {
				// for _, node := range service.Nodes {
				// 	registryService.service.Nodes = append(registryService.service.Nodes, node)
				// }
				registryService.service.Nodes = append(registryService.service.Nodes, service.Nodes...)
				registryService.registered = false
				break
			}
			//2. 如果要注册的服务不在已注册的map里, 那就放到map中去
			registryService = &RegisterService{
				service: service,
			}
			e.registryServiceMap[service.Name] = registryService
		case <-ticker.C:
			//如果触发定时器就同步etcd服务到缓存
			fmt.Println("触发定时器")
			e.syncServiceFromEtcd()
		default:
			//如果管道里没有, 那就给map里的服务进行注册或者续约
			e.registerOrKeepAlive() //开始注册或者续约
			time.Sleep(time.Second) //防止返回空map在select里死循环
		}
	}

}

//轮询map, 判断是否注册.
func (e *EtcdRegistry) registerOrKeepAlive() {
	//遍历map, 挨个进行续约或者注册
	for _, registryService := range e.registryServiceMap {
		if registryService.registered {
			e.keepAlive(registryService) //如果已注册, 那就续约
			continue
		}
		e.registerService(registryService) //未注册, 那就注册
	}
}

//续约
func (e *EtcdRegistry) keepAlive(registryService *RegisterService) {
	//续约就是从keepalive返回的管道把它读出来即可
	//由于在注册的时候已经进行了续约, 所以这里没必要再次续约.
	//如果管道里为空, 说明到期或者出问题了. 把registered设置为false, 让它再次进入run()走default分支注册且续约即可
	resp := <-registryService.keepAliveChan
	if resp == nil {
		registryService.registered = false
		return
	}

	fmt.Printf("service: %v node: %v ttl:%v\n", registryService.service.Name, registryService.service.Nodes[0].IP, registryService.service.Nodes[0].Port)
}

//注册服务
//把服务注册到etcd的操作
func (e *EtcdRegistry) registerService(registryService *RegisterService) (err error) {
	resp, err := e.client.Grant(context.TODO(), e.options.HeartBeat) //租约授权
	if err != nil {
		logs.Error(context.TODO(), "租约授权失败, 错误: %v", err)
		return
	}
	registryService.id = resp.ID //要注册的服务的id = clientv3里的租期id

	for _, node := range registryService.service.Nodes {
		//遍历服务节点
		tempService := &registry.Service{
			Name: registryService.service.Name,
			Nodes: []*registry.Node{
				node,
			},
		}

		data, err := json.Marshal(tempService)
		if err != nil {
			logs.Warn(context.TODO(), "json marshal failed, err :%v", err)
			continue
		}
		key := e.serviceNodePath(tempService)
		fmt.Printf("register key : %s\n", key)
		//把服务的节点信息存到etcd中去  key=节点路径  value=节点node信息
		_, err = e.client.Put(context.TODO(), key, string(data), clientv3.WithLease(resp.ID))
		if err != nil {
			continue
		}
		//注册也需要续约. 这里使用永久续约.
		ch, err := e.client.KeepAlive(context.TODO(), resp.ID)
		if err != nil {
			continue
		}
		registryService.keepAliveChan = ch
		registryService.registered = true
	}
	return
}

//获取服务在etcd里的路径
func (e *EtcdRegistry) serviceNodePath(registryService *registry.Service) string {
	nodeIP := fmt.Sprintf("%s:%d", registryService.Nodes[0].IP, registryService.Nodes[0].Port) //获取ip+port
	return path.Join(e.options.RegistryPath, registryService.Name, nodeIP)                     //拼接服务路径+服务名+ip+port
}

//获取服务在etcd里的路径前缀 (就是没有ip+port的部分)
func (e *EtcdRegistry) servicePath(name string) string {
	return path.Join(e.options.RegistryPath, name)
}

//服务发现
func (e *EtcdRegistry) GetService(ctx context.Context, name string) (service *registry.Service, err error) {
	//取出value, 根据name查看该服务是否存在.
	//一般情况下, 都会从缓存读取
	service, ok := e.getServiceFromCache(ctx, name)
	if ok {
		return
	}

	//如果缓存中没有service, 就从etcd中读取.
	//加锁保护一下etcd, 防止多个请求同时过来
	//如果多个请求过来, 只允许一个请求操作etcd, 将etcd放到缓存中去.
	//然后其他请求就可以到缓存里取服务信息.
	e.lock.Lock()
	defer e.lock.Unlock()
	service, ok = e.getServiceFromCache(ctx, name)
	if ok {
		return
	}

	//从etcd中读取指定服务名字的服务信息
	key := e.servicePath(name)
	resp, err := e.client.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		return
	}

	//把读取到的信息放到service结构体里
	service = &registry.Service{
		Name: name,
	}
	for _, kv := range resp.Kvs {
		// fmt.Printf("index: %v, key: %v, value: %v\n", index, string(kv.Key), string(kv.Value))
		//打印结果: index: 1, key: /ibinarytree/koala/comment_service/127.0.0.2:8801, value: {"name":"comment_service","nodes":[{"id":"","ip":"127.0.0.3","port":8801}]}
		value := kv.Value
		var tempService registry.Service
		err = json.Unmarshal(value, &tempService)
		if err != nil {
			return
		}
		service.Nodes = append(service.Nodes, tempService.Nodes...)
		// for _, node := range tempService.Nodes {
		// 	service.Nodes = append(service.Nodes, node)
		// }
	}

	//放到AllServiceInfo里缓存起来, 这样别的请求就能直接从缓存里读到该服务
	oldAllServiceInfo := e.value.Load().(*AllServiceInfo) //先把服务读出来
	var newAllServiceInfo = &AllServiceInfo{
		serviceMap: make(map[string]*registry.Service, MaxServiceNum),
	}
	for key, val := range oldAllServiceInfo.serviceMap {
		newAllServiceInfo.serviceMap[key] = val
	} //存到AllServiceInfo结构体里

	newAllServiceInfo.serviceMap[name] = service
	//存到etcdRegistry的value字段里
	e.value.Store(newAllServiceInfo)

	return
}

//封装从缓存中获取服务
func (e *EtcdRegistry) getServiceFromCache(ctx context.Context, name string) (service *registry.Service, ok bool) {
	allServiceInfo := e.value.Load().(*AllServiceInfo)
	service, ok = allServiceInfo.serviceMap[name]
	return
}

//从etcd把服务同步到缓存中
func (e *EtcdRegistry) syncServiceFromEtcd() {
	var allServiceInfoNew = &AllServiceInfo{
		serviceMap: make(map[string]*registry.Service, MaxServiceNum),
	}
	ctx := context.TODO()
	allServiceInfo := e.value.Load().(*AllServiceInfo)
	//每个服务都要在etcd中更新
	for _, service := range allServiceInfo.serviceMap {
		key := e.servicePath(service.Name)
		resp, err := e.client.Get(ctx, key, clientv3.WithPrefix())
		if err != nil {
			allServiceInfoNew.serviceMap[service.Name] = service //更新失败的话, 把旧的服务存回去
			continue
		}

		//把读取到的信息放到service结构体里
		serviceNew := &registry.Service{
			Name: service.Name,
		}
		for _, kv := range resp.Kvs {
			fmt.Printf("key : %s, value : %s\n", string(kv.Key), string(kv.Value))
			value := kv.Value
			var tempService registry.Service
			err = json.Unmarshal(value, &tempService)
			if err != nil {
				fmt.Println("unmarshal failed , err : ", err, "value: ", kv.Value)
				return
			}
			serviceNew.Nodes = append(serviceNew.Nodes, tempService.Nodes...)
		}
		//加入缓存
		// allServiceInfoNew.serviceMap[serviceNew.Name] = serviceNew

		//加入缓存
		allServiceInfoOld := e.value.Load().(*AllServiceInfo)

		for _, service := range allServiceInfoOld.serviceMap {
			allServiceInfoNew.serviceMap[service.Name] = service
		}

		e.value.Store(allServiceInfoNew)
	}
	// e.value.Store(allServiceInfoNew)
	fmt.Printf("更新所有service成功, len:%d, serviceMap:%v\n", len(allServiceInfoNew.serviceMap), allServiceInfoNew.serviceMap)
}

func InitEtcd(ctx context.Context, opts ...registry.OptionHandler) {
	etcdRegistry.Init(ctx, opts...)
}
