package main

import (
	"fmt"
	"html/template"
	"modtest/gostudy/lesson2/ibinarytree/koala/util"
	"os"
	"path"
)

type MainGenerator struct {
}

func (m *MainGenerator) Run(opt *Option, metaData *ServiceMetaData) (err error) {
	filename := path.Join("./", opt.Output, "main/main.go")

	exist := util.IsFileExist(filename)
	if exist {
		return
	}

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		fmt.Println("open file ", filename, " failed, err : ", err)
		return
	}

	defer file.Close()
	//渲染模板
	err = m.render(file, main_template, metaData)
	if err != nil {
		fmt.Println("render template failed, err :", err)
		return
	}
	return
}

func (m *MainGenerator) render(file *os.File, data string, metaData *ServiceMetaData) (err error) {
	temp := template.New("main").Funcs(templateFuncMap)
	temp, err = temp.Parse(data)
	if err != nil {
		return
	}

	temp.Execute(file, metaData)
	return
}

func init() {
	m := &MainGenerator{}
	RegisterServerGenerator("main generator", m)
}
