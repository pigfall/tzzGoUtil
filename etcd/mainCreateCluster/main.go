//----------
// AUTHOR: tzz
// DESC : 快速创建 etcd 集群
//----------
package main

import (
	"crypto/x509/pkix"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/Peanuttown/tzzGoUtil/certs"
	"github.com/Peanuttown/tzzGoUtil/encoding"
	"github.com/Peanuttown/tzzGoUtil/etcd"
	"github.com/Peanuttown/tzzGoUtil/output"
	"os"
	"time"
)

func main() {
	var memCfgPath string
	flag.StringVar(&memCfgPath, "members", "", "etcd 成员")
	flag.Parse()
	if len(memCfgPath) == 0 {
		fmt.Println("please input etcd members file")
	}

	// load cfg
	var members []*etcd.MemberCfg
	err := encoding.UnMarshalByFile(memCfgPath, &members, json.Unmarshal)
	if err != nil {
		output.Err(fmt.Errorf("解析 成员配置文件失败: %w", err))
		os.Exit(1)
	}
	const keySize = certs.PrivateKeyBitSize_2048
	var validDuation = time.Hour * 24 * 365
	now := time.Now()
	const caCrtPath = "tz_ca.crt"
	// << 生成 ca 证书
	caPrivKey, err := certs.RSAGenPrivateKey(keySize)
	if err != nil {
		output.Err(fmt.Sprintf("生成 etcd ca 私钥失败: %w", err))
		os.Exit(1)
	}
	caCrtTpl := certs.NewX509CaCrtTpl(
		pkix.Name{
			CommonName: "tz_etcd_ca",
		},
		now,
		validDuation,
		nil,
	)
	caCrt, err := certs.SignSelf(caCrtTpl, caPrivKey)
	if err != nil {
		output.Err(fmt.Sprintf("生成 ca 根证书失败: %w", err))
		os.Exit(1)
	}
	err = certs.PemX509Save(caCrtPath, caCrt.Raw)
	if err != nil {
		output.Err(fmt.Sprintf("保存 ca 根证书失败: %w", err))
		os.Exit(1)
	}
	// >>

	for _, mem := range members {
		// < 生成证书
		// << server crt
		certs.RSAGenPrivateKey()
		// >>
		// >
		//生成启动脚本

		// 拷贝资源
	}
}
