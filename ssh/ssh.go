package ssh

import (
	"bytes"
	"errors"
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
	cfg     *DialCfg
	addr    string
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
		err = fmt.Errorf("dial ssh addr :%s 失败: %w", fmt.Sprintf("%s@%s", cfg.User, addr), err)
		return nil, err
	}
	sftpClt, err := sftp.NewClient(clt)
	if err != nil {
		clt.Close()
		return nil, err
	}
	return &Client{
			Client: clt, sftpClt: sftpClt,
			cfg:  cfg,
			addr: addr,
		},
		nil
}

func (c *Client) Copy(local string, remote string) error {
	c.sftpClt.MkdirAll(path.Dir(remote))
	if exist, err := c.PathExist(remote); err != nil {
		return err
	} else if exist {
		err = c.ForceRemove(remote)
		if err != nil {
			return err
		}
	}
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
	return c.do(func(session *Session) error {
		out, err := session.CombinedOutput(fmt.Sprintf("rm %s", filepath))
		if err != nil {
			return fmt.Errorf("删除文件 %s 失败: %v , %s", filepath, err, string(out))
		}
		return nil
	})
}

func (c *Client) do(f func(*Session) error) error {
	session, err := c.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()
	return f(session)
}

func (c *Client) NewSession() (*Session, error) {
	session, err := c.Client.NewSession()
	if err != nil {
		return nil, fmt.Errorf("New session failed:%w", err)
	}
	return newSession(session, c.sftpClt), nil
}

func (this *Client) Unlink(path string) error {
	return this.do(
		func(s *Session) error {
			out, err := s.CombinedOutput(fmt.Sprintf("unlink %s", path))
			if err != nil {
				return fmt.Errorf("unlink 文件 %s 失败: %w", path, fmt.Errorf("%s,%w", out, err))
			}
			return nil
		},
	)
}

func (this *Client) UnlinkIfExist(path string) error {
	return this.do(
		func(s *Session) error {
			_, err := this.sftpClt.Lstat(path)
			if err != nil {
				if errors.Is(err, os.ErrNotExist) {
					return nil
				}
				return err
			}
			return this.Unlink(path)
		},
	)
}
func (this *Client) PathExist(p string) (exist bool, err error) {
	_, err = this.sftpClt.Stat(p)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (this *Client) ExecIfExist(p string, f func() error) (err error) {
	e, err := this.PathExist(p)
	if err != nil {
		return err
	}
	if e {
		return f()
	}
	return nil
}

func (this *Client) RunWithErrOut(cmd string) (err error) {
	return this.do(
		func(s *Session) error {
			b := bytes.Buffer{}
			s.Stderr = &b
			err = s.Run(cmd)
			if err != nil {
				return err
				//return fmt.Errorf("执行命令 %s 失败: %w", cmd, fmt.Errorf("%s,%w", b.String(), err))
			}
			return nil
		},
	)
}

func (this *Client) ForceRemove(p string) error {
	return this.RunWithErrOut(fmt.Sprintf("rm -rf %s", p))
}
