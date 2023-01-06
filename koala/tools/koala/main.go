package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	var opt Option

	app := cli.NewApp()
	app.Flags = []cli.Flag{
		//idl文件名
		cli.StringFlag{
			Name:        "f",
			Value:       "./test.proto",
			Usage:       "idl filename",
			Destination: &opt.Proto3Filename,
		},
		//输出目录
		cli.StringFlag{
			Name:        "o",
			Value:       "./output/",
			Usage:       "output directory",
			Destination: &opt.Output,
		},
		//是否生成grpc 客户端代码
		cli.BoolFlag{
			Name:        "c",
			Usage:       "generate grpc client code",
			Destination: &opt.GenClientCode,
		},
		//是否生成grpc 服务端代码
		cli.BoolFlag{
			Name:        "s",
			Usage:       "generate grpc server code",
			Destination: &opt.GenServerCode,
		},
		//导入包的路径前缀
		cli.StringFlag{
			Name:        "p",
			Value:       "",
			Usage:       "prefix of package",
			Destination: &opt.Prefix,
		},
	}

	app.Action = func(c *cli.Context) error {
		//命令行程序代码入口
		err := genMgr.Run(&opt)
		if err != nil {
			fmt.Printf("code generator failed, err : %v\n", err)
			return err
		}
		fmt.Println("code generate succ")
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
