package loadbalance

import (
	"context"
	"math/rand"
	"modtest/gostudy/lesson2/ibinarytree/koala/registry"
)

//负载均衡-随机算法
type RandomBalance struct {
	name string
}

func NewRandomBalance() LoadBalance {
	return &RandomBalance{
		name: "random",
	}
}

//获取算法名
func (r *RandomBalance) Name() string {
	return "random"
}

//随机算法获取节点
func (r *RandomBalance) Select(ctx context.Context, nodes []*registry.Node) (node *registry.Node, err error) {
	if len(nodes) == 0 {
		err = ErrNotHaveNodes
		return
	}

	defer func() {
		if node != nil {
			setSelected(ctx, node)
		}
	}()

	var newNodes = filterNodes(ctx, nodes)
	if len(newNodes) == 0 {
		err = ErrAllNodeFailed
		return
	}

	var totalWeight int
	for _, node := range newNodes {
		if node.Weight == 0 {
			node.Weight = DefaultNodeWeight
		}
		totalWeight += node.Weight
	}

	curWeight := rand.Intn(totalWeight)
	curIndex := -1 //设置一个不存在的index下标
	for index, node := range nodes {
		curWeight -= node.Weight
		if curWeight < 0 {
			curIndex = index
			break
		}
	}
	if curIndex == -1 {
		//如果还等于-1, 那就说明没找到节点
		err = ErrNotHaveNodes
		return
	}
	node = nodes[curIndex]
	return
}

//随机算法获取节点
func (r *RandomBalance) SelectV1(ctx context.Context, nodes []*registry.Node) (node *registry.Node, err error) {
	if len(nodes) == 0 {
		err = ErrNotHaveNodes
		return
	}

	index := rand.Intn(len(nodes))
	node = nodes[index]
	return
}
