package rb

import(
		"context"
		"time"
		"fmt"
		"github.com/streadway/amqp"
		"github.com/pigfall/tzzGoUtil/log"
		"github.com/pigfall/tzzGoUtil/uuid"
)




type Consumer struct{
	connCfg  *ConnCfg
	logger log.LoggerLite
}

func NewConsumer(
	connCfg *ConnCfg,
	logger log.LoggerLite,
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
	ctx context.Context,do func(ctx context.Context, msg *amqp.Delivery) error,
	exchangeId string,
	queueId string,
	newClientOptions []NewClientOption,
	onNewConn []func(ctx context.Context,ch *amqp.Channel)error,
){
	logger := this.logger
	tickerToRetry := time.NewTicker(time.Millisecond)
	for{
		select{
		case <-ctx.Done():
			logger.Info("Context Done : %v",ctx.Err())
			return 
		default:
			err := newConnToComsume(ctx,this.logger,this.connCfg,do,exchangeId,queueId,newClientOptions,onNewConn)
			if err != nil{
				logger.Error(err)
			}
			tickerToRetry.Reset(time.Second*2)
			logger.Info("ReConnect To Server")
		}
	}
}

func  newConnToComsume (
	ctx context.Context,
	logger log.LoggerLite,
	cfg *ConnCfg,
	do func(ctx context.Context,msg *amqp.Delivery)error,
	exchangeId string,
	queueId string,
	newClientOptions []NewClientOption,
	onNewConn []func(ctx context.Context, ch *amqp.Channel)error,
)(err error){
	cli,err := NewClient(cfg,newClientOptions...)
	if err != nil{
		logger.Error(err)
		return err
	}
	defer cli.Close()
	channel,err := cli.Channel()
	if err != nil{
		logger.Error(err)
		return err
	}
	defer channel.Close()
	for _,v := range onNewConn{
		err := v(ctx,channel)
		if err != nil{
			logger.Error(err)
			return err
		}
	}
	err = channel.QueueBind(
		queueId,
		"",
		exchangeId,
		false,
		nil,
	)
	if err != nil{
		err = fmt.Errorf("绑定队列 %s, %s 失败: %w",exchangeId,queueId,err)
		logger.Error(err)
		return err
	}
	uuid,err := uuid.GenUUID()
	if err != nil{
		logger.Error(err)
		return err
	}
	msgChannel,err := channel.Consume(
		queueId,
		uuid,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		err = fmt.Errorf("消费队列 %s 失败: %w",queueId,err)
		logger.Error(err)
		return err
	}
	logger.Debug("loop to wait msg")
	for{
		select{
		case <-ctx.Done():
			return fmt.Errorf("Ctx done with err %w",ctx.Err())
		case msg,ok:=<-msgChannel:
			logger.Debug("Rcv Msg")
			if !ok{
				return fmt.Errorf("channel closed")
			}
			err := do(ctx,&msg)
			if err != nil{
				logger.Error(err)
				//err := msg.Acknowledger
				//if err != nil{
				//	logger.Error("Msg Reject")
				//}
				
			}else{
				err:= 	msg.Ack(false)
				if err != nil{
					logger.Error("Msg Ack err :%v",err)
				}
			}
		}
	}
}
