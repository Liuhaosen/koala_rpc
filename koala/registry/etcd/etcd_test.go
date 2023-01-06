package etcd

import (
	"context"
	"fmt"
	"modtest/gostudy/lesson2/ibinarytree/koala/registry"

	"testing"
	"time"
)

func TestRegister(t *testing.T) {
	registryInst, err := registry.InitPlugin(context.TODO(), "etcd",
		registry.WithAddrs([]string{"127.0.0.1:2379"}),
		registry.WithTimeout(time.Second),
		registry.WithRegistryPath("/ibinarytree/koala/"),
		registry.WithHeartBeat(5),
	)
	if err != nil {
		t.Errorf("初始化插件失败, 错误:%v", err)
		return
	}

	service := &registry.Service{
		Name: "comment_service",
	}
	service.Nodes = append(service.Nodes, &registry.Node{
		IP:   "127.0.0.1",
		Port: 8801,
	},
		&registry.Node{
			IP:   "127.0.0.2",
			Port: 8801,
		})
	err = registryInst.Register(context.TODO(), service)
	if err != nil {
		t.Errorf("服务注册失败, 错误:%v", err)
		return
	}
	go func() {
		time.Sleep(time.Second * 10)
		fmt.Println("添加新的节点")
		service.Nodes = append(service.Nodes, &registry.Node{
			IP:   "127.0.0.3",
			Port: 8801,
		},
			&registry.Node{
				IP:   "127.0.0.4",
				Port: 8801,
			},
		)
		err = registryInst.Register(context.TODO(), service)
		if err != nil {
			t.Errorf("服务注册失败, 错误:%v", err)
			return
		}
	}()
	// time.Sleep(time.Second * 3)
	for {
		//服务进程开启后就不能退出了. 写个for循环
		service, err := registryInst.GetService(context.TODO(), "comment_service")
		if err != nil {
			t.Errorf("get service failed, err : %v", err)
			return
		}
		for _, node := range service.Nodes {
			fmt.Printf("service:%s, node : %#v\n", service.Name, node)
		}
		fmt.Println()
		time.Sleep(5 * time.Second)
	}

}
