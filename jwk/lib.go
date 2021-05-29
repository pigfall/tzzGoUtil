package jwk

import (
	"crypto/rsa"
	"encoding/json"
	"github.com/lestrrat-go/jwx/jwk"
)

func ConvertPublicKeyToJWKStr(pubKey *rsa.PublicKey) (string, error) {
	jwkKey, err := jwk.New(pubKey)
	if err != nil {
		return "", err
	}
	bytes, err := json.Marshal(jwkKey)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func BuildSymmetricJWK(symmetricKey []byte) (string, error) {
	jwkKey, err := jwk.New(symmetricKey)
	if err != nil {
		return "", err
	}
	bytes, err := json.Marshal(jwkKey)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
