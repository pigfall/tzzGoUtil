package rb

import(
	"time"
		"context"
		"github.com/streadway/amqp"
		"net/url"
		"github.com/Peanuttown/tzzGoUtil/log"
)



type Consumer struct{
	connCfg  *ConnCfg
	logger log.LoggerI
}

func NewConsumer(
	connCfg *ConnCfg,
	logger log.LoggerI,
)*Consumer{
	return &Consumer{
		connCfg:connCfg,
		logger:logger,
	}
}

type ConsumeReqRequired struct{
	QueueName string
	ConsumerName string
	AutoAck bool
	NoWait bool
}

func (this *Consumer) ConsumeLoop (
	ctx context.Context,do func(ctx context.Context, msg amqp.Delivery),
){
	logger := this.logger
	for{
		select{
		case <-ctx.Done():
			logger.Info("Context Done : %v",ctx.Err())
			return 
		default:
			err := newConnToComsume(ctx,this.logger,this.connCfg,do)
			if err != nil{
				logger.Error(err)
			}
			logger.Info("ReConnect To Server")
		}
	}
}

func  newConnToComsume(
	ctx context.Context,
	logger log.LoggerI,
	cfg *ConnCfg,
	do func(ctx context.Context,msg amqp.Delivery),
)(err error){
	urlToConnect := cfg.ToUrl()
	conn,err := amqp.Dial(urlToConnect)
	if err != nil{
		return err
	}
	defer conn.Close()
	ch,err := conn.Channel()
	if err != nil{
		return err
	}
	consumeCh,err := ch.Consume()
	if err != nil{
		return err
	}
	for msg := range consumeCh{
		do(ctx,msg)
	}

	return nil
}
