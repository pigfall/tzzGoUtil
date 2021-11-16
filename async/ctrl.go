package async

import(
		"sync"
		"context"
)

type Ctrl struct{
	wg sync.WaitGroup
	cancels []func()
	onRoutineQuit func(jobName string)
}

func NewCtrl()*Ctrl{
	return &Ctrl{
		cancels:make([]func(),0),
	}
}




func (this *Ctrl) AsyncDo(
	ctx context.Context,
	jobName string,
	do func(ctx context.Context),
){
	this.wg.Add(1)
	go func(ctx context.Context){
		defer func(){
			if this.onRoutineQuit != nil{
				this.onRoutineQuit(jobName)
			}
			this.wg.Done()
		}()
		do(ctx)
	}(ctx)

}

func (this *Ctrl) AppendCancelFuncs(cancels ...func()){
	if this.cancels == nil{
		this.cancels = cancels
	}else{
		this.cancels = append(this.cancels,cancels...)
	}

}

func (this *Ctrl) Cancel(){
	for _,cancel := range this.cancels{
		cancel()
	}

}

func (this *Ctrl) Wait(){
	this.wg.Wait()
}

func (this *Ctrl) OnRoutineQuit(f func(jobName string)){
	this.onRoutineQuit = f
}

func (this *Ctrl) WaitCtx(ctx context.Context,doWhenCtxDone func()){
	<-ctx.Done()
	doWhenCtxDone()
}
