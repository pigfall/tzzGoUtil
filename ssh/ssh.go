package ssh

import (
	"github.com/Peanuttown/tzzGoUtil/fs"
	"github.com/pkg/sftp"
	gossh "golang.org/x/crypto/ssh"
	"io"
	"os"
	"path"
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
