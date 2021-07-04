package rb

import(
	"errors"
	"time"
	"context"
	"fmt"
		"github.com/streadway/amqp"
		"github.com/Peanuttown/tzzGoUtil/log"
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

func WaitExchangeIsNilThenDeleteIt(ctx context.Context,connCfg *ConnCfg,options []NewClientOption,exchangeName, queueName string,logger log.LoggerLite)error{
	for{
		select{
		case <-ctx.Done():
			return ctx.Err()
		default:
			err := waitExchangeIsNilThenDeleteIt(ctx,connCfg,options,exchangeName,queueName,logger)
			if err != nil{
				logger.Error(err)
				return err
			}
			return nil
		}
	}
}

func  waitExchangeIsNilThenDeleteIt(ctx context.Context,connCfg *ConnCfg,options []NewClientOption, exchangeName string,queueName string,logger log.LoggerLite)error{
	ctx,cancel := context.WithCancel(ctx)
	defer cancel()
	connCloseNotifer := make(chan *amqp.Error,2)
	cli,err := NewClient(connCfg,options...)
	if err != nil{
		return err
	}
	defer cli.Close()
	connCloseNotifer = cli.NotifyClose(connCloseNotifer)
	deleteSuccess := make(chan struct{},1)
	go func(){
		defer cancel()
		select{
			case <-connCloseNotifer:
				return
			case <-ctx.Done():
				return
		}
	}()
	go func(){ // delete queue
		ticker := time.NewTicker(time.Second*2)
		defer close(deleteSuccess)
		for {
			select{
			case <-ctx.Done():
				return
			case <-ticker.C:
				err = DoThenCloseChannelWithCli(
					ctx,
					cli,
					func(ctx context.Context,mqCh *amqp.Channel)error{
						// try to delete queue
						// check if exist
						logger.Debugf("尝试删除队列 %s",queueName)
						_,err := mqCh.QueueDeclarePassive(queueName,true,false,false,false,nil)
						if err != nil{
							var amqpErr *amqp.Error
							if errors.As(err,&amqpErr){
								if amqpErr.Code == amqp.NotFound{
									logger.Debug("队列 %s 原本就不存在",queueName)
									// queue not exist
									return nil
								}
							}
							return err
						}else{
							_,err = mqCh.QueueDelete(queueName,false,true,false)
							if err != nil{
								err = fmt.Errorf("Try delete quque  %s failed: %w",queueName,err)
								logger.Error(err)
								return err
							}
							logger.Debugf("删除队列 %s 成功",queueName)
							return nil
						}
					},
				)
				if err != nil{
					logger.Error(err)
					continue
				}
				err = DoThenCloseChannelWithCli(
					ctx,
					cli,
					func(ctx context.Context,mqCh *amqp.Channel)error{

						err = mqCh.QueueUnbind(queueName,"",exchangeName,nil)
						if  err != nil{
							logger.Errorf("Tru unbind queue failed: %v",err)
							// goon
						}else{
							logger.Debug("unbind queue %s success",queueName)
						}
						err = mqCh.ExchangeDelete(exchangeName,false,false)
						if err != nil{
							err = fmt.Errorf("Try delete exchange %s failed: %w",exchangeName,err)
							logger.Error(err)
							return err
						}
						logger.Debugf("删除 exchange %s 成功",exchangeName)
						// >
						return nil
					},
				)
				if err != nil{
					logger.Error(err)
					continue
				}
				deleteSuccess<-struct{}{}
				return
			}
		}
	}()
	select{
	case <-ctx.Done():
		return ctx.Err()
	case _,ok:=<-deleteSuccess:
		if ok{
			logger.Debug("删除 exchange,queue 成功")
			return nil
		}else{
			return fmt.Errorf("删除 exchange,queue 失败")
		}
	}
}


func DoThenCloseChannel(ctx context.Context,ch *amqp.Channel,do func(ctx context.Context,ch *amqp.Channel)error)error{
	defer ch.Close()
	return do(ctx,ch)
}

func DoThenCloseChannelWithCli(ctx context.Context,cli *Client,do func(ctx context.Context,ch *amqp.Channel)error)error{
	ch,err := cli.Channel()
	if err != nil{
		return err
	}
	defer ch.Close()
	return do(ctx,ch)
}
