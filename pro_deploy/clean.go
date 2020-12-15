package main

import (
	"fmt"
	"github.com/Peanuttown/tzzGoUtil/ssh"
	"path"
)

const (
	ETCD    = "etcd"
	DOCKER  = "docker"
	RAMMER  = "rammer"
	EIKUBE  = "eikube"
	DRACOFS = "dracofs"
)

var validItems = []string{
	ETCD,
	DOCKER,
	RAMMER,
	EIKUBE,
	DRACOFS,
}

type cleanFlags struct {
	showClean bool
	toCleans  []string
}

func DoCleanCluster(cfg *Cfg, cleanFlags *cleanFlags) error {
	if cleanFlags.showClean {
		fmt.Printf("%v\n", validItems)
		return nil
	}
	for _, node := range cfg.ProNodes {
		err := cleanNode(node, cleanFlags.toCleans)
		if err != nil {
			return err
		}
	}
	return nil
}

func cleanNode(node *ProNode, toCleans []string) error {
	sshCli, err := ssh.Dial(node.GetAddr(), &ssh.DialCfg{
		User:   node.SSHUser,
		Passwd: node.SSHPasswd,
	})
	if err != nil {
		return err
	}
	defer sshCli.Close()

	cleanF := make(map[string]func(*ssh.Client) error)
	cleanF[ETCD] = cleanEtcd
	cleanF[DOCKER] = cleanDocker
	cleanF[RAMMER] = cleanRammer
	cleanF[EIKUBE] = cleanEikube
	cleanF[DRACOFS] = cleanDracofs

	for _, v := range toCleans {
		if v == "all" {
			for _, f := range cleanF {
				err = f(sshCli)
				if err != nil {
					return err
				}
			}
			return nil
		}

		if f := cleanF[v]; f == nil {
			return fmt.Errorf("有效的清理选项为 %v", validItems)
		} else {
			err = f(sshCli)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func cleanEtcd(session *ssh.Client) error {
	fmt.Println("开始清理  etcd ")
	const uninstallEtcd = "/usr/local/eietcd/uninstall"
	session.ExecIfExist(uninstallEtcd, func() error {
		return session.RunWithErrOut(fmt.Sprintf("sh %s", uninstallEtcd))
	})

	fmt.Println("unlink etcd")
	err := session.UnlinkIfExist("/usr/sbin/eietcd")
	if err != nil {
		return err
	}

	fmt.Println("remove etcd dir")
	err = session.ForceRemove(path.Dir(uninstallEtcd))
	if err != nil {
		return err
	}

	fmt.Println("remove etcd data dir")
	err = session.ForceRemove(fmt.Sprintf("/var/lib/etcd/"))
	if err != nil {
		return err
	}
	session.ForceRemove("/etc/eietcd")
	return nil
}

func cleanDocker(session *ssh.Client) error {
	fmt.Println("开始清理 docker ")
	const uninstallDocker = "/usr/local/eidocker/uninstall"
	fmt.Println("执行 uninstall ")
	session.ExecIfExist(uninstallDocker, func() error {
		return session.RunWithErrOut(fmt.Sprintf("sh %s", uninstallDocker))
	})

	fmt.Println("删除 eidocker dir")
	err := session.ForceRemove(path.Dir(uninstallDocker))
	if err != nil {
		return err
	}

	fmt.Println("unlink eidocker")
	err = session.UnlinkIfExist("/usr/sbin/eidocker")
	if err != nil {
		return err
	}

	session.ForceRemove("/etc/eidocker")
	return nil
}

func cleanRammer(session *ssh.Client) error {
	const uninstallRammer = "/usr/local/rammer/uninstall"
	session.ExecIfExist(uninstallRammer, func() error {
		return session.RunWithErrOut(fmt.Sprintf("sh %s", uninstallRammer))
	})

	session.ForceRemove(path.Dir(uninstallRammer))

	return session.UnlinkIfExist("/usr/sbin/rammer")
}

func cleanEikube(session *ssh.Client) error {
	// <
	session.RunWithErrOut("eikube --sure=true deploy-clean")
	// >
	session.ForceRemove("/etc/eikube")
	session.ExecIfExist("/usr/local/eikube/uninstall", func() error {
		return session.RunWithErrOut(fmt.Sprintf("sh %s", "/usr/local/eikube/uninstall"))
	})
	return nil
}

func cleanDracofs(session *ssh.Client) error {
	const uninstallPath = "/usr/local/dracofs-adm"
	session.ExecIfExist(
		uninstallPath,
		func() error {
			return session.RunWithErrOut(fmt.Sprintf("sh %s", uninstallPath))
		},
	)
	err := session.ForceRemove((uninstallPath))
	if err != nil {
		return err
	}

	session.ForceRemove("/etc/dracofs-adm")
	session.ForceRemove("/etc/dracofs")
	session.ForceRemove("/usr/sbin/dracofs-adm")
	return nil
}
