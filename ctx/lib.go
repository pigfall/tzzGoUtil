package ctx 

import(
	"time"
	"context"
)



func SelectDoUtilSuc(ctx context.Context,delay time.Duration,do func(ctx context.Context)error)error{
	err := do(ctx)
	if err  == nil{
		return nil
	}
	ticker := time.NewTicker(delay)
	for {
		select{
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			err := do(ctx)
			if err != nil{
				ticker.Reset(delay)
				continue
			}
			return nil
		}
	}
}
