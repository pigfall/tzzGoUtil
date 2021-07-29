package encoding

import (
	"io/ioutil"
	"os"
)

func UnMarshalByFile(filepath string, v interface{}, unmarshal func([]byte, interface{}) error) error {
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}
	return unmarshal(bytes, v)

}

func MarshalToFile(filepath string, v interface{},marshalFunc func(v interface{})([]byte,error))error{
	bytes,err := marshalFunc(v)
	if err != nil{
		return err
	}
	return ioutil.WriteFile(filepath,bytes,os.ModePerm)
}
