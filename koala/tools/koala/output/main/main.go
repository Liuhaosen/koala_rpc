package main

import (
	"context"
	"modtest/gostudy/lesson2/ibinarytree/koala/logs"
	"modtest/gostudy/lesson2/ibinarytree/koala/server"

	"modtest/gostudy/lesson2/ibinarytree/koala/tools/koala/output/generate/hello"
	"modtest/gostudy/lesson2/ibinarytree/koala/tools/koala/output/router"
)

var routerServer = &router.RouterServer{}

func main() {
	err := server.Init("hello")
	if err != nil {
		logs.Error(context.TODO(), "初始化服务失败, 错误: %v", err)
		logs.Stop()
		return
	}
	// hello.RegisterHelloServiceServer()
	hello.RegisterHelloServiceServer(server.GRPCServer(), routerServer)
	server.Run()
}
