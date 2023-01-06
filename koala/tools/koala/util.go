package main

import "html/template"

//模板函数映射.
var templateFuncMap = template.FuncMap{
	"Capitalize": Capitalize, //key是模板的函数名, value是这里实际的函数名
}

//模板函数映射的函数
//使用方法: 在模板中使用{{Capitalize .Package.Name}}
func Capitalize(str string) string {
	var output string
	chars := []rune(str)
	for i := 0; i < len(chars); i++ {
		if i == 0 {
			if chars[i] < 'a' || chars[i] > 'z' {
				return output
			}
			chars[i] -= 32
			output += string(chars[i])
		} else {
			output += string(chars[i])
		}
	}
	return output
}
