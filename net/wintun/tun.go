package wintun

import(
	"fmt"
	"golang.zx2c4.com/wireguard/tun"	
	"github.com/pigfall/tzzGoUtil/net"
	"github.com/pigfall/tzzGoUtil/process"
)


type tunDev struct{
	tun.Device
}

func NewTun(ifceName string,mtu int)(net.TunIfce,error) {
	dev,err := tun.NewTun(ifceName,mtu)
	if err != nil{
		return nil,err
	}

	return &tunDev {
		Device:dev,
	},nil
}



func (this *tunDev) SetIp(ip ...string)error{
	ipNet,err := net.FromIpSlashMask(ip[0])
	if err != nil{
		return err
	}
	devName,err := this.Name()
	if err != nil {
		return err
	}
	out,errOut,err := process.ExeOutput("netsh","interface","ip","set","address",devName,"static",ipNet.Ip.String(),net.MaskFormatTo255(ipNet.Mask))
	if err != nil{
		return fmt.Errorf("%w, %v, %v",err,errOut,out)
	}
	return nil
}

func (this *tunDev) Read(buf []byte)(int,error){
	return this.Device.Read(buf,0)
}

func (this *tunDev) Write(p []byte)(int,error){
	return this.Device.Write(p,0)
}
