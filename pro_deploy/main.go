package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/Peanuttown/tzzGoUtil/encoding/yaml"
	"github.com/Peanuttown/tzzGoUtil/output"
	"github.com/Peanuttown/tzzGoUtil/process"
	"github.com/Peanuttown/tzzGoUtil/ssh"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
)

const (
	START_FROM_DOWNLOAD_PRO_PKG = "downloadProPkg"
	START_FROM_DOWNLOAD_CRT     = "downloadCrt"
)

type Cfg struct {
	Original   *Node      `json:"original"`
	ProPkgPath string     `json:"pro_pkg_path"`
	ProNodes   []*ProNode `json:"pro_nodes"`
}

type Node struct {
	Addr      string `json:"addr"`
	SSHUser   string `json:"ssh_user"`
	SSHPasswd string `json:"ssh_passwd"`
}

type ProNode struct {
	Hostname  string `json:"hostname"`
	Ip        string `json:"ip"`
	SSHUser   string `json:"ssh_user"`
	SSHPasswd string `ssh_passwd`
}

func handleErr(err error) {
	if err != nil {
		output.Err(err)
		output.Err("❌ ❌  部署失败")
		os.Exit(1)
	}
}

func main() {
	var start bool
	var cfgPath string
	flag.StringVar(&cfgPath, "config", "config.json", "配置文件路径")
	var startFrom string
	flag.StringVar(&startFrom, "startFrom", START_FROM_DOWNLOAD_PRO_PKG, "从那一部开始")
	flag.Parse()
	var cfg = &Cfg{}
	cfgRaw, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		output.Err("read config file [%s] failed, %v", err)
		os.Exit(1)
	}
	err = json.Unmarshal(cfgRaw, cfg)
	if err != nil {
		output.Errf("unmarshal config content failed: %v", err)
		os.Exit(1)
	}
	//create ssh client
	sshCli, err := ssh.Dial(cfg.Original.Addr, &ssh.DialCfg{User: cfg.Original.SSHUser, Passwd: cfg.Original.SSHPasswd})
	if err != nil {
		output.Errf("create ssh client to %s failed: %v", cfg.Original.Addr, err)
		os.Exit(1)
	}
	defer sshCli.Close()

	if startFrom == START_FROM_DOWNLOAD_PRO_PKG {
		start = true
		err = DownloadProPkg(sshCli, cfg.ProPkgPath)
		handleErr(err)
	}

	if startFrom == START_FROM_DOWNLOAD_CRT || start {
		start = true
		err = DownloadLicense(sshCli)
		handleErr(err)
		err = DownloadCaCrt(sshCli)
		handleErr(err)
	}

	// < တ link to bash
	os.Rename("/bin/sh", "/bin/sh.back")
	os.Remove("/bin/sh")
	err = os.Link("/bin/bash", "/bin/sh")
	handleErr(err)
	// > ဈ

	// < တ uncompress

	err = process.ExecWithErrOutput(os.Stderr, "tar", "-xf", FHMC_TAR_PATH)
	handleErr(err)
	// > ဈ

	var uncompressPath = path.Join(path.Dir(FHMC_TAR_PATH), path.Base(cfg.ProPkgPath))

	err = process.ExecWithErrOutput(os.Stderr, "bash", "-c", fmt.Sprintf("cd %s && ./fhmc-guide deploy-export --storage-type gfs", path.Join(uncompressPath)))
	handleErr(err)

	// < config pro node
	cfgData, err := ConvertToYamlBytes(cfg.ProNodes)
	handleErr(err)
	sshCli.
	// >
}

func ConvertToYamlBytes(nodes []*ProNode) ([]byte, error) {
	var units = make([]map[string]interface{}, 0)
	for _, node := range nodes {
		var unit = make(map[string]interface{})
		subMap := make(map[string]interface{})
		subMap["ansible_host"] = node.Ip
		subMap["ansible_ssh_pass"] = node.SSHPasswd
		subMap["ansible_ssh_user"] = node.SSHUser
		unit[node.Hostname] = subMap
		units = append(units, unit)
	}
	return yaml.Marshal(units)
}
