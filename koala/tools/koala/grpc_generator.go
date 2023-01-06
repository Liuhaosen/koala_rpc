package main

//使用官方工具, 生成generate目录下的grpc官方代码
import (
	"fmt"
	"os"
	"os/exec"
	"path"
)

type GrpcGenerator struct {
}

func (g *GrpcGenerator) Run(opt *Option, metaData *ServiceMetaData) (err error) {
	//protoc --go_out=plugins=grpc:. .\hello.proto
	//创建目录
	dir := path.Join(opt.Output, "generate", metaData.Package.Name)
	err = os.MkdirAll(dir, 0755)
	outputParams := fmt.Sprintf("plugins=grpc:%s/generate/%s", opt.Output, metaData.Package.Name)

	if err != nil {
		fmt.Println("创建目录失败, 错误:", err)
		return
	}
	//执行命令行
	cmd := exec.Command("protoc", "--go_out", outputParams, opt.Proto3Filename)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		fmt.Printf("grpc generator failed, err : %v\n", err)
		return
	}
	return
}

func init() {
	grpc := &GrpcGenerator{}
	RegisterServerGenerator("grpc generator", grpc)
	RegisterClientGenerator("grpc generator", grpc)
}
