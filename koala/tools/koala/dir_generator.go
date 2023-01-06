package main

// 生成项目的目录
import (
	"fmt"
	"os"
	"path"
)

// var (
// 	AllDirList []string = []string{
// 		"controller",
// 		"idl",
// 		"main",
// 		"scripts",
// 		"conf",
// 		"app/router",
// 		"app/config",
// 		"model",
// 		"generate",
// 		"router",
// 	}
// )

//目录生成
type DirGenerator struct {
	dirList []string //目录
}

func (d *DirGenerator) Run(opt *Option) (err error) {
	for _, dir := range d.dirList {
		//拼接路径
		fullDir := path.Join(opt.Output, dir)
		err = os.MkdirAll(fullDir, 0755)
		if err != nil {
			fmt.Printf("make dir %s failed, err: %v\n", dir, err)
			return
		}
	}
	return
}

func init() {
	// dir := &DirGenerator{
	// 	dirList: AllDirList,
	// }
	// Register("dir generator", dir)
}
