package loadbalance

import (
	"context"
	"fmt"
	"modtest/gostudy/lesson2/ibinarytree/koala/logs"
	"modtest/gostudy/lesson2/ibinarytree/koala/registry"
)

type selectedNodes struct {
	selectedNodeMap map[string]bool
}

type loadbalaceFilterNodes struct{}

func WithBalanceContext(ctx context.Context) context.Context {
	sel := &selectedNodes{
		selectedNodeMap: make(map[string]bool),
	}
	return context.WithValue(ctx, loadbalaceFilterNodes{}, sel)
}

func GetSelectedNodes(ctx context.Context) *selectedNodes {
	sel, ok := ctx.Value(loadbalaceFilterNodes{}).(*selectedNodes)
	if !ok {
		return nil
	}
	return sel
}

//过滤已选择的节点
func filterNodes(ctx context.Context, nodes []*registry.Node) []*registry.Node {
	var newNodes []*registry.Node
	sel := GetSelectedNodes(ctx)
	if sel == nil {
		return newNodes
	}

	for _, node := range nodes {
		addr := fmt.Sprintf("%s:%d", node.IP, node.Port)
		_, ok := sel.selectedNodeMap[addr]
		if ok {
			logs.Debug(ctx, "addr:%s ok", addr)
			continue
		}
		newNodes = append(newNodes, node)
	}
	return newNodes
}

func setSelected(ctx context.Context, node *registry.Node) {
	sel := GetSelectedNodes(ctx)
	if sel == nil {
		return
	}

	addr := fmt.Sprintf("%s:%d", node.IP, node.Port)
	logs.Debug(ctx, "filter node:%s", addr)
	sel.selectedNodeMap[addr] = true
}
