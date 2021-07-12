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
