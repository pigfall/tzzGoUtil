package main

import(
    "encoding/json"
    "fmt"
    "github.com/Peanuttown/tzzGoUtil/certs"

)

const SUB_CMD_NAME_RSA="rsa"

type subCmdRsaConf struct{
    OutPrivateKeyPath string `json:"out_private_key_path"`
    OutPublickeyPath string `json:"out_publick_key_path"`
    KeyBitSize int `json:"key_bit_size"`
}

const demo_rsa_conf =`
{
    "out_private_key_path":"/tmp/town_private.key",
    "out_publick_key_path": "/tmp/town_public.pub",
    "key_bit_size": 256
}
`

func subCmdRsa(cfgBytes []byte) error {
    conf := &subCmdRsaConf{}
    err := json.Unmarshal(cfgBytes,conf)
    if err != nil{
        err = fmt.Errorf("< UnMarshal config `%s` failed >: %w",string(cfgBytes),err)
        return err
    }

    var outPrivateKeyPath  = "town_rsa.key"
    var outPublicKeyPath = "town_rsa.pub"
    // < check config content
    if len(conf.OutPrivateKeyPath) > 0{
        outPrivateKeyPath = conf.OutPrivateKeyPath
    }
    if len(conf.OutPublickeyPath) > 0 {
        outPublicKeyPath = conf.OutPublickeyPath
    }
    if conf.KeyBitSize == 0{
        return fmt.Errorf("key bit can not be 0")
    }

    // >

    privKey,err := certs.RSAGenPrivateKey(certs.PrivateKeyBitSize(conf.KeyBitSize))
    if err != nil{
        return err
    }

    // < save to file as pem format
    err = certs.PemSaveRSAPrivateKey(outPrivateKeyPath,privKey)
    if err != nil{
        return fmt.Errorf("< Save privateKey to %s failed >: %w",outPrivateKeyPath,err)
    }
    fmt.Printf("< Success save privateKey to %s >\n",outPrivateKeyPath)

    err = certs.PemSaveRSAPublicKey(outPublicKeyPath,&privKey.PublicKey)
    if err != nil{
        return fmt.Errorf("< Save publick key to %s failed >: %w",outPublicKeyPath,err)
    }
    fmt.Printf("< Success save publicKey to %s >\n",outPublicKeyPath)
    // >

    return nil
}
