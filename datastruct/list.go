package datastruct

import(
		"container/list"
)

type List struct{
	*list.List
}


func NewList()*List{
	return &List{
		List:list.New(),
	}
}

func (this *List) ForEach(do func (e *list.Element)error)error{
	 e:= this.List.Front()
	 for (e!=nil){
		 err := do(e)
		 if err != nil{
			 return err
		 }
		 e = e.Next()
	 }

	 return nil
}
