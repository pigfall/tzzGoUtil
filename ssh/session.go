package ssh

import (
	//"bytes"
	"github.com/pkg/sftp"
	gossh "golang.org/x/crypto/ssh"
)

type Session struct {
	*gossh.Session
	sftpClt *sftp.Client
	clt     *Client
}

func newSession(s *gossh.Session, sftpClt *sftp.Client) *Session {
	return &Session{
		Session: s,
		sftpClt: sftpClt,
	}
}
