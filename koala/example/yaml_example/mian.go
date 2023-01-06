package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Site  SiteConfig  `yaml:"site"`
	Nginx NginxConfig `yaml:"nginx"`
}

type SiteConfig struct {
	Port      int    `yaml:"port"`
	HttpsOn   bool   `yaml:"https_on"`
	Domain    string `yaml:"domain"`
	HttpsPort int    `yaml:"https_port"`
}

type NginxConfig struct {
	Port     int      `yaml:"port"`
	LogPath  string   `yaml:"log_path"`
	SiteName string   `yaml:"site_name"`
	SiteAddr string   `yaml:"site_addr"`
	Upstream []string `yaml:"upstream"`
}

func main() {
	fmt.Printf("os.args[0]= %s\n", os.Args[0])
	//1. 把配置文件读到内存中
	data, err := ioutil.ReadFile("./test.yaml")
	if err != nil {
		fmt.Println("read file failed, err:", err)
		return
	}

	//2. 创建结构体
	var config Config

	//3. 解析配置文件,把数据放到config结构体中
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		fmt.Println("unmarshal failed, err : ", err)
		return
	}

	// fmt.Printf("config: %#v\n", config)
	// fmt.Printf("site port: %d\n", config.Site.Port)
}
