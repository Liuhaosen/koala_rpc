package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/afex/hystrix-go/hystrix"
)

func main() {
	hystrix.ConfigureCommand("koala_rpc", hystrix.CommandConfig{
		Timeout:               1000,
		MaxConcurrentRequests: 1000,
		ErrorPercentThreshold: 20,
	})

	//开启一个for循环, 正常连接一会然后断网.
	//就会看到

	for {
		//参数1: 服务名,
		//参数2: 访问服务的代码,
		//参数3: 熔断触发或者报错后的操作. 也可以是降级预案
		err := hystrix.Do("get_baidu", func() error {
			//talk to other services
			_, err := http.Get("https://www.baidu.com/")
			if err != nil {
				fmt.Println("get error")
				return err
			}
			return nil
		}, func(err error) error {
			fmt.Println("get an error, handle it , err :", err)
			return err
		})
		if err == nil {
			fmt.Println("request success")
		}
		time.Sleep(time.Second) //调用go方法就是起了一个goroutine, 这里要sleep一下, 不然看不到效果
	}
}
