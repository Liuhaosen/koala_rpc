package main

import (
	"fmt"
	"math"
	"time"
)

//漏桶限流器
type BucketLimit struct {
	bucketSize float64 //漏桶最多能装的水大小
	rate       float64 //楼同种水的流出速率, 就是qps, 每秒能处理的请求数
	curWater   float64 //当前桶里的水
	unixNano   int64   //unix时间戳
}

func NewBucketLimit(bucketSize float64, rate float64) *BucketLimit {
	return &BucketLimit{
		bucketSize: bucketSize,
		rate:       rate,
		curWater:   0,
		unixNano:   time.Now().UnixNano(),
	}
}

func (b *BucketLimit) reflesh() {
	now := time.Now().UnixNano()
	//时间差,把纳秒换算成秒
	diffSec := float64(now-b.unixNano) / 1000 / 1000 / 1000
	b.curWater = math.Max(0, b.curWater-diffSec*b.rate) //qps * 时间就是这段时间内请求数的总量
	b.unixNano = now
}

func (b *BucketLimit) Allow() bool {
	//计算桶里的水
	b.reflesh()
	//如果当前的水量小于桶的最大承受水量, 则继续装水
	if b.curWater < b.bucketSize {
		b.curWater = b.curWater + 1
		return true
	}
	//满了就拒绝
	return false
}

func main() {
	//限速50qps, 桶大小100
	m := make(map[int]bool)
	limit := NewBucketLimit(100, 50)
	for i := 0; i < 1000; i++ {
		allow := limit.Allow()
		if allow {
			// fmt.Printf("i = %d is allowed\n", i)
			// continue
			m[i] = true

			continue
		}
		// fmt.Printf("i= %d is not allow\n", i)
		m[i] = false
		time.Sleep(time.Millisecond * 10)
	}

	for i := 0; i < 1000; i++ {
		fmt.Printf("i = %d allow= %v\n", i, m[i])
	}
}
