package main

type Option struct {
	Proto3Filename string //proto文件名
	Output         string //指定代码生成的路径
	GenClientCode  bool   //客户端代码
	GenServerCode  bool   //生成服务端代码
	Prefix         string //代码生成路径
}
