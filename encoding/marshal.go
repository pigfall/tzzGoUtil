package encoding

import (
	"io/ioutil"
)

func UnMarshalByFile(filepath string, v interface{}, unmarshal func([]byte, interface{}) error) error {
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}
	return unmarshal(bytes, v)

}
