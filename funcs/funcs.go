package funcs


type  Funcs struct{
	funcs []func()
}

func NewFuncs(funcs ...func())*Funcs{
	funcsWrapper := &Funcs{
		funcs :make([]func(),0),
	}
	funcsWrapper.AddFunc(funcs...)

	return funcsWrapper
}

func (this *Funcs) AddFunc(funcs ...func()){
	this.funcs= append(this.funcs,funcs...)
}

func (this *Funcs) Call(){
	for _,f := range this.funcs{
		f()
	}
}
