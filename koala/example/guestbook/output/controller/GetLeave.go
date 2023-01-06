package controller

import (
	"context"

	"modtest/gostudy/lesson2/ibinarytree/koala/example/guestbook/output/generate/guestbook"
)

type GetLeaveController struct{}

//检查请求参数, 如果该函数返回错误, 则Run函数不会执行

func (s *GetLeaveController) CheckParams(ctx context.Context, r *guestbook.GetLeaveRequest) (err error) {
	return
}

func (s *GetLeaveController) Run(ctx context.Context, r *guestbook.GetLeaveRequest) (resp *guestbook.GetLeaveResponse, err error) {
	resp = &guestbook.GetLeaveResponse{}
	return
}
