package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/emicklei/proto"
)

var AllDirList []string = []string{
	"controller",
	"idl",
	"main",
	"scripts",
	"conf/product",
	"conf/test",
	"app/router",
	"app/config",
	"model",
	"generate",
	"router",
}

var genMgr *GeneratorMgr = &GeneratorMgr{
	genClientMap: make(map[string]Generator),
	genServerMap: make(map[string]Generator),
	metaData:     &ServiceMetaData{},
}

//管理所有generator
type GeneratorMgr struct {
	genClientMap map[string]Generator
	genServerMap map[string]Generator
	metaData     *ServiceMetaData
}

//解析Service
func (g *GeneratorMgr) parseService(opt *Option) (err error) {
	//打开proto文件
	reader, err := os.Open(opt.Proto3Filename)
	if err != nil {
		fmt.Printf("open file %s failed, err : %v\n", opt.Proto3Filename, err)
		return
	}
	defer reader.Close()
	//解析proto文件
	parser := proto.NewParser(reader)
	definition, err := parser.Parse()
	if err != nil {
		fmt.Printf("解析%s失败, 错误:%v", opt.Proto3Filename, err)
		return
	}
	proto.Walk(definition,
		//解析到service信息, 就使用handleService回调函数
		proto.WithService(g.handleService),
		proto.WithMessage(g.handleMessage),
		proto.WithRPC(g.handleRPC),
		proto.WithPackage(g.handlePackage),
	)
	return
}

//创建目录
func (g *GeneratorMgr) createAllDir(opt *Option) (err error) {
	for _, dir := range AllDirList {
		//拼接路径
		fullDir := path.Join(opt.Output, dir)
		err = os.MkdirAll(fullDir, 0755)
		if err != nil {
			err = fmt.Errorf("make dir %s failed, err: %v", dir, err)
			return
		}
	}
	return
}

//处理output输出路径, 使用prefix字段代替.
//处理反斜杠
func (g *GeneratorMgr) initOutputDir(opt *Option) (err error) {
	//如果是在gopath下写的代码. 那就去掉gopath, 如果不是, 要查看当前文件的目录.
	// goPath := os.Getenv("GOPATH")
	cmd := exec.Command("go", "env", "GOMOD")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("grpc generator failed, err : %v\n", err)
		return
	}
	defer stdout.Close()
	err = cmd.Start()
	data, err := ioutil.ReadAll(stdout)
	goPath := string(data)
	goPath = strings.Replace(goPath, `\`, `/`, -1)
	lastindex := strings.LastIndex(goPath, "modtest")
	if lastindex < 0 {
		fmt.Printf("invalid exe path: %v", goPath)
		return
	}
	goPath = strings.ToLower(goPath[0:lastindex])

	//(1) 假如用户指定了prefix = github.com/ibinarytree/koala/example
	if len(opt.Prefix) > 0 {
		//output = $GOPATH/$prefix
		if strings.LastIndex(goPath, "modtest") < 0 {
			goPath = path.Join(goPath, "modtest")
		}
		opt.Output = path.Join(goPath, opt.Prefix)
		return
	}

	//(2) 没有指定, 使用当前路径作为包的路径以及代码输出目录
	exeFilePath := os.Args[0]
	exeFilePath, err = filepath.Abs(exeFilePath)
	if err != nil {
		fmt.Println("get absolute path failed, err:", err)
		return
	}

	if runtime.GOOS == "windows" {
		exeFilePath = strings.Replace(exeFilePath, "\\", "/", -1)
	}

	lastIndex := strings.LastIndex(exeFilePath, "/")
	if lastIndex < 0 {
		err = fmt.Errorf("invalid exe path: %v", exeFilePath)
		return
	}
	opt.Output = strings.ToLower(exeFilePath[0:lastIndex])
	// srcPath := path.Join(goPath, "src/")
	//去掉gopath前缀
	opt.Prefix = strings.Replace(opt.Output, goPath, "", -1)
	fmt.Printf("opt output: %s, prefix: %s\n", opt.Output, opt.Prefix)
	return
}

//运行服务
func (g *GeneratorMgr) Run(opt *Option) (err error) {
	err = g.initOutputDir(opt)
	if err != nil {
		return
	}
	//解析服务数据
	err = g.parseService(opt)
	if err != nil {
		fmt.Println("meta data parse failed, err :", err)
		return
	}

	g.metaData.Prefix = opt.Prefix
	if opt.GenClientCode {
		//遍历所有已注册的生成器, 运行.
		for _, gen := range g.genClientMap {
			err = gen.Run(opt, g.metaData)
			if err != nil {
				return
			}
		}
		return
	}

	if opt.GenServerCode {
		//先创建好目录
		err = g.createAllDir(opt)
		if err != nil {
			fmt.Println(err)
			return
		}

		//遍历所有已注册的生成器, 运行.
		for _, gen := range g.genServerMap {
			err = gen.Run(opt, g.metaData)
			if err != nil {
				return
			}
		}
		return
	}

	return
}

//注册客户端
func RegisterClientGenerator(name string, gen Generator) (err error) {
	_, ok := genMgr.genClientMap[name]
	if ok {
		err = fmt.Errorf("generator %s is esists", name)
	}
	genMgr.genClientMap[name] = gen
	return
}

//注册服务端
func RegisterServerGenerator(name string, gen Generator) (err error) {
	_, ok := genMgr.genServerMap[name]
	if ok {
		err = fmt.Errorf("generator %s is esists", name)
	}
	genMgr.genServerMap[name] = gen
	return
}

func (g *GeneratorMgr) handleService(s *proto.Service) {
	g.metaData.Service = s
}

func (g *GeneratorMgr) handleMessage(m *proto.Message) {
	g.metaData.Messages = append(g.metaData.Messages, m)
}

func (g *GeneratorMgr) handleRPC(r *proto.RPC) {
	g.metaData.Rpc = append(g.metaData.Rpc, r)
}

func (g *GeneratorMgr) handlePackage(p *proto.Package) {
	g.metaData.Package = p
}
