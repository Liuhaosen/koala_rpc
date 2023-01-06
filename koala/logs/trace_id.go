package logs

import (
	"context"
	"fmt"
	"math/rand"
	"modtest/gostudy/lesson2/ibinarytree/koala/util"
	"time"
)

var (
	MaxTraceId = 100000000
)

type traceIdKey struct{}

func init() {
	rand.Seed(time.Now().UnixNano())
}

//获取trace_id, 打印日志时调用.
func GetTraceId(ctx context.Context) (traceId string) {
	traceId, ok := ctx.Value(traceIdKey{}).(string)
	if !ok {
		traceId = "-"
	}
	return
}

//生成traceId, 时间戳 + 随机数
func GenTraceId() (traceId string) {
	now := time.Now()
	traceId = fmt.Sprintf("%04d%02d%02d%02d%02d%02d%08d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), rand.Int31n(int32(MaxTraceId)))
	util.TraceID = traceId
	return
}

//设置traceId, 在程序入口使用
func WithTraceId(ctx context.Context, traceId string) context.Context {
	return context.WithValue(ctx, traceIdKey{}, traceId)
}
