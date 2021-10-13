package wintun

import(
	"golang.zx2c4.com/wireguard/tun"	
	"github.com/pigfall/tzzGoUtil/net"
)


type tunDev struct{
	tun.Device
}

func NewTun(ifceName string,mtu int)(net.TunIfce,error){
	dev,err := tun.NewTun(ifceName,mtu)
	if err != nil{
		return nil,err
	}

	return &tunDev {
		Device:dev,
	},nil
}



func (this *tunDev) SetIp(ip ...string)error{
	panic("TODO")
}

func (this *tunDev) Read(buf []byte)(int,error){
	return this.Device.Read(buf,0)
}

func (this *tunDev) Write(p []byte)(int,error){
	return this.Device.Write(p,0)
}
