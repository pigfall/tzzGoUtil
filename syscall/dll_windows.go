package syscall

import(
	"syscall"
)


type DLL struct{
	*syscall.DLL
}

type Procdure struct{
	*syscall.Proc
}



func LoadDLL(dllPath string)(*DLL,error){
	dll,err := syscall.LoadDLL(dllPath)
	if err != nil {
		return nil,err
	}
	return &DLL{
		DLL:dll,
	},nil
}

func (this *DLL) FindProcure(produreName string)(*Procdure,error){
	proc,err := this.FindProc(produreName)
	if err !=nil{
		return nil,err
	}
	return &Procdure{
		Proc:proc,
	},nil
}

func (this *Procdure) Call(a ...uintptr)(r1,r2 uintptr,err error){
	r1,r2,err = this.Proc.Call(a...)
	if r1 == 0{
		return  r1,r2,err
	}
	return r1,r2,nil
}
