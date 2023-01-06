package logs

import (
	"bytes"
	"fmt"
	"runtime"
	"time"
)

type LogData struct {
	curTime     time.Time
	message     string
	timeStr     string
	level       LogLevel
	filename    string
	lineNo      int
	traceId     string
	serviceName string
	fields      *LogField //通过AddField增加字段
}

func writeField(buffer *bytes.Buffer, field, sep string) {
	buffer.WriteString(field) //写入字段
	buffer.WriteString(sep)   //写入分隔符
}

func writeKv(buffer *bytes.Buffer, key, val string) {
	buffer.WriteString(key)
	buffer.WriteString("=")
	buffer.WriteString(val)
}

//将LogData转换为字节流
func (l *LogData) Bytes() []byte {
	var buffer bytes.Buffer
	levelStr := getLevelText(l.level)
	//格式: 时间 日志级别 服务名 文件名: 行号 traceId
	writeField(&buffer, l.timeStr, SpaceSep)
	writeField(&buffer, levelStr, SpaceSep)
	writeField(&buffer, l.serviceName, SpaceSep)

	writeField(&buffer, l.filename, ColonSep)
	writeField(&buffer, fmt.Sprintf("%d", l.lineNo), SpaceSep)
	writeField(&buffer, l.traceId, SpaceSep)

	//如果是access访问日志, 把用户增加的field也都输出到缓冲区里
	if l.level == LogLevelAccess && l.fields != nil {
		for _, field := range l.fields.kvs {
			writeField(&buffer, fmt.Sprintf("%v = %v", field.key, field.val), SpaceSep)
		}
	}

	writeField(&buffer, l.message, LineSep)
	//将缓冲内容转换为字节数组
	return buffer.Bytes()
}

func GetLineInfo() (fileName string, lineNo int) {
	_, fileName, lineNo, _ = runtime.Caller(3)
	//例如, 先调Debug(), 然后是Debug()里的writeLog, 然后再调GetLineInfo. 所以Caller参数是3. 代表堆栈升路=3
	//就能获取调用Api的地方的文件名和行号.
	return
}
