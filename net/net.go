package net

import (
	"fmt"
	goNet "net"
	"os"
)

func GetIp() string{

	addrs, err := goNet.InterfaceAddrs()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, address := range addrs {

		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*goNet.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				fmt.Println(ipnet.IP.String())
                return ipnet.IP.String()
			}
		}
	}
    return ""
}
