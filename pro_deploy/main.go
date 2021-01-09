package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/Peanuttown/tzzGoUtil/encoding/yaml"
	"github.com/Peanuttown/tzzGoUtil/fs"
	"github.com/Peanuttown/tzzGoUtil/output"
	"github.com/Peanuttown/tzzGoUtil/process"
	"github.com/Peanuttown/tzzGoUtil/ssh"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	START_FROM_DOWNLOAD_PRO_PKG   = "downloadProPkg"
	START_FROM_DOWNLOAD_CRT       = "downloadCrt"
	START_FROM_UNCOMPRESS         = "uncompress"
	START_FROM_EXPORT_CFG         = "export_cfg"
	START_FROM_EXPORT_CLUSTER_CFG = "export_cluster_cfg"
	START_FROM_CHANGE_CLUSTER_CFG = "change_cluster_cfg"
	START_FROM_INIT_CLUSTER       = "init_cluster"
)

func StepDo(startFrom string) func(from string, f func()) {
	var start bool
	return func(from string, f func()) {
		if startFrom == from || start {
			start = true
			fmt.Println("🌸 🌸 🌸  开始: ", from)
			f()
			fmt.Printf("\n")
		}
	}
}

func ShowProgress(overChan chan struct{}) {
	fmt.Println("😈 😈  开始部署")
	start := time.Now()
	tick := time.NewTicker(time.Second * 2)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			// TODO
			//fmt.Println("show progress")
		case <-overChan:
			fmt.Printf(" 🎉 🎉 🎉  pro 部署成功 ! ! ! 耗时: %v\n min ", time.Now().Sub(start).Minutes())
			return
		}
	}
}

