package main

import (
	"flag"
	"log"
	"math"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	addr              = flag.String("listen-address", ":8080", "the address to listen for http request")
	uniformDomain     = flag.Float64("uniform.domain", 0.0002, "the domain for the uniform distribution") //均匀分布的域
	normDomain        = flag.Float64("normal.domain", 0.0002, "the domain for the normal distribution")   //正态分布的域
	normMean          = flag.Float64("normal.mean", 0.00001, "the mean for the normal distribution")      //正态分布的平均值
	oscillationPeriod = flag.Duration("oscillation-period", 10*time.Minute, "The duration of the rate oscillation period")
)

var (
	rpcDurations = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "rpc_durations_seconds",
			Help:       "RPC latency distributions.",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001}, //50%的请求是0.05秒解决.
		},
		[]string{"service"},
	)
)

func init() {
	prometheus.MustRegister(rpcDurations)
	prometheus.MustRegister(prometheus.NewBuildInfoCollector())
}

func main() {
	flag.Parse()
	start := time.Now()
	oscillationFactor := func() float64 {
		return 2 + math.Sin(math.Sin(2*math.Pi*float64(time.Since(start))/float64(*oscillationPeriod)))
	}

	go func() {
		for {
			v := rand.Float64() * *uniformDomain
			rpcDurations.WithLabelValues("uniform").Observe(v)
			time.Sleep(time.Duration(100*oscillationFactor()) * time.Millisecond)
			//当前请求耗时传入ovserve, 统计出来.
		}
	}()

	// go func() {
	// 	for {
	// 		v := rand.ExpFloat64() / 1e6
	// 		rpcDurations.WithLabelValues("exponential", "303").Observe(v)
	// 		time.Sleep(time.Duration(50*oscillationFactor()) * time.Millisecond)
	// 	}
	// }()

	// Expose the registered metrics via HTTP.
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}
