package reflect

import(
		"reflect"
		"fmt"
)

func getStructEmptyField(obj interface{},isEmptyField func (rv reflect.Value,strcutField reflect.StructField)(bool))([]string){
	rawRv := reflect.ValueOf(obj)
	nonPtrRv := IndirectUntilNonPtr(rawRv)
	if nonPtrRv.Kind() != reflect.Struct{
		return nil
	}
	return	getStructEmptyFieldHelper("",nonPtrRv,isEmptyField)
}

func getStructEmptyFieldHelper(parentFieldName string,rv reflect.Value,isEmptyField func(rv reflect.Value,fieldStruct reflect.StructField)bool)([]string){
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

func getFieldEmpty(parentFieldName string,rvField reflect.Value,fieldStruct reflect.StructField,isEmptyField func(reflect.Value,reflect.StructField)bool)[]string{
	if isEmptyField(rvField,fieldStruct){
		return []string{fmt.Sprintf("%s.%s",parentFieldName,fieldStruct.Name)}
	}
	if rvField.Kind() == reflect.Struct{
		return getStructEmptyFieldHelper(fmt.Sprintf("%s.%s",parentFieldName,fieldStruct.Name),rvField,isEmptyField)
	}
	return nil
}

func fieldEmptyIfIsDefaultValue(rv reflect.Value,structField reflect.StructField)(bool){
	return rv.IsZero()
}

func GetStructEmptyField(obj interface{})([]string){
	return getStructEmptyField(obj,fieldEmptyIfIsDefaultValue)
}
