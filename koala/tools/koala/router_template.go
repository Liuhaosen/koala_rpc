package main

var router_template = `package router
	import(
		"context"
		"fmt"
		"modtest/gostudy/lesson2/ibinarytree/koala/meta"
		"modtest/gostudy/lesson2/ibinarytree/koala/server"
		{{if not .Prefix }}
		"generate/{{.Package.Name}}"
		"controller"
		{{else}}
		"{{.Prefix}}/generate/{{.Package.Name}}"
		"{{.Prefix}}/controller"
		{{end}}
		
	)
	type RouterServer struct{}

	{{range .Rpc}}
	func (s *RouterServer) {{.Name}} (ctx context.Context, r *{{$.Package.Name}}.{{.RequestType}}) (resp *{{$.Package.Name}}.{{.ReturnsType}}, err error){
		ctx = meta.InitServerMeta(ctx, "{{$.Package.Name}}", "{{.Name}}") //将元信息设置到context里了, 传入了服务名和方法名
		mwFunc := server.BuildServerMiddleware(mw{{.Name}})
		mwResp, err := mwFunc(ctx, r)
		if err != nil{
			return 
		}
		resp = mwResp.(*{{$.Package.Name}}.{{.ReturnsType}})
		return
	}


	func mw{{.Name}}(ctx context.Context, request interface{}) (resp interface{}, err error) {
		ctrl := &controller.{{.Name}}Controller{}
		r := request.(*{{$.Package.Name}}.{{.RequestType}})
		err = ctrl.CheckParams(ctx, r)
		if err != nil {
			fmt.Println("检查参数有误, err : ", err)
			return
		}
	
		resp, err = ctrl.Run(ctx, r)
		return
	}
	{{end}}
`
