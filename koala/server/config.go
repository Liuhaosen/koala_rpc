package server

import (
	"fmt"
	"io/ioutil"
	"modtest/gostudy/lesson2/ibinarytree/koala/util"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

var (
	koalaConf = &KoalaConf{
		Port: 8080,
		Prometheus: PrometheusConf{
			SwitchOn: true,
			Port:     8081,
		},
		ServiceName: "koala_server",
		Register: RegisterConf{
			SwitchOn: false,
		},
		Log: LogConf{
			Level: "debug",
			Dir:   "./logs/",
		},
		Limit: LimitConf{
			SwitchOn: true,
			QPSLimit: 50000,
		},
	}
)

type KoalaConf struct {
	Port        int            `yaml:"port"`
	Prometheus  PrometheusConf `yaml:"prometheus"`
	ServiceName string         `yaml:"service_name"`
	Register    RegisterConf   `yaml:"register"`
	Log         LogConf        `yaml:"log"`
	Limit       LimitConf      `yaml:"limit"`
	Trace       TraceConf      `yaml:"trace"`

	//内部配置项
	ConfigDir  string `yaml:"-"`
	RootDir    string `yaml:"-"`
	ConfigFile string `yaml:"-"`
}

type TraceConf struct {
	SwitchOn   bool    `yaml:"switch_on"`
	ReportAddr string  `yaml:"report_addr"`
	SampleType string  `yaml:"sample_type"`
	SampleRate float64 `yaml:"sample_rate"`
}

type LimitConf struct {
	QPSLimit int  `yaml:"qps"`
	SwitchOn bool `yaml:"switch_on"`
}

type PrometheusConf struct {
	SwitchOn bool `yaml:"switch_on"`
	Port     int  `yaml:"port"`
}

type RegisterConf struct {
	SwitchOn     bool          `yaml:"switch_on"`
	RegisterPath string        `yaml:"register_path"`
	Timeout      time.Duration `yaml:"timeout"`
	HeartBeat    int64         `yaml:"heart_beat"`
	RegisterName string        `yaml:"register_name"`
	RegisterAddr string        `yaml:"register_addr"`
}

type LogConf struct {
	Level    string `yaml:"level"`
	Dir      string `yaml:"path"`
	ChanSize int    `yaml:"chan_size"`
}

func initDir(serviceName string) (err error) {

	exeFilePath, err := filepath.Abs(os.Args[0])
	if err != nil {
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

	//d:/code/goproject/modtest/xxx
	koalaConf.RootDir = path.Join(strings.ToLower(exeFilePath[0:lastIndex]), "..")
	koalaConf.ConfigDir = path.Join(koalaConf.RootDir, "./conf/", util.GetEnv())
	koalaConf.ConfigFile = path.Join(koalaConf.ConfigDir, fmt.Sprintf("%s.yaml", serviceName))
	return
}

func InitConfig(serviceName string) (err error) {
	err = initDir(serviceName)
	if err != nil {
		return
	}

	data, err := ioutil.ReadFile(koalaConf.ConfigFile)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(data, &koalaConf)
	if err != nil {
		return
	}
	fmt.Printf("初始化koala配置项成功, conf: %#v\n", koalaConf)
	return
}

func GetConfigDir() string {
	return koalaConf.ConfigDir
}

func GetRootDir() string {
	return koalaConf.RootDir
}

func GetServerPort() int {
	return koalaConf.Port
}

func GetConf() *KoalaConf {
	return koalaConf
}
