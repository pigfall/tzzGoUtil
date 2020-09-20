package yaml

import (
	libYaml "gopkg.in/yaml.v2"
)


func Marshal(v interface{})([]byte,error){
	return libYaml.Marshal(v)
}

func UnMarshal(data []byte,v interface{})(error){
	return libYaml.Unmarshal(data,v)
}
