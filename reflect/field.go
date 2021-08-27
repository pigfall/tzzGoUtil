package reflect

import(
	"reflect"
)


func FieldStringIsNilOrLenZero(obj interface{})([]string){
	var strLenZeroOrNil= []string{}
	rawRv := reflect.ValueOf(obj)
	rvNonPtr := IndirectUntilNonPtr(rawRv)
	if rvNonPtr.Kind() != reflect.Struct{
		return nil
	}

	rvStruct := rvNonPtr
	rtStruct := rvStruct.Type()
	for fieldIndex:=0;fieldIndex<rtStruct.NumField();fieldIndex++{
		rawStructField := rvStruct.Field(fieldIndex)
		rawRtStructField := rtStruct.Field(fieldIndex)
		if rawRtStructField.Type.Kind() == reflect.String{
			if s := rawStructField.String(); len(s) == 0{
				strLenZeroOrNil = append(strLenZeroOrNil,rtStruct.Field(fieldIndex).Name)
			}
		}else if rawRtStructField.Type.Kind() == reflect.Ptr{
			if rawStructField.IsNil(){
				strLenZeroOrNil = append(strLenZeroOrNil,rtStruct.Field(fieldIndex).Name)
			}
		}
	}

	return strLenZeroOrNil
}

