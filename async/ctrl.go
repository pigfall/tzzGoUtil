package async

import(
		"sync"
		"context"
)

type Ctrl struct{
	wg sync.WaitGroup
	cancels []func()
}




func (this *Ctrl) AsyncDo(
	ctx context.Context,
	do func(ctx context.Context),
){
	this.wg.Add(1)
	go func(ctx context.Context){
		defer this.wg.Done()
		do(ctx)
	}(ctx)

}

func (this *Ctrl) AppendCancelFunc(cancel func()){
	if this.cancels == nil{
		this.cancels = []func(){cancel}
	}else{
		this.cancels = append(this.cancels,cancel)
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
