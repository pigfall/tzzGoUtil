package certs

import(
    "time"
    "testing"
    "crypto/x509/pkix"
    "crypto/x509"
    "crypto/rand"
    "crypto/rsa"
)

func TestCaCrt(t *testing.T){
    caKey,err := rsa.GenerateKey(rand.Reader,2048)
    if err != nil{
        t.Fatal(err)
    }
    var now =time.Now()
    var hour = time.Hour
    caTpl := NewX509CaCrtTpl(
        pkix.Name{CommonName:"testCa"},
        now,
        hour,
        nil,
    )
    caRaw,err := x509.CreateCertificate(
        rand.Reader,
        caTpl,
        caTpl,
        &caKey.PublicKey,
        caKey,
    )
    if err != nil{
        t.Fatal(err)
    }
    caCrt,err :=x509.ParseCertificate(caRaw)
    if err != nil{
        t.Fatal(err)
    }

    serverCrtTpl := NewX509CertTpl(
        pkix.Name{
            CommonName:"tzzServer",
        },
        now,
        hour,
        nil,
    )
    serverPrivKey,err:=rsa.GenerateKey(rand.Reader,2048)
    if err != nil{
        t.Fatal(err)
    }
    serverCrt,err:=SignCrt(caCrt,serverCrtTpl,&serverPrivKey.PublicKey,caKey)
    if err != nil{
        t.Fatal(err)
    }

    crtPool :=x509.NewCertPool()
    crtPool.AddCert(caCrt)
    _,err = serverCrt.Verify(x509.VerifyOptions{Roots:crtPool})
    if err != nil{
        t.Fatal(err)
    }
    t.Log("证书验证通过")
}
