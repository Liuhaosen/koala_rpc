package main

import (
	"context"
	"modtest/gostudy/lesson2/ibinarytree/koala/logs"
	_ "modtest/gostudy/lesson2/ibinarytree/koala/registry/etcd"
	"modtest/gostudy/lesson2/ibinarytree/koala/rpc"
	"modtest/gostudy/lesson2/ibinarytree/koala/tools/koala/client_example/generate/client/helloc"
	"modtest/gostudy/lesson2/ibinarytree/koala/tools/koala/client_example/generate/hello"
	"time"
)

const (
	address     = "0.0.0.0:8080"
	defaultName = "world"
)

func main() {
	client := helloc.NewHelloClient("hello",
		rpc.WithLimitQPS(10000),
		rpc.WithClientServiceName("hello-client-example"),
	)
	var count int
	ctx := context.Background()
	// 尝试禁止循环调用caller
	for {
		count++
		resp, err := client.SayHello(ctx, &hello.HelloRequest{Name: "test my client"})
		if err != nil {
			if count%100 == 0 {
				// fmt.Printf("could not greet : %v\n", err)
				logs.Error(ctx, "could not greet : %v", err)
				return
			}
			logs.Warn(ctx, "err : ", err)
			time.Sleep(time.Second)
			continue
		}
		if count%100 == 0 {
			// fmt.Printf("Greeting: %s\n", resp.Reply)
			logs.Info(ctx, "Greeting: %s", resp.Reply)
		}

	}
	// logs.Stop()
}

/* func main() {
	client := helloc.NewHelloClient("hello",
		rpc.WithClientServiceName("hello-client-example"),
	)
	ctx := context.TODO()
	resp, err := client.SayHello(ctx, &hello.HelloRequest{Name: "test my client"})
	if err != nil {

		logs.Error(ctx, "could not greet : %v", err)

		return
	}

	logs.Info(ctx, "Greeting: %s", resp.Reply)

	logs.Stop()
}
*/
