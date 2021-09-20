package main

import (
	"fmt"
	"github.com/pigfall/tzzGoUtil/certs"
	"github.com/pigfall/tzzGoUtil/jwk"

	"os"
)

func main() {
	jwkKeyTpe := os.Args[1]
	switch jwkKeyTpe {
	case "rsaPubKey":
		pemFilePath := os.Args[2]
		// targetFilePath := os.Args[2]

		pubKey, err := certs.PemLoadRSAPublicKey(pemFilePath)
		if err != nil {
			panic(err)
		}
		jwkStr, err := jwk.ConvertPublicKeyToJWKStr(pubKey)
		if err != nil {
			panic(err)
		}
		fmt.Println(jwkStr)
	case "symmetric":
		jwkStr, err := jwk.BuildSymmetricJWK([]byte(os.Args[2]))
		if err != nil {
			panic(err)
		}
		fmt.Println(jwkStr)

	}
	return
}
