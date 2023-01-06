package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

//计数器限流
type CounterLimit struct {
	counter      int64 //计数器
	limit        int64 //指定时间窗口内允许的最大请求数
	intervalNano int64 //指定的时间窗口
	unixNano     int64 //unix时间戳, 单位为纳秒
}

//初始化限流计数器
func NewCounterLimit(interval time.Duration, limit int64) *CounterLimit {
	return &CounterLimit{
		counter:      0,
		limit:        limit,
		intervalNano: int64(interval),
		unixNano:     time.Now().UnixNano(),
	}
}

//是否限流
func (c *CounterLimit) Allow() bool {
	now := time.Now().UnixNano()

	//1. 如果时间超过了时间窗口, 就重新计数
	if now-c.unixNano > c.intervalNano {
		//原子操作, 不需要加锁
		atomic.StoreInt64(&c.counter, 0)
		atomic.StoreInt64(&c.unixNano, now)
		return true
	}

	//2. 如果没超过时间窗口, 判断是否超过了最大请求数
	atomic.AddInt64(&c.counter, 1)
	return c.counter < c.limit
}

func main() {
	limit := NewCounterLimit(time.Second, 100)
	m := make(map[int]bool)
	for i := 0; i < 3000; i++ {
		// time.Sleep(time.Microsecond * 20)
		allow := limit.Allow()
		if allow {
			// fmt.Printf("i = %d is allowed\n", i)
			m[i] = true
		} else {
			// fmt.Printf("i = %d is not allowed\n", i)
			m[i] = false
		}
	}

	for i := 0; i < 3000; i++ {
		fmt.Printf("i = %d allow = %v\n", i, m[i])
	}
}
