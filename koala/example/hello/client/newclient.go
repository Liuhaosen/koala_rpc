package main

import (
	"context"
	"fmt"
	"modtest/gostudy/lesson2/ibinarytree/koala/example/hello/hello"
	"modtest/gostudy/lesson2/ibinarytree/koala/logs"
	"modtest/gostudy/lesson2/ibinarytree/koala/rpc"

	"google.golang.org/grpc"
)

type HelloClient struct {
	//继承HelloServiceClient
	serviceName string
	gclient     hello.HelloServiceClient
}

func NewHelloClient(serviceName string) *HelloClient {
	return &HelloClient{
		serviceName: serviceName,
	}
}

func (h *HelloClient) SayHelloV1(ctx context.Context, in *hello.HelloRequest, opts ...grpc.CallOption) (*hello.HelloResponse, error) {
	//1. 建立连接
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		logs.Error(context.Background(), "did not connect : %v", err)
		return nil, err
	}
	defer conn.Close()

	//2. 生成client实例
	h.gclient = hello.NewHelloServiceClient(conn)

	//3. 调用服务的方法
	return h.gclient.SayHello(ctx, in, opts...)
}

func (h *HelloClient) SayHello(ctx context.Context, in *hello.HelloRequest, opts ...grpc.CallOption) (*hello.HelloResponse, error) {
	//使用中间件的方式
	middlewareFunc := rpc.BuildClientMiddleware(mwClientSayHello)
	mkResp, err := middlewareFunc(ctx, in)
	if err != nil {
		return nil, err
	}

	resp, ok := mkResp.(*hello.HelloResponse)
	if !ok {
		err = fmt.Errorf("invalid resp, not *hello.HttpResponse")
		return nil, err
	}

	return resp, err
}

func mwClientSayHello(ctx context.Context, request interface{}) (resp interface{}, err error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		logs.Error(ctx, "did not connect: %v", err)
		return nil, err
	}

	req := request.(*hello.HelloRequest)
	defer conn.Close()
	client := hello.NewHelloServiceClient(conn)
	return client.SayHello(ctx, req)
}
