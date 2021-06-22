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
