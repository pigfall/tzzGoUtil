package atomic


import(

	"sync/atomic"
)

type Int64 struct{
	value int64
}

func (this* Int64) Add(dt int64)(now int64){
	return	atomic.AddInt64(&this.value,dt)
}
