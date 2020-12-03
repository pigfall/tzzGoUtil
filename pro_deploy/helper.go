package main

import (
	"github.com/Peanuttown/tzzGoUtil/encoding/yaml"
)

func ConvertToYamlBytes(nodes []*ProNode) ([]byte, error) {
	var st = make(map[string]interface{}, 0)
	for _, node := range nodes {
		var unit = make(map[string]interface{})
		unit["ansible_host"] = node.Ip
		unit["ansible_ssh_pass"] = node.SSHPasswd
		unit["ansible_ssh_user"] = node.SSHUser
		st[node.Hostname] = unit
	}
	return yaml.Marshal(st)
}
