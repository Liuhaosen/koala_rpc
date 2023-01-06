package main

import (
	"fmt"
	"html/template"
	"modtest/gostudy/lesson2/ibinarytree/koala/util"
	"os"
	"path"

	"github.com/emicklei/proto"
)

// 生成controller的接口实现代码
// (1) 使用github.com/emicklei/proto解析proto3的idl文件
// (2) 拿到idl定义的方法和结构体等元数据
// (3) 遍历元数据, 生成接口实现代码

type CtrlGenerator struct {
	// service  *proto.Service
	// messages []*proto.Message
	// rpc      []*proto.RPC
}

type RpcMeta struct {
	Rpc     *proto.RPC
	Package *proto.Package
	Prefix  string
}

func init() {
	ctrl := &CtrlGenerator{}
	RegisterServerGenerator("controller generator", ctrl)
}

func (c *CtrlGenerator) Run(opt *Option, metaData *ServiceMetaData) (err error) {
	reader, err := os.Open(opt.Proto3Filename)
	if err != nil {
		fmt.Printf("open file %s failed, err : %v\n", opt.Proto3Filename, err)
		return
	}
	defer reader.Close()
	//解析proto文件
	// parser := proto.NewParser(reader)
	// definition, err := parser.Parse()
	// if err != nil {
	// 	fmt.Printf("解析%s失败, 错误:%v", opt.Proto3Filename, err)
	// 	return
	// }
	// proto.Walk(definition,
	// 	//解析到service信息, 就使用handleService回调函数
	// 	proto.WithService(c.handleService),
	// 	proto.WithMessage(c.handleMessage),
	// 	proto.WithRPC(c.handleRPC),
	// )
	// fmt.Printf("parse protoc succ, rpc:%#v\n", c.rpc)
	return c.generateRpc(opt, metaData)
}

func (c *CtrlGenerator) generateRpc(opt *Option, metaData *ServiceMetaData) (err error) {
	for _, rpc := range metaData.Rpc {
		//1. 生成文件名
		filename := path.Join("./", opt.Output, "controller", rpc.Name)

		//2. 判断是否存在, 防止覆盖写入
		exist := util.IsFileExist(filename)
		if exist {
			continue
		}

		//3. 打开文件句柄
		file, err := os.OpenFile(filename+".go", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			fmt.Printf("open file :%s failed, err : %v\n", filename, err)
			return err
		}

		//4. 获取元数据
		rpcMeta := &RpcMeta{}
		rpcMeta.Package = metaData.Package
		rpcMeta.Rpc = rpc
		rpcMeta.Prefix = metaData.Prefix

		//5. 渲染写入.
		err = c.render(file, ctrl_template, rpcMeta)
		if err != nil {
			fmt.Println("render ctrl_template failed , err :", err)
			return err
		}
		defer file.Close()
	}
	return
}

//渲染文件. 把template的内容渲染到文件里去.
func (c *CtrlGenerator) render(file *os.File, data string, metaData *RpcMeta) (err error) {
	t := template.New("main")
	t, err = t.Parse(data)
	if err != nil {
		fmt.Println("parse ctrl_template failed, err:", err)
		return
	}
	err = t.Execute(file, metaData)
	return
}
