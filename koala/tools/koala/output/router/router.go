package router
	import(
		"context"
		"fmt"
		"modtest/gostudy/lesson2/ibinarytree/koala/meta"
		"modtest/gostudy/lesson2/ibinarytree/koala/server"
		
		"modtest/gostudy/lesson2/ibinarytree/koala/tools/koala/output/generate/hello"
		"modtest/gostudy/lesson2/ibinarytree/koala/tools/koala/output/controller"
		
		
	)
	type RouterServer struct{}

	
	func (s *RouterServer) SayHello (ctx context.Context, r *hello.HelloRequest) (resp *hello.HelloResponse, err error){
		ctx = meta.InitServerMeta(ctx, "hello", "SayHello") //将元信息设置到context里了, 传入了服务名和方法名
		mwFunc := server.BuildServerMiddleware(mwSayHello)
		mwResp, err := mwFunc(ctx, r)
		if err != nil{
			return 
		}
		resp = mwResp.(*hello.HelloResponse)
		return
	}


	func mwSayHello(ctx context.Context, request interface{}) (resp interface{}, err error) {
		ctrl := &controller.SayHelloController{}
		r := request.(*hello.HelloRequest)
		err = ctrl.CheckParams(ctx, r)
		if err != nil {
			fmt.Println("检查参数有误, err : ", err)
			return
		}
	
		resp, err = ctrl.Run(ctx, r)
		return
	}
	
