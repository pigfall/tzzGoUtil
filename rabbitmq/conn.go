package rb

import(
	"fmt"
		"github.com/streadway/amqp"
		"net/url"
)

type Client struct{
	*amqp.Connection
}


type NewClientOption func(*amqp.Config)

func NewClientOptionVHost(vhost string)func(*amqp.Config){
	return func(cfg *amqp.Config){
		cfg.Vhost = vhost
	}
}

func NewClient(connCfg *ConnCfg,options ...NewClientOption)(*Client,error){
	cfg := amqp.Config{}
	for _,op := range options{
		op(&cfg)
	}
	conn,err := amqp.DialConfig(connCfg.ToUrl(),cfg)
	if err != nil{
		return nil,err
	}
	return &Client{Connection:conn},nil
}

type ConnCfg struct{
	InSecureConn bool
	User string
	Password string
	Host string
	Port string
}

type Conn struct{
	*amqp.Connection
}

func NewConnCfg(insecure bool,user, password, host, port string)*ConnCfg{
	return &ConnCfg{
		InSecureConn:insecure,
		User:user,
		Password:password,
		Host:host,
		Port:port,
	}
}

func (this *ConnCfg) ToUrl()string{
	u := url.URL{}
	u.Host = fmt.Sprintf("%s:%s",this.Host,this.Port)
	u.User = url.UserPassword(this.User,this.Password)
	if this.InSecureConn {
		u.Scheme = "amqp"
	}else{
		u.Scheme = "amqps"
	}
	return u.String()
}

func ValidConnLoopDo(cfg *ConnCfg,do func(conn *Conn)(stop bool,connectionClosed bool)){
}


