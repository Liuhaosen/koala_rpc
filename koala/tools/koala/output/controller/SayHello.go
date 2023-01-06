package controller

import (
	"context"
	
	"modtest/gostudy/lesson2/ibinarytree/koala/tools/koala/output/generate/hello"
	
)

type SayHelloController struct{}

//检查请求参数, 如果该函数返回错误, 则Run函数不会执行

func (s *SayHelloController) CheckParams (ctx context.Context, r *hello.HelloRequest) (err error){
	return
}

func (s *SayHelloController) Run(ctx context.Context, r *hello.HelloRequest) (resp *hello.HelloResponse, err error) {
	resp = &hello.HelloResponse{
		Reply: "我是通过代码生成器生成的脚手架",
	}
	return
}
