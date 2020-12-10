package main

import (
	"bytes"
	"html/template"
)

const pro_deploy_base_tpl = `
base:
  ntpserver: {{.Base}}
dfs:
  hosts_chunkserver: ALL
  hosts_client: ALL
  hosts_master: {{.Base}}
  hosts_metalogger: ALL
docker:
  hosts: ALL
etcd:
  hosts_client: ALL
  hosts_server: ALL
k8s:
  hosts_business: ALL
  hosts_fhmc: ALL
  hosts_gate: ALL
  hosts_logs: ALL
  hosts_master: ALL
  hosts_monitor: ALL
  hosts_node: ALL
  hosts_vip: ALL
`

type ProDeployBaseTpl struct {
	Base string
}

func (this *ProDeployBaseTpl) Marshal() ([]byte, error) {
	t, err := template.New("").Parse(pro_deploy_base_tpl)
	if err != nil {
		return nil, err
	}
	buffer := bytes.Buffer{}
	err = t.Execute(&buffer, this)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

const configbase = `
dfs:
  DFS_HDD_PATH: /data/vol_*
  HDDS:
  - /data/vol_a/dfs
  - /data/vol_b/dfs
k8s:
  CLUSTER_DOMAIN_CN: 根域
  CLUSTER_NAME: pro.cluster101
  CLUSTER_NAME_CN: 微云101
  CLUSTER_VIP: 172.16.2.85
  CNI_BUSINESS_VLANS: '172.31.0.0/16'
  KUBE_APISERVER_VIP: 172.16.2.86
  NET_PLUGIN: fhmc-cni-plugins
  REGISTRY_CA_PATH: /root/deploy/ca
usual:
  LICENSE_PATH: /root/deploy/license.fhmc
  MOUNTPOINT: /mnt/shareroot
`

//`
//k8s:
//  FEDE_TAISHI_ENDPOINT: ''
//  FHMC_CNI_BUSINESS_VLANS: ''
//  FHMC_CNI_MACVLAN: ''
//  FHMC_CNI_MACVLAN_GATEWAY: ''
//  FHMC_DATACELL_CONTROLLER_MYSQL_HOST: ''
//  FHMC_DATACELL_CONTROLLER_MYSQL_PASSWORD: ''
//  FHMC_DATACELL_CONTROLLER_MYSQL_PORT: ''
//  FHMC_DATACELL_CONTROLLER_MYSQL_USER: ''
//  FHMC_DATACELL_CONTROLLER_PG_HOST: ''
//  FHMC_DATACELL_CONTROLLER_PG_PASSWORD: ''
//  FHMC_DATACELL_CONTROLLER_PG_PORT: ''
//  FHMC_DATACELL_CONTROLLER_PG_USER: ''
//  KUBE_APISERVER_VIP: 192.168.10.101
//  NET_PLUGIN: fhmc-cni-plugins
//  REGISTRY_CA_PATH: /root/deploy/ca
//usual:
//  LICENSE_PATH: /root/deploy/license.fhmc
//  MOUNTPOINT: /mnt/shareroot
//
//`
