package logs

import (
	"os"
)

//控制台输出器
type ConsoleOutputer struct {
}

//控制台输出器构造函数
func NewConsoleOutputer() (log Outputer) {
	log = &ConsoleOutputer{}
	return
}

//输出到控制台
func (c *ConsoleOutputer) Write(data *LogData) {
	color := getLevelColor(data.level) //根据级别拿到日志对应颜色

	text := color.Add(string(data.Bytes()))
	os.Stdout.Write([]byte(text))

}

func (c *ConsoleOutputer) Close() {

}
