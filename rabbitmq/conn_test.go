package rb

import(
	"testing"
)

func TestNewClient(t *testing.T){
	cfg := NewConnCfg(
		true,
		"MjphbXFwLWNuLTJyNDI4YmlpYzAwODpMVEFJNXQ3WHZkYWoxUFA4aE1RVFRqemI=",
		"MjlBODdERDIzMURBODAwNUEyMzQyQkQzOThGMEJGOUE5M0Y1QTdBRDoxNjI0NDI3MjE3MTQ1",
		"amqp-cn-2r428biic008.mq-amqp.cn-shanghai-867405-a.aliyuncs.com",
		"5672",
	)
	t.Log(cfg.ToUrl())
//NewClientOptionVHost("sweetlucky")
	_,err := NewClient(cfg)
	if err != nil{
		t.Fatal(err)
	}
}
