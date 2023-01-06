package logs

import (
	"context"
	"testing"
)

func TestFileLogger(t *testing.T) {
	outputer, err := NewFileOutputer("d:/tmp/logs/test.log")
	if err != nil {
		t.Errorf("init file outputer failed, err : %v", err)
		return
	}

	initLogger(LogLevelDebug, 10000, "account")
	AddOutputer(outputer) //把文件输出器, 加到日志库里

	Debug(context.Background(), "this is a debug test")
	Trace(context.Background(), "this is a trace test")
	Info(context.Background(), "this is a info test")
	Access(context.Background(), "this is a access test")
	Warn(context.Background(), "this is a warn test")
	Error(context.Background(), "this is a error test")
	Stop()
}

func TestConsoleLogger(t *testing.T) {
	//这里没有初始化, 默认输出到控制台
	initLogger(LogLevelDebug, 10000, "account")
	ctx := context.Background()
	ctx = WithFieldContext(ctx) //生成用来保存当前key-val的数据结构, 设置到context里
	ctx = WithTraceId(ctx, GenTraceId())
	AddField(ctx, "user_id", 833323423)
	AddField(ctx, "name", "kswss")

	Debug(context.TODO(), "this is a debug test")
	Trace(ctx, "this is a trace test")
	Info(ctx, "this is a info test")
	Access(ctx, "this is a access test")
	Warn(ctx, "this is a warn test")
	Error(ctx, "this is a error test")
	Stop()
}
