package utils

import (
	"github.com/Peanuttown/tzzGoUtil/net"
	"sync"
	"github.com/Peanuttown/tzzGoUtil/log"
	"github.com/Peanuttown/tzzGoUtil/async"
	"context"
)

type MsgReadFromTun []byte

type MsgWillWriteToTun []byte


// < loop read write for tun interface
func LoopRWTun(ctx context.Context,logger log.LoggerLite,tun net.TunIfce,msgReadFromTun chan MsgReadFromTun,msgWillWriteToTun chan MsgWillWriteToTun)error{
	wg := sync.WaitGroup{}
	// < { read from tun then write to readFromTun channel
	async.AsyncDo(
		ctx,
		&wg,
		func(ctx context.Context){
			err := loopReadTun(ctx,logger,tun,msgReadFromTun)
			if err != nil{
				logger.Error(err)
			}
		},
	)
	// > }

	// < { when rcv msg from channl: msgWillWriteToTun , write to tun
	async.AsyncDo(
		ctx,
		&wg,
		func(ctx context.Context){
			err := loopWriteToTun(ctx,logger,tun,msgWillWriteToTun)
			if err != nil{
				logger.Error(err)
			}
		},
	)
	// > }
	wg.Wait()
	logger.Info("Loop for tun rw quit")
	return nil
}

func loopReadTun(ctx context.Context,logger log.LoggerLite,tun net.TunIfce, msgReadFromTun chan MsgReadFromTun)error{
	var buf = make([]byte,1024*5)
	for {
		num,err := tun.Read(buf)
		if err != nil{
			logger.Error("Failed to read from tun")
			return err
		}
		bufToSend := make([]byte,num)
		copy(bufToSend,buf[:num])
		logger.Debug("Will send msg to readFromTunChannel ")
		msgReadFromTun<-bufToSend
		logger.Debug("Ok to send msg to readFromTunChannel ")
	}
}

func loopWriteToTun(ctx context.Context,logger log.LoggerLite,tun net.TunIfce,msgWillWriteToTun chan MsgWillWriteToTun)error{
	for {
		select{
		case <-ctx.Done():
			return ctx.Err()
		case msg :=<-msgWillWriteToTun:
			logger.Debug("Read one msg from WillWriteToTunChannel, will wirte to tun")
			_,err := tun.Write(msg)
			if err != nil{
				logger.Error(err)
			}else{
				logger.Info("Ok write bytes to tun")
			}
		}
	}
}
