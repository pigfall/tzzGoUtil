package main

import (
	"flag"
	"os/signal"
	"syscall"
	"sync"
	"context"
	"os"
	"net/http"
	"net"
"github.com/Peanuttown/tzzGoUtil/log"		
"github.com/Peanuttown/tzzGoUtil/async"		
	"github.com/Peanuttown/tzzGoUtil/net/water_tun_wrap/cmd_test/utils"	
)


func main() {
	var logger log.LoggerLite= log.NewLogger()
	// {  parse params
	var listenAt string
	flag.StringVar(&listenAt,"listenAt","","listen at")
	flag.Parse()
	if  len(listenAt) == 0{
		listenAt = ":10101"
	}
	// }

	ctx,cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()
	var cancelFuncs = []func(){func(){cancelCtx()}}
	wg := sync.WaitGroup{}

	// {　create tun interface
	tun,err := createTun("10.10.0.1/16",logger)
	if err != nil{
		logger.Error(err)
		os.Exit(1)
	}
  // }
	// { // listen port
	listener,err := net.Listen("tcp",listenAt)
	if err != nil{
		logger.Error(err)
		os.Exit(1)
	}
	cancelFuncs = append(cancelFuncs,func(){listener.Close()})
	// }

	var msgReadFromTun = make(chan utils.MsgReadFromTun,100)
	var msgWillWriteToTun = make(chan utils.MsgWillWriteToTun,100)
	async.AsyncDo(
		ctx,
		&wg,
		func(ctx context.Context){
			utils.LoopRWTun(ctx,logger,tun,msgReadFromTun,msgWillWriteToTun)
		},
	)

	// {　start server
	async.AsyncDo(
		ctx,
		&wg,
		func(ctx context.Context){
			svr := http.NewServeMux()
			svr.Handle("/",NewConnHandler(logger))
			err :=http.Serve(listener,svr)
			if err != nil{
				logger.Error(err)
			}
		},
	)
	// }
	wg.Wait()
}



func elegantWaitAndQuit(ctx context.Context,logger log.LoggerLite,cancel func(),wg *sync.WaitGroup,taskOver chan struct{}){
	sig := make(chan os.Signal ,1)
	signal.Notify(sig,syscall.SIGTERM,syscall.SIGINT)
	select{
	case <-taskOver:
		logger.Info("All async tasks over")
		break
	case <-ctx.Done():
		logger.Info("Context Done")
		break
	case <-sig:
		logger.Info("Rcv quit signal")
	}
	logger.Infof("Cancel context, app will quit")
	cancel()
	logger.Infof("Wait all task done")
	wg.Wait()
	logger.Info("Quit app")
}
