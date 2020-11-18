package main

import (
	"fmt"
	"github.com/Peanuttown/tzzGoUtil/ssh"
)

const (
	LICENSE_PATH  = "/root/deploy/eicas-license"
	CRT_PATH      = "/root/deploy/eicas-ca"
	FHMC_TAR_PATH = "/root/fhmc-pro.tgz"
)

func DownloadProPkg(sshCli *ssh.Client, pkgPath string) error {
	fmt.Printf("开始从 %s 下载 %s 到 %s\n", sshCli.RemoteAddr().String(), pkgPath, "/root")
	return sshCli.DownloadFile(FHMC_TAR_PATH, pkgPath)
}

func DownloadLicense(sshCli *ssh.Client) error {
	fmt.Printf("开始下载 许可证\n")
	return sshCli.DownloadDir(LICENSE_PATH, "/home/pstore/packages/other/common/certs/eicas-license")
<<<<<<< HEAD

=======
>>>>>>> update
}
func DownloadCaCrt(sshCli *ssh.Client) error {
	fmt.Printf("开始下载 证书 \n")
	return sshCli.DownloadDir(CRT_PATH, "/home/pstore/packages/other/common/certs/eicas-ca")
}
