package hmac

import(
	stdhmac  "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
)


func EnCryptAndEncodeToHex(secret []byte,contentToEncrypt []byte)(string,error){
	h :=stdhmac.New(sha256.New,secret)
	_,err := h.Write(contentToEncrypt)
	if err != nil{
		return "",err
	}
	return hex.EncodeToString(h.Sum(nil)),nil
}

func EnCrypt(secret []byte,contentToEncrypt []byte)([]byte,error){
	h :=stdhmac.New(sha256.New,secret)
	_,err := h.Write(contentToEncrypt)
	if err != nil{
		return nil,err
	}
	return h.Sum(nil),nil
}
