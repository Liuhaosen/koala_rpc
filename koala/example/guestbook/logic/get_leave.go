package logic

import (
	"context"
	"modtest/gostudy/lesson2/ibinarytree/koala/example/guestbook/dal"
	"modtest/gostudy/lesson2/ibinarytree/koala/example/guestbook/model"
)

type GetLeaveLogic struct {
	offset uint32
	limit  uint32
}

func NewGetLeaveLogic(offset, limit uint32) *GetLeaveLogic {
	return &GetLeaveLogic{
		offset: offset,
		limit:  limit,
	}
}

func (g *GetLeaveLogic) Execute(ctx context.Context) (result []*model.Leave, err error) {
	return dal.GetLeave(ctx, g.offset, g.limit)
}
