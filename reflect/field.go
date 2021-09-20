package reflect

import(
	"strconv"
	"strings"
		"reflect"
		"fmt"
		"github.com/pigfall/tzzGoUtil/ascii"
)

func getStructEmptyField(obj interface{},isEmptyField func (rv reflect.Value,strcutField reflect.StructField)(isEmpty bool,notRecursive bool))([]string){
	fmt.Println("obj is ",reflect.TypeOf(obj).Name())
	if obj == nil{
		return []string{fmt.Sprintf("%s is nil",reflect.TypeOf(obj).Name())}
	}
	rawRv := reflect.ValueOf(obj)
	nonPtrRv := IndirectUntilNonPtr(rawRv)
	if nonPtrRv.Kind() != reflect.Struct{
		fmt.Println("not struct ",nonPtrRv.Type().Name())
		return nil
	}
	return	getStructEmptyFieldHelper("",nonPtrRv,isEmptyField)
}

func getStructEmptyFieldHelper(parentFieldName string,rv reflect.Value,isEmptyField func(rv reflect.Value,fieldStruct reflect.StructField)(isEmpty bool,notRecursive bool))([]string){
	fmt.Println(parentFieldName,"  ",rv.Type().Name())
	rets := make([]string,0)
	for i:=0;i<rv.NumField();i++{
		fieldRv := rv.Field(i)
		fieldStruct := rv.Type().Field(i)
		emptyReport := getFieldEmpty(parentFieldName,fieldRv,fieldStruct,isEmptyField)
		if len(emptyReport) > 0 {
			rets =append(rets,emptyReport...)
		}
	}
	return rets
}

func getFieldEmpty(parentFieldName string,rvField reflect.Value,fieldStruct reflect.StructField,isEmptyField func(reflect.Value,reflect.StructField)(isEmpty bool,notRecursive bool))[]string{
	isEmpty,notRecursive :=  isEmptyField(rvField,fieldStruct)
	if isEmpty{
		return []string{fmt.Sprintf("%s.%s",parentFieldName,fieldStruct.Name)}
	}
	if notRecursive {
		return nil
	}
	rvNonPtr := IndirectUntilNonPtr(rvField)
	if  rvNonPtr.Kind()== reflect.Struct{
		return getStructEmptyFieldHelper(fmt.Sprintf("%s.%s",parentFieldName,fieldStruct.Name),rvNonPtr,isEmptyField)
	}
	return nil
}

func fieldEmptyIfIsDefaultValue(rv reflect.Value,structField reflect.StructField)(isEmpty bool,notRecursive bool){
	return rv.IsZero(),false
}
func publicFieldEmptyIfIsDefaultValue(rv reflect.Value,structField reflect.StructField)(isEmpty bool,notRecursive bool){
	fmt.Println(structField.Name)
	if ascii.IsLowerAlpha(structField.Name[0]){
		return false,true
	}
	return fieldEmptyIfIsDefaultValue(rv,structField)
}

func GetStructEmptyField(obj interface{})([]string){
	return getStructEmptyField(obj,fieldEmptyIfIsDefaultValue)
}

func GetStructPublicEmptyField(obj interface{})([]string){
	fields := getStructEmptyField(obj,publicFieldEmptyIfIsDefaultValue)
	rets := make([]string,0,len(fields))
	for _,f := range fields{
		fTrimDot := strings.TrimPrefix(f,".")
		if ascii.IsLowerAlpha(fTrimDot[0]){
			continue
		}
		rets = append(rets,fTrimDot)
	}
	return rets
}

func ForEachField(obj interface{},do func(rvField reflect.Value,structFields reflect.StructField)error)error{
	rawRv := reflect.ValueOf(obj)
	rvIndirect := IndirectValue(rawRv)
	rtIndirect := rvIndirect.Type()
	for i:=0;i < rtIndirect.NumField();i++ {
		err := do(rvIndirect.Field(i),rtIndirect.Field(i))
		if err != nil{
			return err
		}
	}
	return nil
}

func FieldNum(structObj interface{})int{
	rv := IndirectStructObj(structObj)
	return rv.NumField()
}

func IndirectStructObj(structObj interface{})reflect.Value{
	rawRv := reflect.ValueOf(structObj)
	return IndirectValue(rawRv)
}

func ToString(rv reflect.Value)(string){
	if IsNumber(rv.Type()){
		return strconv.FormatInt(rv.Int(),10)
	}
	if IsUnsignedNumber(rv.Type()){
		return strconv.FormatUint(rv.Uint(),10)
	}

	if IsFloatNumber(rv.Type()){
		return fmt.Sprintf("%.2f",rv.Float())
	}
	if rv.Kind() == reflect.String{
		return rv.String()
	}
	return rv.Type().String()
}
