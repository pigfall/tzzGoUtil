package async

import(
	"context"
	"sync"
)


func AsyncDo(ctx context.Context,wg *sync.WaitGroup,do func(ctx context.Context)){
	wg.Add(1)
	go func(){
		defer wg.Done()
		do(ctx)
	}()
}

func  AsyncDoWithCancel(ctx context.Context,wg *sync.WaitGroup,do func(ctx context.Context),cancel context.CancelFunc){
	wg.Add(1)
	go func(ctx context.Context){
		defer func(){
			cancel()
			wg.Done()
		}()
		do(ctx)
	}(ctx)
}


func AsyncNotifyDone(wg *sync.WaitGroup, ch chan struct{}){
	go func(){
		wg.Wait()
		close(ch)
	}()
}

