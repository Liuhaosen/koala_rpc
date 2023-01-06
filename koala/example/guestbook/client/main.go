package main

import (
	"context"
	"fmt"
	"modtest/gostudy/lesson2/ibinarytree/koala/example/guestbook/client/generate/client/guestbookc"
	"modtest/gostudy/lesson2/ibinarytree/koala/example/guestbook/client/generate/guestbook"
)

func main() {
	client := guestbookc.NewGuestbookClient("guestbook")
	r := &guestbook.AddLeaveRequest{
		Leave: &guestbook.Leave{
			Email:   "test@qq.com",
			Content: "测试",
		},
	}
	_, err := client.AddLeave(context.TODO(), r)

	fmt.Println("add leave result : ", err)

	getReq := &guestbook.GetLeaveRequest{
		Offset: 0,
		Limit:  10,
	}
	result, err := client.GetLeave(context.TODO(), getReq)
	if err != nil {
		fmt.Println("getleave failed, ", err)
		return
	}

	for _, leave := range result.Leaves {
		fmt.Println("email:", leave.Email, "content: ", leave.Content)
	}
}
