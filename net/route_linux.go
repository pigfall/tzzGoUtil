package net

import(
	"fmt"
	"net"
	"github.com/pigfall/tzzGoUtil/process"
)



func AddRoute(target net.IP,devName string,via net.IP)error{
	return addRoute(target.String(),devName,via)
}

func AddRouteIpNet(target *IpWithMask,devName string,via net.IP)error{
	return addRoute(target.FormatAsIpSlashMask(),devName,via)
}

func addRoute(target string,devName string,via net.IP)error{
	var errOut string
	var err error
	if via != nil{
		_,errOut,err = process.ExeOutput("ip","route","add",target,"via",via.String(),"dev",devName)
	}else{
		_,errOut,err = process.ExeOutput("ip","route","add",target,"dev",devName)
	}

	if err != nil{
		return fmt.Errorf("%s, %w",errOut,err)
	}

	return err
}
