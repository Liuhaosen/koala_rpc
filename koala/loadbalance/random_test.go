package loadbalance

import (
	"context"
	"fmt"
	"modtest/gostudy/lesson2/ibinarytree/koala/registry"
	"testing"
)

func TestSelect(t *testing.T) {
	balance := &RandomBalance{}
	var weights = [3]int{50, 100, 150}
	var nodes []*registry.Node
	for i := 0; i < 8; i++ {
		node := &registry.Node{
			IP:     fmt.Sprintf("127.0.0.%d", i),
			Port:   8080,
			Weight: weights[i%3],
		}
		fmt.Printf("node:%#v\n", node)
		nodes = append(nodes, node)
	}
	countStat := make(map[string]int)
	for i := 0; i < 1000; i++ {
		node, err := balance.Select(context.TODO(), nodes)
		if err != nil {
			t.Fatalf("select failed, err :%v", err)
			continue
		}
		countStat[node.IP]++
	}
	for key, val := range countStat {
		fmt.Printf("ip:%s, count: %v\n", key, val)
	}
}
