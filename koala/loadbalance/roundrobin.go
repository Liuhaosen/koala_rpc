package loadbalance

import (
	"context"
	"modtest/gostudy/lesson2/ibinarytree/koala/registry"
)

//负载均衡-轮询算法
type RoundRobinBalance struct {
	index int
	name  string
}

func NewRoundRobinBalance() LoadBalance {
	return &RoundRobinBalance{
		name: "roundrobin",
	}
}

func (r *RoundRobinBalance) Name() string {
	return "roundrobin"
}

func (r *RoundRobinBalance) Select(ctx context.Context, nodes []*registry.Node) (node *registry.Node, err error) {
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

	r.index = (r.index + 1) % len(nodes)
	node = nodes[r.index]
	return
}
