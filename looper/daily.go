package looper

import(
		"time"
		"context"
)


// hour [0,23]
func LoopDailyDo(ctx context.Context,hour int,do func(ctx context.Context)error)error{
	if hour < 0 || hour > 23{
		panic("input hour must in [0,23]")
	}
	durationTick := time.Hour/2
	ticker := time.NewTicker(durationTick)
	lastExecDay := -1
	for{
		select{
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:

		}
		now := time.Now()
		if lastExecDay ==now.Day(){
			ticker.Reset(durationTick)
			continue
		}
		nowHour := now.Hour()
		if hour != nowHour{
			ticker.Reset(durationTick)
			continue
		}

		lastExecDay = now.Day()
		err := do(ctx)
		if err != nil{
			return err
		}
	}
}
