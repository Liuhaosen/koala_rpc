package main

import (
	"context"
	"fmt"
	pb "modtest/gostudy/lesson2/ibinarytree/koala/example/hello/hello"
	"net"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type server struct{}

//实现hello.proto里的SayHello方法
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{
		Reply: "你好 " + in.Name,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", port) //监听端口
	if err != nil {
		fmt.Println("listen failed, err :", err)
	}

	s := grpc.NewServer()                       //生成新的服务实例
	pb.RegisterHelloServiceServer(s, &server{}) //注册到hello包里
	s.Serve(lis)                                //开启服务
}
