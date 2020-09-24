package certs

import(
    "os"
    "crypto/x509"
    "encoding/pem"
    "github.com/Peanuttown/tzzGoUtil/fs"
)


func PemX509Save(filepath string,crtRaw []byte)(error){
    return fs.CreateThen(
        filepath,
        func(file *os.File)(error){
            return pem.Encode(
                file,
                &pem.Block{
                    Type:"CERTIFICATE",
                    Bytes:crtRaw,
                },
            )
        },
    )
}

func PemX509FromFile(filepath string)(*x509.Certificate,error){
    var crt *x509.Certificate
    err := fs.ReadAllThen(
        filepath,
        func(c []byte)(error){
            block,_:=pem.Decode(c)
            var err error
            crt,err =x509.ParseCertificate(block.Bytes)
            return err
        },
    )
    return crt,err
}
