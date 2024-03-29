package main

import(
	"context"
	"sync"
		"net/http"
	"github.com/pigfall/tzzGoUtil/net/water_tun_wrap/cmd_test/utils"	

	"github.com/pigfall/tzzGoUtil/log"
	"github.com/pigfall/tzzGoUtil/async"
	ws "github.com/gorilla/websocket"
)



func NewConnHandler(logger log.LoggerLite,ctx context.Context,msgReadFromTun chan utils.MsgReadFromTun,msgWillWriteToTun chan utils.MsgWillWriteToTun)http.Handler{
	return &ConnHandler{
		logger:logger,
		ctx:ctx,
		msgReadFromTun:msgReadFromTun,
		msgWillWriteToTun:msgWillWriteToTun,
	}

}


type ConnHandler struct{
	logger log.LoggerLite
	ctx context.Context
	msgReadFromTun chan utils.MsgReadFromTun
	msgWillWriteToTun chan utils.MsgWillWriteToTun
}


func (this *ConnHandler)  ServeHTTP(w http.ResponseWriter, req *http.Request){
	logger := this.logger
	logger.Info("rcv one http request")
	upgrader := ws.Upgrader{}
	conn,err := upgrader.Upgrade(w,req,nil)
	if err != nil{
		logger.Error("upgrad to web socket failed:",err)
		return
	}
	defer conn.Close()
	wg := sync.WaitGroup{}
	msgWillWriteToConn :=make(chan []byte,100)
	msgReadFromConn := make(chan []byte,100)
	async.AsyncDo(
		this.ctx,
		&wg,
		func(ctx context.Context){
			utils.LoopReadMsgFromConn(conn,logger,msgReadFromConn)
		},
	)
	async.AsyncDo(
		this.ctx,
		&wg,
		func(ctx context.Context){
			utils.LoopWriteMsgToConn(ctx,logger,conn,msgWillWriteToConn)
		},
			
	)
	async.AsyncDo(
		this.ctx,
		&wg,
		func(ctx context.Context){
			<-ctx.Done()
			conn.Close()
			close(msgWillWriteToConn)
		},
	)
	async.AsyncDo(
		this.ctx,
		&wg,
		func(ctx context.Context){
			utils.HandleDataFromConn(
				ctx,
				logger,
				msgReadFromConn,
				this.msgWillWriteToTun,

			)
		},
	)
	async.AsyncDo(
		this.ctx,
		&wg,
		func(ctx context.Context){
			utils.HandleDataFromTun(
				ctx,
				logger,
				this.msgReadFromTun,
				msgWillWriteToConn,
			)
		},
	)
	wg.Wait()
}



 
