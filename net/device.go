package net

import (
	"fmt"
	"net"
	"github.com/pigfall/tzzGoUtil/process"
	"github.com/pkg/errors"
)

type VIRTUAL_DEV int

const (
	TUN VIRTUAL_DEV = iota + 1
	TAP
)

type device struct {
	ifce net.Interface
}

func ListDevices()([]DeviceIfce,error){
	ifces,err := net.Interfaces()
	if err != nil{
		return nil,err
	}
	devs :=make([]DeviceIfce,0,len(ifces))
	for _,ifce := range ifces {
		devs = append(devs,&device{ifce:ifce})
	}
	return devs,nil
}


func (this *device) Addrs()([]IpWithMask,error){
	addrs,err := this.ifce.Addrs()
	if err != nil{
		return nil,err
	}
	addrsToRet := make([]IpWithMask,0,len(addrs))
	for _,addr := range addrs{
		cidrStr := addr.String()
		ip,_,err := net.ParseCIDR(cidrStr)
		if err != nil{
			return nil,fmt.Errorf("Parse cidr from %s failed: %v",cidrStr,err)
		}
		if IsIpv4(ip){
			ipWithMask,err := FromIpSlashMask(cidrStr)
			if err != nil{
				return nil,err
			}
			addrsToRet = append(addrsToRet,*ipWithMask)
		}
	}
	return addrsToRet,nil
}


func SetIp(devName string, ip IpWithMask) error {
	_, errOut, err := process.ExeOutput("ip", "addr", "add",ip.FormatAsIpSlashMask(), "dev", devName)
	if err != nil {
		return errors.Wrap(err, errOut)
	}
	return nil
}

func DevUp(devName string) error {
	_, errOut, err := process.ExeOutput("ip", "link", "set", devName, "up")
	if err != nil {
		return errors.Wrap(err, errOut)
	}
	return nil
}
