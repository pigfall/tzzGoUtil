package reflect


import(

	"fmt"
	"reflect"
)


func IndirectUntilNonPtr(rvInput reflect.Value)(reflect.Value){
	var rv = rvInput
	for{
		fmt.Println(rv.Kind())
		if rv.Kind() == reflect.Ptr{
			rv = rv.Elem()
			continue
		}
		return rv
	}
}
