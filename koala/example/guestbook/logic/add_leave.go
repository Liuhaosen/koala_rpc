package logic

import (
	"context"
	"modtest/gostudy/lesson2/ibinarytree/koala/example/guestbook/dal"
	"modtest/gostudy/lesson2/ibinarytree/koala/example/guestbook/model"
	"time"
)

type AddLeaveLogic struct {
	email   string
	content string
}

func NewAddLeaveLogic(email, content string) *AddLeaveLogic {
	return &AddLeaveLogic{
		email:   email,
		content: content,
	}
}

func (a *AddLeaveLogic) Execute(ctx context.Context) (err error) {
	leave := &model.Leave{
		Email:     a.email,
		Content:   a.content,
		Timestamp: time.Now().Unix(),
	}
	return dal.AddLeave(ctx, leave)
}
