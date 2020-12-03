package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/Peanuttown/tzzGoUtil/output"
	"github.com/Peanuttown/tzzGoUtil/process"
	"github.com/Peanuttown/tzzGoUtil/ssh"
	"io"
	"io/ioutil"
	"os"
	"path"
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
	var showStartFrom bool
	var verbose bool
	var cfgPath string
	var help bool
	flag.StringVar(&cfgPath, "config", "config.json", "配置文件路径")
	var startFrom string
	flag.StringVar(&startFrom, "startFrom", START_FROM_DOWNLOAD_PRO_PKG, "从那一部开始")
	flag.BoolVar(&verbose, "verbose", false, "verbose output")
	flag.BoolVar(&showStartFrom, "showStartFrom", false, "showStartFrom")
	flag.BoolVar(&help, "help", false, "help")
	flag.Parse()
	if help {
		flag.Usage()
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
		output.Err("read config file [%s] failed, %v", err)
		os.Exit(1)
	}
	err = json.Unmarshal(cfgRaw, cfg)
	if err != nil {
		output.Errf("unmarshal config content failed: %v", err)
		os.Exit(1)
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

	var uncompressPath = path.Join(path.Dir(FHMC_TAR_PATH), strings.TrimSuffix(path.Base(cfg.ProPkgPath), ".tgz"))
	var proNodeCfgPath = path.Join(uncompressPath, "fhmc-guide-deploy", "fhmc-guide.role.yaml")

	// < တ export node cfg
	doF(
		START_FROM_EXPORT_CFG,
		func() {
			err = process.ExecOutput(verboseWriter, os.Stderr, "bash", "-c", fmt.Sprintf("cd %s && ./fhmc-guide deploy-export --storage-type dfs", path.Join(uncompressPath)))
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
	doF(
		START_FROM_CHANGE_CLUSTER_CFG,
		func() {
			err = ioutil.WriteFile(path.Join(uncompressPath, "fhmc-guide-config", "fhmc-guide.configbase.yaml"), []byte(configbase), os.ModePerm)
			handleErr(err)
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
