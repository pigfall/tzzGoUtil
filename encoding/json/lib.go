package json

import(
		"github.com/Peanuttown/tzzGoUtil/encoding"
		std_json"encoding/json"
)

func JsonUnmarshalFromFile(filepath string,entity interface{})error{
	return encoding.UnMarshalByFile(filepath,entity,std_json.Unmarshal)
}
