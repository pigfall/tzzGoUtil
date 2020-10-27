package certs

import (
	crand "crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"math/rand"
	"time"
)

/*
NewX509CaCrtTpl: 生成 ca 证书模板
NewX509CertTpl :x509 证书模板
SignSelf: 自签名证书
SignCrt: 签名证书
*/

// 生成 ca 证书模板, 只要区别就是:
//c.BasicConstraintsValid= true
//c.IsCA= true
func NewX509CaCrtTpl(
	subject pkix.Name,
	notBefore time.Time,
	validDuration time.Duration,
	option func(c *x509.Certificate),
) *x509.Certificate {
	return NewX509CertTpl(
		subject,
		notBefore,
		validDuration,
		func(c *x509.Certificate) {
			c.BasicConstraintsValid = true
			c.IsCA = true
			if option != nil {
				option(c)
			}
		},
	)
}

// x509 证书模板
func NewX509CertTpl(
	subject pkix.Name,
	notBefore time.Time,
	validDuration time.Duration,
	option func(c *x509.Certificate),
) *x509.Certificate {
	cert := &x509.Certificate{
		Subject:      subject,
		NotBefore:    notBefore,
		NotAfter:     notBefore.Add(validDuration),
		SerialNumber: big.NewInt(rand.Int63()),
	}
	if option != nil {
		option(cert)
	}
	return cert
}

func SignSelf(crtTpl *x509.Certificate, privKey *rsa.PrivateKey) (*x509.Certificate, error) {
	return SignCrt(crtTpl, crtTpl, &privKey.PublicKey, privKey)
}

func SignCrt(parentCrt *x509.Certificate, crtTpl *x509.Certificate, hisPublicKey *rsa.PublicKey, parentPrivateKey *rsa.PrivateKey) (signedCrt *x509.Certificate, err error) {
	crtRaw, err := x509.CreateCertificate(
		crand.Reader,
		crtTpl,
		parentCrt,
		hisPublicKey,
		parentPrivateKey,
	)
	if err != nil {
		return nil, err
	}
	signedCrt, err = x509.ParseCertificate(crtRaw)
	return signedCrt, err
}
