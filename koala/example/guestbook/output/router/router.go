package router
	import(
		"context"
		"fmt"
		"modtest/gostudy/lesson2/ibinarytree/koala/meta"
		"modtest/gostudy/lesson2/ibinarytree/koala/server"
		
		"modtest/gostudy/lesson2/ibinarytree/koala/example/guestbook/output/generate/guestbook"
		"modtest/gostudy/lesson2/ibinarytree/koala/example/guestbook/output/controller"
		
		
	)
	type RouterServer struct{}

	
	func (s *RouterServer) AddLeave (ctx context.Context, r *guestbook.AddLeaveRequest) (resp *guestbook.AddLeaveResponse, err error){
		ctx = meta.InitServerMeta(ctx, "guestbook", "AddLeave") //将元信息设置到context里了, 传入了服务名和方法名
		mwFunc := server.BuildServerMiddleware(mwAddLeave)
		mwResp, err := mwFunc(ctx, r)
		if err != nil{
			return 
		}
		resp = mwResp.(*guestbook.AddLeaveResponse)
		return
	}


	func mwAddLeave(ctx context.Context, request interface{}) (resp interface{}, err error) {
		ctrl := &controller.AddLeaveController{}
		r := request.(*guestbook.AddLeaveRequest)
		err = ctrl.CheckParams(ctx, r)
		if err != nil {
			fmt.Println("检查参数有误, err : ", err)
			return
		}
	
		resp, err = ctrl.Run(ctx, r)
		return
	}
	
	func (s *RouterServer) GetLeave (ctx context.Context, r *guestbook.GetLeaveRequest) (resp *guestbook.GetLeaveResponse, err error){
		ctx = meta.InitServerMeta(ctx, "guestbook", "GetLeave") //将元信息设置到context里了, 传入了服务名和方法名
		mwFunc := server.BuildServerMiddleware(mwGetLeave)
		mwResp, err := mwFunc(ctx, r)
		if err != nil{
			return 
		}
		resp = mwResp.(*guestbook.GetLeaveResponse)
		return
	}


	func mwGetLeave(ctx context.Context, request interface{}) (resp interface{}, err error) {
		ctrl := &controller.GetLeaveController{}
		r := request.(*guestbook.GetLeaveRequest)
		err = ctrl.CheckParams(ctx, r)
		if err != nil {
			fmt.Println("检查参数有误, err : ", err)
			return
		}
	
		resp, err = ctrl.Run(ctx, r)
		return
	}
	
