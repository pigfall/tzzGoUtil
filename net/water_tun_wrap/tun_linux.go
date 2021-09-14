package water_wrap

import(
	"github.com/Peanuttown/tzzGoUtil/net"
	"github.com/Peanuttown/tzzGoUtil/process"
	"fmt"
	"github.com/songgao/water"
)

type tun struct {
	ifce *water.Interface

}

func NewTun() (net.TunIfce,error){
	ifce,err := water.New(
		water.Config{
			DeviceType:water.TUN,
		},
	)
	if err != nil{
		return nil,err
	}
	return &tun{
		ifce:ifce,
	},nil
}

func(this *tun)Write(p []byte) (n int, err error){
	return this.ifce.Write(p)
}
func(this *tun)Read(p []byte) (n int, err error){
	return this.ifce.Read(p)
}
func(this *tun)Close() error{
	return this.ifce.Close()
}
func(this *tun)Name() string{
	return this.ifce.Name()
}

func (this *tun) SetIp(ip ...string)error{
	_,errOut,err := process.ExeOutput(
		"ip",
		"addr",
		"add",
		ip[0],
		"dev",
		this.Name(),
	)
	if err != nil{
		return fmt.Errorf("Set ip failed: %v, %v",err,errOut)
	}

	_,errOut,err  = process.ExeOutput(
		"ip",
		"link",
		"set",
		this.Name(),
		"up",
	)
	if err != nil{
		return fmt.Errorf("Enable ifce up failed: %w, %v",err,errOut)
	}

	return nil
}
