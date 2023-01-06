package guestbookc

import (
	"context"
	"fmt"

	"modtest/gostudy/lesson2/ibinarytree/koala/example/guestbook/client/generate/guestbook"

	"modtest/gostudy/lesson2/ibinarytree/koala/errno"
	"modtest/gostudy/lesson2/ibinarytree/koala/meta"
	"modtest/gostudy/lesson2/ibinarytree/koala/rpc"

	"google.golang.org/grpc"
)

type GuestbookClient struct {
	//继承HelloServiceClient
	serviceName string
	client      *rpc.KoalaClient
}

func NewGuestbookClient(serviceName string, opts ...rpc.RpcOptionFunc) *GuestbookClient {
	return &GuestbookClient{
		serviceName: serviceName,
		client:      rpc.NewKoalaClient(serviceName, opts...),
	}
}

func (s *GuestbookClient) AddLeave(ctx context.Context, r *guestbook.AddLeaveRequest, opts ...grpc.CallOption) (resp *guestbook.AddLeaveResponse, err error) {
	/*
		//使用中间件的方式
		middlewareFunc := rpc.BuildClientMiddleware(mwClientAddLeave)
		mkResp, err := middlewareFunc(ctx, r)
		if err != nil {
			return nil, err
		}
	*/

	mkResp, err := s.client.Call(ctx, "AddLeave", r, mwClientAddLeave)
	if err != nil {
		return nil, err
	}

	resp, ok := mkResp.(*guestbook.AddLeaveResponse)
	if !ok {
		err = fmt.Errorf("invalid resp, not *guestbook.AddLeaveResponse")
		return nil, err
	}

	return resp, err
}

func mwClientAddLeave(ctx context.Context, request interface{}) (resp interface{}, err error) {
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

	req := request.(*guestbook.AddLeaveRequest)
	client := guestbook.NewGuestBookServiceClient(rpcMeta.Conn)
	return client.AddLeave(ctx, req)
}

func (s *GuestbookClient) GetLeave(ctx context.Context, r *guestbook.GetLeaveRequest, opts ...grpc.CallOption) (resp *guestbook.GetLeaveResponse, err error) {
	/*
		//使用中间件的方式
		middlewareFunc := rpc.BuildClientMiddleware(mwClientGetLeave)
		mkResp, err := middlewareFunc(ctx, r)
		if err != nil {
			return nil, err
		}
	*/

	mkResp, err := s.client.Call(ctx, "GetLeave", r, mwClientGetLeave)
	if err != nil {
		return nil, err
	}

	resp, ok := mkResp.(*guestbook.GetLeaveResponse)
	if !ok {
		err = fmt.Errorf("invalid resp, not *guestbook.GetLeaveResponse")
		return nil, err
	}

	return resp, err
}

func mwClientGetLeave(ctx context.Context, request interface{}) (resp interface{}, err error) {
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

	req := request.(*guestbook.GetLeaveRequest)
	client := guestbook.NewGuestBookServiceClient(rpcMeta.Conn)
	return client.GetLeave(ctx, req)
}
