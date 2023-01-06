package loadbalance

//负载均衡-加权随机算法
import (
	"math/rand"
	"modtest/gostudy/lesson2/ibinarytree/koala/registry"

	"golang.org/x/net/context"
)

type WeightedRandomBalance struct {
}

func (w *WeightedRandomBalance) Name() string {
	return "weighted_random"
}

//加权随机算法获取节点
func (w *WeightedRandomBalance) Select(ctx context.Context, nodes []*registry.Node) (node *registry.Node, err error) {
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
	for _, node := range nodes {
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
