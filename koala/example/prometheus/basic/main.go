package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var addr = flag.String("listen-address", ":8080", "the address to listen on for http requests")

func main() {
	flag.Parse()
	http.Handle("/metrics", promhttp.Handler()) //监听/metrics地址. 然后执行promhttp里的handler方法, 会自动把go里的一些常用监控加上
	log.Fatal(http.ListenAndServe(*addr, nil))  //监听该地址, 仅供prometheus使用.
}
