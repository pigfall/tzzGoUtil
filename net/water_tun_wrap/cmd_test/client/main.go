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
	flag.StringVar(&serverAddr,"svr","","server address")
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
			logger.Infof("Will retry after %d second",retryDelaySec)
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

 msgReadFromConn := make(chan []byte,100)
 msgWillWriteToConn  := make(chan []byte,100)
	msgWillWriteToTun := make(chan utils.MsgWillWriteToTun,100)
	msgReadFromTun := make(chan utils.MsgReadFromTun,100)
	cancelFunc := func(){
		cancel()
		close(msgWillWriteToTun )
		close(msgReadFromTun)
		close(msgReadFromConn)
		close(msgWillWriteToConn)
	}
	// < keep connection and loop for read write{
	async.AsyncDo(
		ctx,
		&wg,
		func (ctx context.Context){
			for{
				select{
				case <-ctx.Done():
					return
				default:
				}
				err := serveConnectionRW(ctx,logger,serverAddr,msgReadFromConn,msgWillWriteToConn)
				if err != nil{
					logger.Error(err)
				}
			}
		},
	)
	// > }


	// < loop for read write tun interface
	async.AsyncDo(
		ctx,
		&wg,
		func (ctx context.Context){
			utils.LoopRWTun(ctx,logger,tun,msgReadFromTun,msgWillWriteToTun)
		},
	)
	// >

	// < { handle data from tun
	async.AsyncDo(ctx,&wg,func(ctx context.Context){
			utils.HandleDataFromTun(ctx,logger,msgReadFromTun,msgWillWriteToConn)
	})
	// > }

	// < { handle data from connection
	async.AsyncDo(ctx,&wg,func(ctx context.Context){
		utils.HandleDataFromConn(ctx,logger,msgReadFromConn,msgWillWriteToTun)
	})
	// > }
	notifyTaskOver := make(chan struct{},1)
	async.AsyncNotifyDone(&wg,notifyTaskOver)


	elegantWaitAndQuit(ctx,logger,cancelFunc,&wg,notifyTaskOver)
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

func serveConnectionRW(ctx context.Context,logger log.LoggerLite,serverAddr string,msgReadFromConn chan []byte,msgWillWriteToConn chan []byte)error{
	var conn *ws.Conn
	var err error
	err = ctx_lib.SelectDoUtilSuc(ctx,time.Second*3,
	func(ctx context.Context)error{
		var err error
		conn,err = readyOnceConn(serverAddr,logger)
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


	// < { read write handle
	wg := sync.WaitGroup{}
	async.AsyncDo(
		ctx,
		&wg,
		func(ctx context.Context){
			loopWriteMsgToConn(ctx,logger,conn,msgWillWriteToConn)
		},
	)
	async.AsyncDo(
		ctx,
		&wg,
		func(ctx context.Context){
			loopReadMsgFromConn(conn,logger,msgReadFromConn)
		},
	)
	wg.Wait()
	// > }
	return nil
}


func loopWriteMsgToConn(ctx context.Context,logger log.LoggerLite ,conn *ws.Conn,msgWillWriteToConn chan []byte)error{
	defer conn.Close()
	for {
		select{
			case <-ctx.Done():
				return ctx.Err()
			case bytes:= <-msgWillWriteToConn:
				logger.Debug("Read on msg from Conn Channel, wil write it to conn")
				err := conn.WriteMessage(ws.BinaryMessage,bytes)
				if err != nil{
					logger.Error("Failed to write bytes to connection ",err)
					return err
				}
				logger.Debug("Ok to write bytes conn")
		}
	}
}

func loopReadMsgFromConn(conn *ws.Conn,logger log.LoggerLite,msgReadFromConn chan []byte)(error){
	defer conn.Close()
	for {
		_,msg,err := conn.ReadMessage()
		if err != nil{
			logger.Error("Failed to read from conn ",err)
			return err
		}
		// < handle msg
		logger.Debug("Will send msg to readFromConnChanel ")
		msgReadFromConn<-msg
		logger.Debug("OK send msg to readFromConnChanel ")
		// >
	}
}


func readyOnceConn(serverAddr string,logger log.LoggerLite)(*ws.Conn,error){
	logger.Info("Connecing to ",serverAddr)
	conn,_,err := ws.DefaultDialer.Dial(serverAddr,nil)
	if err !=nil{
		logger.Error(fmt.Errorf("Failed to connect to %s, %w",serverAddr,err))
		return nil,err
	}
	logger.Info("Ok to connect to server ", serverAddr)

	return conn,nil
}
