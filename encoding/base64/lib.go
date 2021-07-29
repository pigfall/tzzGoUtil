package base64

import(
		"encoding/base64"
)

func Decode(toDecode string) (string,error){
	decodedBytes,err := base64.StdEncoding.DecodeString(toDecode)
	if err != nil{
		return "", err
	}
	return string(decodedBytes), nil
}
