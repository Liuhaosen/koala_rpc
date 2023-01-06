package controller

import (
	"context"
	"modtest/gostudy/lesson2/ibinarytree/koala/example/guestbook/client/generate/guestbook"
	"modtest/gostudy/lesson2/ibinarytree/koala/example/guestbook/logic"
	"modtest/gostudy/lesson2/ibinarytree/koala/logs"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AddLeaveController struct {
}

//检查请求参数, 如果该函数返回错误, 则Run函数不会执行
func (c *AddLeaveController) CheckParams(ctx context.Context, r *guestbook.AddLeaveRequest) (err error) {
	if len(r.Leave.Email) == 0 || len(r.Leave.Content) == 0 {
		err = status.Errorf(codes.InvalidArgument, "add leave failed")
		return
	}
	return
}

//SayHello函数的实现
func (c *AddLeaveController) Run(ctx context.Context, r *guestbook.AddLeaveRequest) (resp *guestbook.AddLeaveResponse, err error) {
	resp = &guestbook.AddLeaveResponse{}
	addLeave := logic.NewAddLeaveLogic(r.Leave.GetEmail(), r.Leave.GetContent())
	err = addLeave.Execute(ctx)
	if err != nil {
		logs.Error(ctx, "add leave failed, err : %v", err)
		return
	}
	return
}
