package logs

import "fmt"

const (
	Black Color = iota + 30
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

type Color uint8

//给字符串加上颜色
func (c Color) Add(s string) string {
	//加了特殊标识, 控制台可以输出颜色
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", uint8(c), s)
}
