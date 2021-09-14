package main

import(
	"context"
	"syscall"
	"flag"
	"sync"
	ws "github.com/gorilla/websocket"
	"os/signal"
	"time"
	"github.com/Peanuttown/tzzGoUtil/log"
	"github.com/Peanuttown/tzzGoUtil/async"
	ctx_lib "github.com/Peanuttown/tzzGoUtil/ctx"
	"os"
	"fmt"
	water_wrap "github.com/Peanuttown/tzzGoUtil/net/water_tun_wrap"	
	"github.com/Peanuttown/tzzGoUtil/net/water_tun_wrap/cmd_test/utils"	
)


func exit(logger log.LoggerLite,exitCode int){
	logger.Info("App exit")
	os.Exit(exitCode)
}

func main() {
	var logger log.LoggerLite= log.NewLogger()
	var serverAddr string
	var tunIp string
	flag.StringVar(&serverAddr,"server address","","server address")
	flag.StringVar(&tunIp,"tun ip","172.168.1.1/16","tunIp")
	flag.Parse()
	if len(serverAddr) ==0 {
		fmt.Println("please input server address")
		os.Exit(1)
	}

	// < { connect to server
	var retryDelaySec int = 3
	var conn *ws.Conn
	var err error
	for{
		logger.Infof("Connecting to server %s",serverAddr)
		conn,_,err = ws.DefaultDialer.Dial(serverAddr,nil)
		if err != nil{
			logger.Errorf("Failed to connect to server %s, %v",serverAddr,err)
			logger.Info("Will retry after %d second",retryDelaySec)
			time.Sleep(time.Second* time.Duration(retryDelaySec))
			continue
		}
		break
	}
	defer conn.Close()
	// > }

	// < { create tun interface
	logger.Info("Creating tun interface")
	tun,err := water_wrap.NewTun()
	if err != nil {
		logger.Error(fmt.Errorf("Create tun interface failed: %w",err))
		exit(logger,1)
	}
	err = tun.SetIp(tunIp)
	if err != nil {
		logger.Error("Set ip { %s } to tun interface { %s } failed: %v",tunIp,tun.Name(),err)
		exit(logger,1)
	}
	logger.Infof("Created tun interface { %s } and set ip { %s } to it",tun.Name(),tunIp)
	// > }

 ctx,cancel := context.WithCancel(context.Background())
 defer cancel()
 wg := sync.WaitGroup{}

	// < {
	async.AsyncDo(
		ctx,
		&wg,
		func (ctx context.Context){
			for{
				select{
					ctx.Done()
				}
			serveConnectionRW()
			}
		}
	)
	serveConnectionRW()
	// > }
	// < { handle data from tun
	async.AsyncDo(ctx,&wg,func(ctx context.Context){
			utils.HandleDataFromTun(ctx,tun,conn)
	})
	// > }

	// < ready msg chan 

	// >

	// < { handle data from connection
	async.AsyncDo(ctx,&wg,func(ctx context.Context){
		utils.HandleDataFromConn(ctx,tun,conn)
	})
	// > }
	notifyTaskOver := make(chan struct{},1)
	async.AsyncNotifyDone(&wg,notifyTaskOver)


	elegantWaitAndQuit(ctx,logger,cancel,&wg,notifyTaskOver)
}

func elegantWaitAndQuit(ctx context.Context,logger log.LoggerLite,cancel context.CancelFunc,wg *sync.WaitGroup,taskOver chan struct{}){
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

func serveConnectionRW(ctx context.Context,serverAddr string)error{
	var conn *ws.Conn
	var err error
	err = ctx_lib.SelectDoUtilSuc(ctx,time.Second*3,
	func(ctx context.Context)error{
		var err error
		conn,err = readyOnceConn(serverAddr)
		if err != nil{
			return err
		}
		return nil
	},
	)

	if err != nil {
		return err
	}
	defer conn.Close()

	var msgReadFromConn = make(chan interface{},100)
	var msgWillWriteToConn  = make(chan interface{},100)

	// < { read write handle
	wg := sync.WaitGroup{}
	async.AsyncDo(
		ctx,
		&wg,
		func(ctx context.Context){
			loopWriteMsgToConn(ctx,conn,msgWillWriteToConn)
		},
	)
	async.AsyncDo(
		ctx,
		&wg,
		func(ctx context.Context){
			loopReadMsgFromConn(conn,msgReadFromConn)
		},
	)
	wg.Wait()
	// > }
	return nil
}


func loopWriteMsgToConn(ctx context.Context ,conn *ws.Conn,msgWillWriteToConn chan interface{})error{
	defer conn.Close()
	for {
		select{
			case <-ctx.Done():
				return ctx.Err()
			case <-msgWillWriteToConn:
				err := conn.WriteMessage()
				if err != nil{
					return err
				}
		}
	}
}

func loopReadMsgFromConn(conn *ws.Conn,msgReadFromConn chan interface{})(error){
	defer conn.Close()
	for {
		msg,err := conn.ReadMessage()
		if err != nil{
			return err
		}
		// < handle msg
		msgReadFromConn<-msg
		// >
	}
}


func readyOnceConn(serverAddr string)(*ws.Conn,error){
	panic("TODO")
}
