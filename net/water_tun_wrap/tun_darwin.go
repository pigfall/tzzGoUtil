package water_wrap

import(
	"log"
		"github.com/Peanuttown/tzzGoUtil/net"
		"github.com/songgao/water"
)

type tun struct{
	ifce *water.Interface
}


func NewTun()(net.TunIfce,error){
	log.Println("creating tun")
	ifce,err := water.New(
		water.Config{
			DeviceType: water.TUN,
		},
	)
	if err != nil{
		return nil, err
	}
	return &tun{
		ifce:ifce,
	},nil
}


func (this *tun) Write(p []byte) (n int, err error){
	return this.ifce.Write(p)
}

func (this *tun) 	Read(p []byte) (n int, err error){
	return this.ifce.Read(p)
}

func (this *tun) Close() error {
	return this.ifce.Close()
}

func (this *tun) Name() string{
	return this.ifce.Name()
}
