package utils

import(
	"context"
"github.com/Peanuttown/tzzGoUtil/log"		
	ws "github.com/gorilla/websocket"
)

func LoopReadMsgFromConn(conn *ws.Conn,logger log.LoggerLite,msgReadFromConn chan []byte)(error){
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

func LoopWriteMsgToConn(ctx context.Context,logger log.LoggerLite ,conn *ws.Conn,msgWillWriteToConn chan []byte)error{
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
