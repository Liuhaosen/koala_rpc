package controller

import (
	"context"

	"modtest/gostudy/lesson2/ibinarytree/koala/example/guestbook/output/generate/guestbook"
)

type AddLeaveController struct{}

//检查请求参数, 如果该函数返回错误, 则Run函数不会执行

func (s *AddLeaveController) CheckParams(ctx context.Context, r *guestbook.AddLeaveRequest) (err error) {
	return
}

func (s *AddLeaveController) Run(ctx context.Context, r *guestbook.AddLeaveRequest) (resp *guestbook.AddLeaveResponse, err error) {
	resp = &guestbook.AddLeaveResponse{}
	return
}
