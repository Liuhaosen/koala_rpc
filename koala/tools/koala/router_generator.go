package main

import (
	"fmt"
	"html/template"
	"os"
	"path"
)

type RouterGenerator struct {
}

func (r *RouterGenerator) Run(opt *Option, metaData *ServiceMetaData) (err error) {
	filename := path.Join("./", opt.Output, "router/router.go")
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		fmt.Printf("open file:%s failed, err :%v\n", filename, err)
		return
	}

	defer file.Close()

	err = r.render(file, router_template, metaData)
	if err != nil {
		fmt.Printf("render failed, err : %v\n", err)
		return
	}
	return
}

func (r *RouterGenerator) render(file *os.File, data string, metaData *ServiceMetaData) (err error) {
	t := template.New("main")
	t, err = t.Parse(data)
	if err != nil {
		fmt.Println("parse template failed, err :", err)
		return
	}
	// fmt.Printf("metaData = %#v\n", metaData)
	err = t.Execute(file, metaData)
	return
}

func init() {
	router := &RouterGenerator{}
	RegisterServerGenerator("router generator", router)
}
