package jwk

import(
    "encoding/json"
    "github.com/lestrrat-go/jwx/jwk"
    "crypto/rsa"
)

func ConvertPublicKeyToJWKStr(pubKey *rsa.PublicKey)(string,error){
    jwkKey,err := jwk.New(pubKey)
    if err != nil{
        return "",err
    }
    bytes,err := json.Marshal(jwkKey)
    if err != nil{
        return "",err
    }
    return string(bytes),nil
}
