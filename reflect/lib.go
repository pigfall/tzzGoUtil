package reflect

import(
    "fmt"
    "reflect"
)

type RtRvI interface{
    Kind() reflect.Kind
}

/*
func SetValue(entity interface{},field string,value interface{})(err error)
func IsNumber(rt reflect.Type)bool
func IsUnsignedNumber(rt reflect.Type)bool
func IndirectType(rt reflect.Type)(reflect.Type)
func IndirectValue(rv reflect.Value)(reflect.Value)
func GetStructName(v interface{}) string 
*/



func SetValue(entity interface{},field string,value interface{})(err error){
    rt := reflect.TypeOf(entity)
    if rt.Kind() != reflect.Ptr{
        return fmt.Errorf("entity must be pointer")
    }
    rt = IndirectType(rt)
    _,exist:= rt.FieldByName(field)
    if !exist{
        return fmt.Errorf("not find field :%s",field)
    }
    rv :=(reflect.ValueOf(entity))
    rvField:= rv.Elem().FieldByName(field)
    if !rvField.CanSet(){
        return fmt.Errorf("field %s can not set",field)
    }
    defer func(){
        r := recover()
        if r != nil{
            err = fmt.Errorf("%v",r)
        }
    }()
    rvField.Set(reflect.ValueOf(value))
    //kind:=rtField.Type.Kind()
    //if kind == reflect.Bool{
    //    rvField.SetBool(value.(bool))
    //}else if IsNumber(rtField.Type){
    //    rvField.SetInt()

    //}
    return
}

func IsNumber(rt reflect.Type)bool{
    return rt.Kind() ==reflect.Int8 || 
        rt.Kind()==reflect.Int16 ||
        rt.Kind()==reflect.Int32 ||
        rt.Kind()==reflect.Int64 ||
        rt.Kind()==reflect.Int 
}

func IsFloatNumber(rt reflect.Type)bool{
	return rt.Kind() == reflect.Float32 || rt.Kind() == reflect.Float64
}

func IsUnsignedNumber(rt reflect.Type)bool{
	return rt.Kind()==reflect.Uint ||
	rt.Kind()==reflect.Uint8||
	rt.Kind()==reflect.Uint16||
	rt.Kind()==reflect.Uint32||
	rt.Kind()==reflect.Uint64

}


func IndirectType(rt reflect.Type)(reflect.Type){
    if rt.Kind() == reflect.Ptr{
        return rt.Elem()
    }
    return rt
}

func IndirectValue(rv reflect.Value)(reflect.Value){
    if rv.Kind() == reflect.Ptr{
        return  rv.Elem()
    }
    return rv
}

func GetStructName(v interface{}) string {
    rt := reflect.TypeOf(v)
    if rt.Kind() == reflect.Ptr{
        rt = rt.Elem()
    }
    return rt.Name()
}

func IndirectUntilNonPtr(rv reflect.Value)reflect.Value{
	var ret = rv
	for{
		if ret.Kind() == reflect.Ptr{
			ret = ret.Elem()
			continue
		}
		return ret
	}
}
