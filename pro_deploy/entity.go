package main

import (
	"fmt"
)

type Cfg struct {
	Original         *Node      `json:"original"`
	ProPkgPath       string     `json:"pro_pkg_path"`
	ProNodes         []*ProNode `json:"pro_nodes"`
	ClusterVip       string     `json:"cluster_vip"`
	KubeApiserverVip string     `json:"kube_apiserver_vip"`
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
	SSHPasswd string `json:"ssh_passwd"`
}

func (this *ProNode) GetAddr() string {
	return fmt.Sprintf("%s:22", this.Ip)

}
