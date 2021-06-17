package rb

import(
	"fmt"
		"github.com/streadway/amqp"
		"net/url"
)

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


