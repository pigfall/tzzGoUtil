package utils



import(
	"context"
	"github.com/Peanuttown/tzzGoUtil/log"
)



func HandleDataFromConn(ctx context.Context,logger log.LoggerLite,msgReadFromConn chan []byte,msgWillWriteToTun  chan MsgWillWriteToTun){
	for{
		select{
		case <-ctx.Done():
			logger.Info("HandleDataFromConn ctx done")
		case msgReadFromConnOne :=<-msgReadFromConn:
			logger.Info("Read one msg from Conn Channel ")
			msgWillWriteToTun<-(msgReadFromConnOne)
		}
	}
}
