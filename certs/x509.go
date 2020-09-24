package certs


import(
    "crypto/x509"
    "crypto/x509/pkix"
    "time"
    "math/big"
)

func NewX509Cert(
    subject pkix.Name,
    notBefore time.Time,
    validDuration time.Duration,
    option func (c *x509.Certificate),
)*x509.Certificate{
    cert := &x509.Certificate{
        Subject:subject,
        NotBefore:notBefore,
        NotAfter:notBefore.Add(validDuration),
        SerialNumber:big.NewInt(1),
    }
    if option != nil{
        option(cert)
    }
    return cert
}


