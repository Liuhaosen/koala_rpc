package util

import (
	"fmt"
	"net"
	"strings"
	"sync/atomic"
)

var ipAuto atomic.Value //缓存ip

func GetLocalIP() (ip string, err error) {
	ip, ok := ipAuto.Load().(string) //查看缓存里是否有ip

	if ok && len(ip) > 0 {
		return
	}

	//如果缓存没有, 那就获取ip
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		err = fmt.Errorf("get local ip failed")
		return
	}
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	ip = strings.Split(localAddr.String(), ":")[0]
	return

}
