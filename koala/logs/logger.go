package logs

import (
	"context"
	"fmt"
	"path"
	"sync"
	"time"
)

var (
	defaultOutputer               = NewConsoleOutputer()
	lm                            = &LoggerMgr{}
	initOnce           *sync.Once = &sync.Once{} //解决多个线程调用的问题
	defaultServiceName            = "default"
)

type LoggerMgr struct {
	outputers   []Outputer
	chanSize    int
	level       LogLevel
	logDataChan chan *LogData
	serviceName string
	wg          sync.WaitGroup
	isClosed    bool
}

func initLogger(level LogLevel, chanSize int, serviceName string) {
	//initOnce.Do : 无论有多少线程去调用, 实际只用一个线程执行, 保证了线程安全
	if len(serviceName) == 0 {
		serviceName = defaultServiceName
	}

	initOnce.Do(func() {
		lm = &LoggerMgr{
			chanSize:    chanSize,
			level:       level,
			serviceName: serviceName,
			logDataChan: make(chan *LogData, chanSize),
		}
		lm.wg.Add(1)

		go lm.run()
	})
}

//对外暴露初始化日志库
func InitLogger(level LogLevel, chanSize int, serviceName string) {
	if chanSize <= 0 {
		chanSize = DefaultLogChanSize
	}
	initLogger(level, chanSize, serviceName)
}

//设置日志级别
func SetLevel(level LogLevel) {
	lm.level = level
}

//增加输出器
func AddOutputer(outputer Outputer) {
	if lm == nil {
		//如果没有初始化就先初始化日志库
		initLogger(LogLevelDebug, DefaultLogChanSize, defaultServiceName)
	}

	lm.outputers = append(lm.outputers, outputer)
	// return
}

//api start
func Debug(ctx context.Context, format string, args ...interface{}) {
	writeLog(ctx, LogLevelDebug, format, args...)
}

func Trace(ctx context.Context, format string, args ...interface{}) {
	writeLog(ctx, LogLevelTrace, format, args...)
}

func Access(ctx context.Context, format string, args ...interface{}) {
	writeLog(ctx, LogLevelAccess, format, args...)
}

func Info(ctx context.Context, format string, args ...interface{}) {
	writeLog(ctx, LogLevelInfo, format, args...)
}

func Warn(ctx context.Context, format string, args ...interface{}) {
	writeLog(ctx, LogLevelWarn, format, args...)
}

func Error(ctx context.Context, format string, args ...interface{}) {
	writeLog(ctx, LogLevelError, format, args...)
}

//api end

func Stop() {
	if !lm.isClosed {
		if lm.logDataChan != nil {
			close(lm.logDataChan)
		}

		lm.isClosed = true
	}

	lm.wg.Wait()

	for _, outputer := range lm.outputers {
		outputer.Close()
	}

	//重新初始化
	initOnce = &sync.Once{}
	lm = nil
}

func (lm *LoggerMgr) run() {

	for data := range lm.logDataChan {

		//1. 如果没有定义输出器, 那么默认输出到控制台

		if len(lm.outputers) == 0 {
			defaultOutputer.Write(data)
			continue
		}

		//2. 否则遍历所有输出器, 通过输出器输出
		for _, outputer := range lm.outputers {
			outputer.Write(data)
		}
	}
	lm.wg.Done()
}

//日志写入
func writeLog(ctx context.Context, level LogLevel, format string, args ...interface{}) {
	if lm == nil {
		initLogger(LogLevelDebug, DefaultLogChanSize, defaultServiceName)
	}
	now := time.Now()
	nowStr := now.Format("2006-01-02 15:04:05.999")
	fileName, lineNo := GetLineInfo()
	fileName = path.Base(fileName)
	msg := fmt.Sprintf(format, args...)

	logData := &LogData{
		message:     msg,
		curTime:     now,
		timeStr:     nowStr,
		level:       level,
		filename:    fileName,
		lineNo:      lineNo,
		traceId:     GetTraceId(ctx),
		serviceName: lm.serviceName,
	}

	//如果是access访问日志, 需要把所有field拉取出来
	//会把所有字段一起合并输出
	if level == LogLevelAccess {
		fields := getFields(ctx)
		if fields != nil {
			logData.fields = fields
		}
	}
	select {
	case lm.logDataChan <- logData:
	default:
		return
	}
}
