package main

var ctrl_template = `package controller

import (
	"context"
	{{if not .Prefix}}
	"generate/{{.Package.Name}}"
	{{else}}
	"{{.Prefix}}/generate/{{.Package.Name}}"
	{{end}}
)

type {{.Rpc.Name}}Controller struct{}

//检查请求参数, 如果该函数返回错误, 则Run函数不会执行

func (s *{{.Rpc.Name}}Controller) CheckParams (ctx context.Context, r *{{.Package.Name}}.{{.Rpc.RequestType}}) (err error){
	return
}

func (s *{{.Rpc.Name}}Controller) Run(ctx context.Context, r *{{.Package.Name}}.{{.Rpc.RequestType}}) (resp *{{.Package.Name}}.{{.Rpc.ReturnsType}}, err error) {
	resp = &{{.Package.Name}}.{{.Rpc.ReturnsType}}{
		Reply: "我是通过代码生成器生成的脚手架",
	}
	return
}
`
