package ctx

import (

	"time"
	"context"
		
)


func TickerDo(ctx context.Context,tickerDuration time.Duration,do func()error)(error){
	err :=do()
	if err != nil{
		return err
	}
	ticker := time.NewTicker(tickerDuration)
	defer ticker.Stop()
	for{
		select{
		case  <-ticker.C:
			err := do()
			if err != nil{
				return err
			}
		case <-ctx.Done():
			return nil
		}
	}
}
