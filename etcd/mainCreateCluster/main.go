//----------
// AUTHOR: tzz
// DESC : 快速创建 etcd 集群
//----------
package main

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/Peanuttown/tzzGoUtil/certs"
	"github.com/Peanuttown/tzzGoUtil/encoding"
	"github.com/Peanuttown/tzzGoUtil/etcd"
	"github.com/Peanuttown/tzzGoUtil/fs"
	"github.com/Peanuttown/tzzGoUtil/output"
	"github.com/Peanuttown/tzzGoUtil/ssh"
	"net"
	"os"
	"path"
	"strings"
	tpl "text/template"
	"time"
)

const caCrtPath = "tz_ca.crt"
const serverCrtPath = "tz_server.crt"
const serverKeyPath = "tz_server.key"
const peerCrtPath = "tz_peer.crt"
const peerKeyPath = "tz_peer.key"
const etcdStartupScritPath = "etcd.sh"

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
	// << 生成 ca 证书
	caPrivKey, err := certs.RSAGenPrivateKey(keySize)
	if err != nil {
		output.Err(fmt.Sprintf("生成 etcd ca 私钥失败: %v", err))
		os.Exit(1)
	}
	//var caPubKey = &caPrivKey.PublicKey
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
		output.Err(fmt.Sprintf("生成 ca 根证书失败: %v", err))
		os.Exit(1)
	}
	err = certs.PemX509Save(caCrtPath, caCrt.Raw)
	if err != nil {
		output.Err(fmt.Sprintf("保存 ca 根证书失败: %v", err))
		os.Exit(1)
	}
	// >>

	for _, mem := range members {
		var commonName = mem.Name
		// < 生成证书
		// << server crt
		serverKey, err := certs.RSAGenPrivateKey(keySize)
		if err != nil {
			output.Err(fmt.Sprintf("生成 etcd 服务证书失败: %v", err))
			os.Exit(1)
		}
		serverCrtTpl := certs.NewX509CertTpl(
			pkix.Name{
				CommonName: commonName,
			},
			now,
			validDuation,
			func(c *x509.Certificate) {
				c.IPAddresses = []net.IP{net.ParseIP(mem.Ip)}
			},
		)
		serverCrt, err := certs.SignCrt(
			caCrt, serverCrtTpl, &serverKey.PublicKey, caPrivKey,
		)
		if err != nil {
			output.Err(fmt.Errorf("签名服务证书失败: %w", err))
			os.Exit(1)
		}
		err = certs.PemSaveRSAPrivateKey(serverKeyPath, serverKey)
		if err != nil {
			output.Err(fmt.Errorf("保存 server 私钥失败: %W", err))
			os.Exit(1)
		}
		err = certs.PemX509Save(serverCrtPath, serverCrt.Raw)
		if err != nil {
			output.Err(fmt.Errorf("保存 server 证书失败: %w", err))
			os.Exit(1)
		}
		// >>

		// << 生成 peer crt
		peerKey, err := certs.RSAGenPrivateKey(keySize)
		if err != nil {
			output.Err(fmt.Errorf("生成 etcd peer 私钥失败: %w", err))
			os.Exit(1)
		}
		peerCrtTpl := certs.NewX509CertTpl(
			pkix.Name{
				CommonName: commonName,
			},
			now,
			validDuation,
			func(c *x509.Certificate) {
				c.IPAddresses = []net.IP{net.ParseIP(mem.Ip)}
			},
		)
		peerCrt, err := certs.SignCrt(caCrt, peerCrtTpl, &peerKey.PublicKey, caPrivKey)
		if err != nil {
			output.Err(fmt.Errorf("签名 peer 证书失败: %w", err))
			os.Exit(1)
		}
		err = certs.PemSaveRSAPrivateKey(peerKeyPath, peerKey)
		if err != nil {
			output.Err(fmt.Errorf("保存 peer 私钥失败: %w", err))
			os.Exit(1)
		}
		err = certs.PemX509Save(peerCrtPath, peerCrt.Raw)
		if err != nil {
			output.Err(fmt.Errorf("保存 peer 证书失败: %w", err))
			os.Exit(1)
		}
		// >>
		// >

		// < 生成启动脚本
		var bootScriptTpl = `
#!/bin/bash
etcd --initial-cluster={{.InitialCluster}} \
--initial-advertise-peer-urls={{.InitialAdvPeerUrls}} \
--listen-client-urls={{.ListenClientUrls}} \
--listen-peer-urls={{.ListenPeerUrls}} \
--advertise-client-urls={{.AdvClientUrls}} \
--trusted-ca-file={{.TrustedCaFile}} \
--peer-trusted-ca-file={{.PeerTrustedCaFile}} \
--cert-file={{.CertFile}} \
--key-file={{.KeyFile}} \
--peer-cert-file={{.PeerCertFile}} \
--peer-key-file={{.PeerKeyFile}} \
--data-dir=data.etcd \
--name={{.Name}}
		`
		var clientPort = 2379
		var peerPort = 2380
		initialClusterList := make([]string, 0, len(members))
		for _, v := range members {
			initialClusterList = append(initialClusterList, fmt.Sprintf("%s=https://%s:%d", v.Name, v.Ip, peerPort))
		}

		argTplStruct := &EtcdStartUpArgTpl{
			InitialCluster:     strings.Join(initialClusterList, ","),
			InitialAdvPeerUrls: fmt.Sprintf("https://%s:%d", mem.Ip, peerPort),
			ListenClientUrls:   fmt.Sprintf("http://127.0.0.1:%d,https://%s:%d", clientPort, mem.Ip, clientPort),
			ListenPeerUrls:     fmt.Sprintf("https://%s:%d", mem.Ip, peerPort),
			AdvClientUrls:      fmt.Sprintf("https://%s:%d", mem.Ip, clientPort),
			TrustedCaFile:      path.Base(caCrtPath),
			PeerTrustedCaFile:  path.Base(caCrtPath),
			CertFile:           path.Base(serverCrtPath),
			KeyFile:            path.Base(serverKeyPath),
			PeerCertFile:       path.Base(peerCrtPath),
			PeerKeyFile:        path.Base(peerKeyPath),
			Name:               mem.Name,
		}
		startupTpl, err := tpl.New("startupTpl").Parse(bootScriptTpl)
		if err != nil {
			// bug
			panic(fmt.Errorf("解析模板失败: %w", err))
		}
		err = fs.CreateThen(etcdStartupScritPath, func(file *os.File) error {
			return startupTpl.Execute(file, argTplStruct)
		})
		if err != nil {
			output.Err("渲染模板失败: %w", err)
			os.Exit(1)
		}
		// >

		// < 拷贝资源
		sshClt, err := ssh.Dial(fmt.Sprintf("%s:%d", mem.Ip, 22), &ssh.DialCfg{User: mem.SSHUser, Passwd: mem.SSHPasswd})
		if err != nil {
			output.Err(fmt.Errorf("创建到 %s 的ssh客户端失败: %w ", mem.Ip, err))
			os.Exit(1)
		}
		defer sshClt.Close()
		const remoteBasePath = "/tmp/tz_etcd"
		localToSync := []string{
			caCrtPath,
			serverCrtPath,
			serverKeyPath,
			peerCrtPath,
			peerKeyPath,
			etcdStartupScritPath,
		}
		for _, sync := range localToSync {
			var remotePath = path.Join(remoteBasePath, path.Base(sync))
			err = sshClt.Copy(sync, remotePath)
			if err != nil {
				output.Err("同步 %s 到 %s 的 %s 失败: %w", sync, mem.Ip, remotePath)
				os.Exit(1)
			}
		}
		// >
	}
}

type EtcdStartUpArgTpl struct {
	InitialCluster     string
	InitialAdvPeerUrls string
	ListenClientUrls   string
	ListenPeerUrls     string
	AdvClientUrls      string
	TrustedCaFile      string
	PeerTrustedCaFile  string
	CertFile           string
	KeyFile            string
	PeerCertFile       string
	PeerKeyFile        string
	Name               string
}
