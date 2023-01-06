package main

import (
	"context"
	"modtest/gostudy/lesson2/ibinarytree/koala/example/hello/hello"
	"modtest/gostudy/lesson2/ibinarytree/koala/logs"
)

const (
	address     = "0.0.0.0:8080"
	defaultName = "world"
)

func myClientExample() {
	client := NewHelloClient("hello")
	ctx := context.Background()
	// resp, err := client.SayHelloV1(ctx, &hello.HelloRequest{Name: "test my client"})//第一版
	resp, err := client.SayHello(ctx, &hello.HelloRequest{Name: "test my client"}) //第二版, 使用中间件架构封装
	if err != nil {
		logs.Error(ctx, "could not greet %v", err)
		return
	}
	logs.Info(ctx, "Greeting : %v", resp.Reply)
}

func main() {
	logs.InitLogger(logs.LogLevelDebug, 20000, defaultName)
	myClientExample()
	logs.Stop()
}
