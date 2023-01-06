package helloc

import (
	"context"
	"fmt"
	
	"modtest/gostudy/lesson2/ibinarytree/koala/tools/koala/client_example/generate/hello"
	
	"modtest/gostudy/lesson2/ibinarytree/koala/rpc"
	"modtest/gostudy/lesson2/ibinarytree/koala/meta"
	"modtest/gostudy/lesson2/ibinarytree/koala/errno"

	"google.golang.org/grpc"
)

type HelloClient struct {
	//继承HelloServiceClient
	serviceName string
	client *rpc.KoalaClient
}

func NewHelloClient(serviceName string, opts ...rpc.RpcOptionFunc) *HelloClient {
	return &HelloClient{
		serviceName: serviceName,
		client: rpc.NewKoalaClient(serviceName, opts...),
	}
}


func (s *HelloClient) SayHello(ctx context.Context, r *hello.HelloRequest, opts ...grpc.CallOption) (resp *hello.HelloResponse, err error) {
	/*
	//使用中间件的方式
	middlewareFunc := rpc.BuildClientMiddleware(mwClientSayHello)
	mkResp, err := middlewareFunc(ctx, r)
	if err != nil {
		return nil, err
	}
	*/

	mkResp, err := s.client.Call(ctx, "SayHello", r, mwClientSayHello)
	if err != nil {
		return nil, err
	}

	resp, ok := mkResp.(*hello.HelloResponse)
	if !ok {
		err = fmt.Errorf("invalid resp, not *hello.HelloResponse")
		return nil, err
	}

	return resp, err
}

func mwClientSayHello(ctx context.Context, request interface{}) (resp interface{}, err error) {
	/*
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		logs.Error(ctx, "did not connect: %v", err)
		return nil, err
	}
	*/

	rpcMeta := meta.GetRpcMeta(ctx)
	if rpcMeta.Conn == nil {
		return nil, errno.ConnFailed
	}

	req := request.(*hello.HelloRequest)
	client := hello.NewHelloServiceClient(rpcMeta.Conn)
	return client.SayHello(ctx, req)
}

