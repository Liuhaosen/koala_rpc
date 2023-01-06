package loadbalance

import (
	"context"
	"errors"
	"modtest/gostudy/lesson2/ibinarytree/koala/registry"
)

var (
	ErrNotHaveNodes  = errors.New("没有节点")
	ErrAllNodeFailed = errors.New("所有节点尝试失败")
)

const (
	DefaultNodeWeight = 100
)

type LoadBalance interface {
	Name() string
	Select(ctx context.Context, nodes []*registry.Node) (node *registry.Node, err error)
}
