package http


type Options struct{
	successStatusCode []int
}

type OptionFiller func(op *Options)

func OptionSuccessStatusCode(successStatusCode []int)OptionFiller{
	return func(op *Options){
		op.successStatusCode = successStatusCode
	}
}


func newOptions(ops ...OptionFiller)*Options{
	op := &Options{}
	for _,filler := range ops{
		filler(op)
	}
	return op
}

func (this *Options) StatusCodeOk(statusCode int)(bool){
	for _,v := range this.successStatusCode{
		if v == statusCode {
			return true
		}
	}
	return false
}