func main() {
	var showClean bool
	var toCleans string
	var showStartFrom bool
	var verbose bool
	var cfgPath string
	var help bool
	var clean bool
	var outDemoCfg bool
	flag.StringVar(&cfgPath, "config", "config.json", "配置文件路径")
	flag.BoolVar(&outDemoCfg, "demo", false, "配置文件模板")
	var startFrom string
	flag.StringVar(&startFrom, "startFrom", START_FROM_DOWNLOAD_PRO_PKG, "从那一部开始")
	flag.BoolVar(&verbose, "verbose", false, "verbose output")
	flag.BoolVar(&showStartFrom, "showStartFrom", false, "showStartFrom")
	flag.BoolVar(&help, "help", false, "help")
	flag.BoolVar(&showClean, "showClean", false, "showClean")
	flag.BoolVar(&clean, "clean", false, "clean")
	flag.StringVar(&toCleans, "toCleans", "", "toCleans")
	flag.Parse()
	if help {
		flag.Usage()
		return
	}
	if outDemoCfg {
		fmt.Println(
			`
{
    "original":{
        "addr":"172.16.1.10:22",
        "ssh_user":"pstore",
        "ssh_passwd":"eicas2019"
    },
    
    "pro_pkg_path":"/home/pstore/packages/other/fhmc-pro/fhmc-integration/unstable/amd64/fhmc-integration_v2.4.2.STDmd64-20201127.211823-git.a39e2c.tgz",
    "pro_nodes":[
        {
            "hostname":"host82",
            "ip":"172.16.2.82",
            "ssh_user":"root",
            "ssh_passwd":"123456"
        },
        {
            "hostname":"host83",
            "ip":"172.16.2.83",
            "ssh_user":"root",
            "ssh_passwd":"123456"
        },
        {
            "hostname":"host84",
            "ip":"172.16.2.84",
            "ssh_user":"root",
            "ssh_passwd":"123456"
        }
	],
	"storage_type":"dfs",
	"cluster_vip":"172.16.3.85",
	"kube_apiserver_vip":"172.16.3.86",
	"internal_cni":"flannel"
}
			`,
		)

		return
	}
	if showStartFrom {
		fmt.Println(
			strings.Join([]string{
				START_FROM_DOWNLOAD_PRO_PKG,
				START_FROM_DOWNLOAD_CRT,
				START_FROM_UNCOMPRESS,
				START_FROM_EXPORT_CFG,
				START_FROM_EXPORT_CLUSTER_CFG,
				START_FROM_CHANGE_CLUSTER_CFG,
				START_FROM_INIT_CLUSTER,
			}, ", "),
		)
		return
	}
	// < ui
	wg := sync.WaitGroup{}
	var overChan = make(chan struct{}, 1)
	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
		}()
		ShowProgress(overChan)
	}()
	// >

	var verboseWriter, err = os.Open("/dev/null")
	handleErr(err)
	defer verboseWriter.Close()
	if verbose {
		verboseWriter = os.Stdout
	}

	doF := StepDo(startFrom)
	// < တ load config
	var cfg = &Cfg{}
	cfgRaw, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		output.Errf("read config file [%s] failed, %v\n", cfgPath, err)
		os.Exit(1)
	}
	err = json.Unmarshal(cfgRaw, cfg)
	if err != nil {
		output.Errf("unmarshal config content failed: %v", err)
		os.Exit(1)
	}
	if clean {
		err := DoCleanCluster(cfg, &cleanFlags{showClean: showClean, toCleans: strings.Split(toCleans, ",")})
		handleErr(err)
		return
	}
	if len(cfg.ProNodes) == 0 {
		handleErr(fmt.Errorf("部署节点数目为0"))
	}

	// > ဈ

	//create ssh client
	sshCli, err := ssh.Dial(cfg.Original.Addr, &ssh.DialCfg{User: cfg.Original.SSHUser, Passwd: cfg.Original.SSHPasswd})
	if err != nil {
		output.Errf("create ssh client to %s failed: %v", cfg.Original.Addr, err)
		os.Exit(1)
	}
	defer sshCli.Close()

	doF(START_FROM_DOWNLOAD_PRO_PKG, func() {
		err = DownloadProPkg(sshCli, cfg.ProPkgPath)
		handleErr(err)
	})

	doF(START_FROM_DOWNLOAD_CRT,
		func() {
			err = DownloadLicense(sshCli)
			handleErr(err)
			err = DownloadCaCrt(sshCli)
			handleErr(err)
		},
	)

	// < တ link to bash
	os.Rename("/bin/sh", "/bin/sh.back")
	os.Remove("/bin/sh")
	err = os.Link("/bin/bash", "/bin/sh")
	handleErr(err)
	// > ဈ

	// < တ uncompress
	doF(START_FROM_UNCOMPRESS, func() {
		err = process.ExecOutput(verboseWriter, os.Stderr, "tar", "-xvzf", FHMC_TAR_PATH)
		handleErr(err)
	})
	// > ဈ

	// < 解析真正的 包名
	cmd := exec.Command("tar", "-tf", FHMC_TAR_PATH)
	pipeRd, err := cmd.StdoutPipe()
	handleErr(err)
	defer pipeRd.Close()
	err = cmd.Start()
	handleErr(err)
	line, _, err := bufio.NewReader(pipeRd).ReadLine()
	handleErr(err)
	pkgDirPath := strings.Split(string(line), "/")[0]
	cmd.Process.Kill()
	// >

	//var uncompressPath = path.Join(path.Dir(FHMC_TAR_PATH), strings.TrimSuffix(path.Base(cfg.ProPkgPath), ".tgz"))
	var uncompressPath = path.Join(path.Dir(FHMC_TAR_PATH), pkgDirPath)
	var proNodeCfgPath = path.Join(uncompressPath, "fhmc-guide-deploy", "fhmc-guide.role.yaml")

	// < တ export node cfg
	doF(
		START_FROM_EXPORT_CFG,
		func() {
			switch cfg.StorageType {
			case "dfs", "gfs", "nfs": // nfs,gfs
			default:
				handleErr(fmt.Errorf("目前仅支持 dfs"))
			}
			err = process.ExecOutput(verboseWriter, os.Stderr, "bash", "-c", fmt.Sprintf("cd %s && ./fhmc-guide deploy-export --storage-type %s", path.Join(uncompressPath), cfg.StorageType))
			handleErr(err)

			cfgData, err := ConvertToYamlBytes(cfg.ProNodes)
			handleErr(err)

			cfgFile, err := os.Create(proNodeCfgPath)
			handleErr(err)

			_, err = io.Copy(cfgFile, bytes.NewReader(cfgData))
			cfgFile.Close()
			handleErr(err)

			// < change deploy base
			deployBase := &ProDeployBaseTpl{Base: cfg.ProNodes[0].Hostname}
			deployBaseData, err := deployBase.Marshal()
			if err != nil {
				handleErr(fmt.Errorf("marshal deploy base failed: %v", err))
			}
			err = ioutil.WriteFile(path.Join(uncompressPath, "fhmc-guide-deploy", "fhmc-guide.deploybase.yaml"), deployBaseData, os.ModePerm)
			handleErr(err)
			// >
		},
	)
	// > ဈ

	// < တ  export cluster cfg
	doF(
		START_FROM_EXPORT_CLUSTER_CFG,
		func() {
			err = process.ExecOutput(verboseWriter, os.Stderr, "bash", "-c", fmt.Sprintf("cd %s && %s", uncompressPath, fmt.Sprintf("./%s config-export", "fhmc-guide")))
			handleErr(err)
		},
	)
	// > ဈ

	// < တ  change cluster cfg
	// TODO
	var pathGuideConfig = path.Join(uncompressPath, "fhmc-guide-config")
	doF(
		START_FROM_CHANGE_CLUSTER_CFG,
		func() {
			if len(cfg.ClusterVip) == 0 || len(cfg.KubeApiserverVip) == 0 {
				handleErr(fmt.Errorf("cluster_vip, kube_apiserver 为空"))
			}
			configBaseTpl := &ConfigBaseTpl{ClusterVip: cfg.ClusterVip, KubeApiserverVip: cfg.KubeApiserverVip}
			data, err := configBaseTpl.Marshal()
			handleErr(err)
			err = ioutil.WriteFile(path.Join(pathGuideConfig, "fhmc-guide.configbase.yaml"), data, os.ModePerm)
			handleErr(err)

			// တ < 修改 configextra
			configExtraYamlPath := path.Join(pathGuideConfig, "fhmc-guide.configextra.yaml")
			err = fs.ReadAllThen(configExtraYamlPath, func(ct []byte) error {
				var configExtra = make(map[string]interface{})
				err := yaml.UnMarshal(ct, &configExtra)
				if err != nil {
					return err
				}

				switch cfg.InternalCNI {
				case CNI_CALICO:
					// check arch
					if runtime.GOARCH == "arm" {
						return fmt.Errorf("arm 架构支持支 flannel")
					}
				case CNI_FLANNEL:
				default:
					return fmt.Errorf("unsupport cni %s. supproted: %v", cfg.InternalCNI, []string{CNI_CALICO, CNI_FLANNEL})
				}

				if k8sCfg := configExtra["k8s"]; k8sCfg == nil {
					return fmt.Errorf("fhmc-guide.configextra.yaml 中不存在 key k8s, \n %s \n", ct)
				} else {
					k8sCfg.(map[interface{}]interface{})["FHMC_CNI_INTERNAL"] = cfg.InternalCNI
					data, err := yaml.Marshal(configExtra)
					if err != nil {
						return fmt.Errorf("marshal configextra failed")
					}
					err = ioutil.WriteFile(configExtraYamlPath, data, os.ModePerm)
					if err != nil {
						return fmt.Errorf("修改文件 %s 失败: %w", configExtraYamlPath, err)
					}
				}
				return nil
			})
			handleErr(err)
			// ဈ >
		},
	)

	// > ဈ

	// < တ init pro cluster
	doF(
		START_FROM_INIT_CLUSTER,
		func() {
			err = process.ExecOutput(verboseWriter, os.Stderr, "bash", "-c", fmt.Sprintf("cd %s && %s", uncompressPath, fmt.Sprintf("./fhmc-guide init-cluster --sure true")))
			handleErr(err)
		},
	)

	// > ဈ
	overChan <- struct{}{}
	wg.Wait()
}
