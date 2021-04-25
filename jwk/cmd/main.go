package main

import(
    "fmt"
    "github.com/Peanuttown/tzzGoUtil/certs"
    "github.com/Peanuttown/tzzGoUtil/jwk"

    "os"
)


func main() {
    pemFilePath := os.Args[1]
    // targetFilePath := os.Args[2]

    pubKey,err := certs.PemLoadRSAPublicKey(pemFilePath)
    if err != nil{
        panic(err)
    }
    jwkStr,err := jwk.ConvertPublicKeyToJWKStr(pubKey)
    if err != nil{
        panic(err)
    }
    fmt.Println(jwkStr)
}
