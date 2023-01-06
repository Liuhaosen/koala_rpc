package main

import (
	"context"
	"modtest/gostudy/lesson2/ibinarytree/koala/logs"
	"modtest/gostudy/lesson2/ibinarytree/koala/server"

	"modtest/gostudy/lesson2/ibinarytree/koala/example/guestbook/output/generate/guestbook"
	"modtest/gostudy/lesson2/ibinarytree/koala/example/guestbook/output/router"
)

var routerServer = &router.RouterServer{}

func main() {
	err := server.Init("guestbook")
	if err != nil {
		logs.Error(context.TODO(), "初始化服务失败, 错误: %v", err)
		logs.Stop()
		return
	}

	guestbook.RegisterGuestBookServiceServer(server.GRPCServer(), routerServer)
	server.Run()
}
