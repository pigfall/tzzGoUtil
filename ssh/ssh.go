package ssh

import (
	"fmt"
	"github.com/Peanuttown/tzzGoUtil/fs"
	"github.com/pkg/sftp"
	gossh "golang.org/x/crypto/ssh"
	"io"
	"os"
	"path"
	"strings"
)

type DialCfg struct {
	User   string
	Passwd string
}

type Client struct {
	*gossh.Client
	sftpClt *sftp.Client
}

func (c *Client) Close() {
	c.sftpClt.Close()
	c.Client.Close()
}

func Dial(addr string, cfg *DialCfg) (*Client, error) {
	clt, err := gossh.Dial("tcp", addr, &gossh.ClientConfig{
		User:            cfg.User,
		Auth:            []gossh.AuthMethod{gossh.Password(cfg.Passwd)},
		HostKeyCallback: gossh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		return nil, err
	}
	sftpClt, err := sftp.NewClient(clt)
	if err != nil {
		clt.Close()
		return nil, err
	}
	return &Client{Client: clt, sftpClt: sftpClt}, nil
}

func (c *Client) Copy(local string, remote string) error {
	c.sftpClt.MkdirAll(path.Dir(remote))
	file, err := c.sftpClt.Create(remote)
	if err != nil {
		return err
	}
	defer file.Close()
	err = fs.OpenThen(
		local,
		func(f *os.File) error {
			_, err := io.Copy(file, f)
			return err
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) DownloadFile(local, remote string) error {
	remoteFile, err := c.sftpClt.Open(remote)
	if err != nil {
		return err
	}
	defer remoteFile.Close()
	if _, err := os.Stat(path.Dir(local)); err != nil {
		os.MkdirAll(path.Dir(local), os.ModePerm)
	}
	localFile, err := os.Create(local)
	if err != nil {
		return err
	}
	defer localFile.Close()
	_, err = io.Copy(localFile, remoteFile)
	return err
}

func (c *Client) DownloadDir(local, remote string) error {
	err := os.MkdirAll(local, os.ModePerm)
	if err != nil {
		return err
	}
	walk := c.sftpClt.Walk(remote)
	for walk.Step() {
		info := walk.Stat()
		filepath := walk.Path()
		if info.IsDir() {
			err = os.MkdirAll(path.Join(local, strings.TrimPrefix(filepath, remote)), os.ModePerm)
			if err != nil {
				return err
			}
		} else {
			remoteFile, err := c.sftpClt.Open(filepath)
			if err != nil {
				return err
			}
			defer remoteFile.Close()
			file, err := os.Create(path.Join(local, strings.TrimPrefix(filepath, remote)))
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(file, remoteFile)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *Client) CreateThenWrite(filepath string, reader io.Reader) error {
	c.sftpClt.Remove(filepath)
	file, err := c.sftpClt.Create(filepath)
	if err != nil {
		err = fmt.Errorf("创建文件 %s 失败: %v", filepath, err)
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, reader)
	if err != nil {
		err = fmt.Errorf("写入文件 %s 失败: %v", filepath, err)
	}
	return err
}

func (c *Client) RemoveFile(filepath string) error {
	session, err := c.NewSession()
	if err != nil {
		return err
	}
	out, err := session.CombinedOutput(fmt.Sprintf("rm %s", filepath))
	if err != nil {
		return fmt.Errorf("删除文件 %s 失败: %v , %s", filepath, err, string(out))
	}
	return nil
}
