package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	cpuTemp = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_temperature_celsius",
		Help: "current temperature of the CPU",
	})
	hdFailures = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "hd_errors_total", //这里是采样打点的名字
			Help: "Number of hard-disk errors.",
		},
		[]string{"device", "service"}, //两个标签, 展示的时候可以通过标签做筛选
	)
)

func init() {
	prometheus.MustRegister(cpuTemp)
	prometheus.MustRegister(hdFailures)
	//注册到普罗米修斯的库里
}

func main() {
	go func() {
		for {
			val := rand.Float64() * 100
			cpuTemp.Set(val)
			hdFailures.With(prometheus.Labels{
				"device":  "/dev/sda",
				"service": "hello.world", //可以通过标签知道哪个磁盘出了问题
			}).Inc()
			time.Sleep(time.Second)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
