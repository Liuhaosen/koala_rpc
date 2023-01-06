package main

var rpc_client_template = `package {{.Package.Name}}c

import (
	"context"
	"fmt"
	{{if not .Prefix}}
	"generate/{{.Package.Name}}"
	{{else}}
	"{{.Prefix}}/generate/{{.Package.Name}}"
	{{end}}
	"modtest/gostudy/lesson2/ibinarytree/koala/rpc"
	"modtest/gostudy/lesson2/ibinarytree/koala/meta"
	"modtest/gostudy/lesson2/ibinarytree/koala/errno"

	"google.golang.org/grpc"
)

type {{Capitalize .Package.Name}}Client struct {
	//继承HelloServiceClient
	serviceName string
	client *rpc.KoalaClient
}

func New{{Capitalize .Package.Name}}Client(serviceName string, opts ...rpc.RpcOptionFunc) *{{Capitalize .Package.Name}}Client {
	return &{{Capitalize .Package.Name}}Client{
		serviceName: serviceName,
		client: rpc.NewKoalaClient(serviceName, opts...),
	}
}

{{range .Rpc}}
func (s *{{Capitalize $.Package.Name}}Client) {{.Name}}(ctx context.Context, r *{{$.Package.Name}}.{{.RequestType}}, opts ...grpc.CallOption) (resp *{{$.Package.Name}}.{{.ReturnsType}}, err error) {
	/*
	//使用中间件的方式
	middlewareFunc := rpc.BuildClientMiddleware(mwClient{{.Name}})
	mkResp, err := middlewareFunc(ctx, r)
	if err != nil {
		return nil, err
	}
	*/

	mkResp, err := s.client.Call(ctx, "{{.Name}}", r, mwClient{{.Name}})
	if err != nil {
		return nil, err
	}

	resp, ok := mkResp.(*{{$.Package.Name}}.{{.ReturnsType}})
	if !ok {
		err = fmt.Errorf("invalid resp, not *{{$.Package.Name}}.{{.ReturnsType}}")
		return nil, err
	}

	return resp, err
}

func mwClient{{.Name}}(ctx context.Context, request interface{}) (resp interface{}, err error) {
	/*
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		logs.Error(ctx, "did not connect: %` + `v", err)
		return nil, err
	}
	*/

	rpcMeta := meta.GetRpcMeta(ctx)
	if rpcMeta.Conn == nil {
		return nil, errno.ConnFailed
	}

	req := request.(*{{$.Package.Name}}.{{.RequestType}})
	client := {{$.Package.Name}}.New{{Capitalize $.Package.Name}}ServiceClient(rpcMeta.Conn)
	return client.{{.Name}}(ctx, req)
}
{{end}}
`
