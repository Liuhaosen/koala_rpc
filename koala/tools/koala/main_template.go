package main

var main_template = `
package main
import (
	"context"
	"modtest/gostudy/lesson2/ibinarytree/koala/logs"
	"modtest/gostudy/lesson2/ibinarytree/koala/server"
	{{if not .Prefix }}
	"router"
	"generate/{{.Package.Name}}"
	{{else}}
	"{{.Prefix}}/generate/{{.Package.Name}}"
	"{{.Prefix}}/router"
	{{end}}	
)	

var routerServer = &router.RouterServer{}

func main(){
	err := server.Init("{{.Package.Name}}")
	if err != nil{
		logs.Error(context.TODO(), "初始化服务失败, 错误: %` + `v", err)
		logs.Stop()
		return
	}
	
	{{.Package.Name}}.Register{{Capitalize .Server.Name}}Server(server.GRPCServer(), routerServer)
	server.Run()
}
`
