package main

import(
	"fmt"
	"github.com/pigfall/tzzGoUtil/net"	
	"github.com/pigfall/tzzGoUtil/log"	
	water_wrap "github.com/pigfall/tzzGoUtil/net/water_tun_wrap"	
)


func createTun(ip string,logger log.LoggerLite)(net.TunIfce,error){
	logger.Info("Creating tun interface")
	tun,err := water_wrap.NewTun()
	if err != nil {
		logger.Error(fmt.Errorf("Create tun interface failed: %w",err))
		return nil,err
	}
	err = tun.SetIp(ip)
	if err != nil {
		logger.Error("Set ip { %s } to tun interface { %s } failed: %v",ip,tun.Name(),err)
		return nil,err
	}
	logger.Infof("Created tun interface { %s } and set ip { %s } to it",tun.Name(),ip)
	return tun,nil
}
