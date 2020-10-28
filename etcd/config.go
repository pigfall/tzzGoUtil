package etcd

type MemberCfg struct {
	SSHUser   string `json:"ssh_user"`
	SSHPasswd string `json:"ssh_passwd"`
	Ip        string `json:"ip"`
	Name      string `json:"name"`
}
