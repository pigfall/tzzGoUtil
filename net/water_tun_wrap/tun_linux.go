package water_wrap

import(
	"github.com/pigfall/tzzGoUtil/net"
	"github.com/pigfall/tzzGoUtil/process"
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
func(this *tun)Name() (string,error){
	return this.ifce.Name(),nil
}

func (this *tun) SetIp(ip ...string)error{
	devName,err := this.Name()
	if err != nil{
		return err
	}
	_,errOut,err := process.ExeOutput(
		"ip",
		"addr",
		"add",
		ip[0],
		"dev",
		devName,
	)
	if err != nil{
		return fmt.Errorf("Set ip failed: %v, %v",err,errOut)
	}

	_,errOut,err  = process.ExeOutput(
		"ip",
		"link",
		"set",
		devName,
		"up",
	)
	if err != nil{
		return fmt.Errorf("Enable ifce up failed: %w, %v",err,errOut)
	}

	return nil
}
