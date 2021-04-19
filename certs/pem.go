package certs

import (
    "io/ioutil"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/Peanuttown/tzzGoUtil/fs"
	"os"
	"path"
)

/*
PemX509Save: 以pem 的形式 保存 x509 证书
PemX509FromFile:  从 pem 格式的文件中加载 x509 证书
PemSaveRSAPrivateKey: 以 pem 的格式保存 rsa 私钥
*/

const (
	PEM_BLOCK_TYPE_RSA_PRIVATE_KEY = "RSA PRIVATE KEY"
	PEM_BLOCK_TYPE_RSA_PUBLIC_KEY  = "PUBLIC KEY"
	PEM_BLOCK_TYPE_CRT             = "CERTIFICATE"
)

func PemX509Save(filepath string, crtRaw []byte) error {
	return fs.CreateThen(
		filepath,
		func(file *os.File) error {
			return pem.Encode(
				file,
				&pem.Block{
					Type:  "CERTIFICATE",
					Bytes: crtRaw,
				},
			)
		},
	)
}

func PemX509FromFile(filepath string) (*x509.Certificate, error) {
	var crt *x509.Certificate
	err := fs.ReadAllThen(
		filepath,
		func(c []byte) error {
			block, _ := pem.Decode(c)
			var err error
			crt, err = x509.ParseCertificate(block.Bytes)
			return err
		},
	)
	return crt, err
}

func PemSaveRSAPrivateKey(savePath string, pk *rsa.PrivateKey) error {
	os.MkdirAll(path.Dir(savePath), os.ModePerm)
	file, err := os.Create(savePath)
	if err != nil {
		return err
	}
	defer file.Close()
	return pem.Encode(
		file,
		&pem.Block{
			Type:  PEM_BLOCK_TYPE_RSA_PRIVATE_KEY,
			Bytes: x509.MarshalPKCS1PrivateKey(pk),
		},
	)
}

func PemLoadRSAPrivateKey(filepath string)(privKey *rsa.PrivateKey,err error){
    var fileContent []byte
    fileContent,err = ioutil.ReadFile(filepath)
    if err != nil{
        return
    }
    pemBlock,_ := pem.Decode(fileContent)
    pemPrivBytes :=pemBlock.Bytes
    privKey,err = x509.ParsePKCS1PrivateKey(pemPrivBytes)
    return
}
