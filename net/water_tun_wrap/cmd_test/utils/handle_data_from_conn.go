package utils



import(
	ws "github.com/gorilla/websocket"
	"github.com/Peanuttown/tzzGoUtil/async"
	"github.com/Peanuttown/tzzGoUtil/net"	
	"context"
	"github.com/Peanuttown/tzzGoUtil/log"
)



func HandleDataFromConn(ctx context.Context,logger log.LoggerLite,tunMsgChan chan <-interface{},connMsgChan <-chan interface{}){
	wg := sync.WaitGroup{}
	async.AsyncDo(
		ctx,&wg,
		func(ctx context.Context){
			for{
					conn.ReadMessage()
			}
		},
	)
	for{
		select{
		case <-ctx.Done():
			logger.Info("Ctx done")
		case 
		}
	}
}
