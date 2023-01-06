package dal

import (
	"context"
	"modtest/gostudy/lesson2/ibinarytree/koala/example/guestbook/model"
)

type LeaveStoreMgr struct {
	leaveList []*model.Leave
}

var leaveStoreMgr = &LeaveStoreMgr{}

//添加留言
func (l *LeaveStoreMgr) AddLeave(ctx context.Context, leave *model.Leave) (err error) {
	l.leaveList = append(l.leaveList, leave)
	return
}

//获取留言
func (l *LeaveStoreMgr) GetLeave(ctx context.Context, offset, limit uint32) (result []*model.Leave, err error) {
	if offset < 0 || limit <= 0 {
		return
	}

	if offset >= uint32(len(l.leaveList)) {
		return
	}

	result = l.leaveList[offset : offset+limit]
	return
}

//添加留言, 这样不用再往外暴露leaveStoreMgr
func AddLeave(ctx context.Context, leave *model.Leave) error {
	return leaveStoreMgr.AddLeave(ctx, leave)
}

func GetLeave(ctx context.Context, offset, limit uint32) (result []*model.Leave, err error) {
	return leaveStoreMgr.GetLeave(ctx, offset, limit)
}
