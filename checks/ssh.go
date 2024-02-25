package checks

import (
	"errors"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/ssh"
	"log"
	"monitoring/spec"
	"net"
	"time"
)

func checkSSH(url string, checkSpec *spec.CheckSpec) error {
	if nil == checkSpec.DSNParams {
		return errors.New("DSNParams user and password must be set for SSH check!")
	}

	user := (*checkSpec.DSNParams)["user"].(string)
	pass := (*checkSpec.DSNParams)["pass"].(string)

	config := &ssh.ClientConfig{
		User:    user,
		Timeout: time.Duration(checkSpec.Timeout) * time.Second,
		Auth: []ssh.AuthMethod{
			ssh.Password(pass),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	conn, err := ssh.Dial("tcp", url, config)
	if nil != err {
		log.Println(err)
		return spec.ConnectionFailed
	}
	defer conn.Close()

	sess, err := conn.NewSession()
	if nil != err {
		return spec.ConnectionFailed
	}
	defer sess.Close()

	return nil
}

func (c *CheckHandler) CheckSSH(spec *spec.CheckSpec) (*map[string]string, error) {
	return nil, checkSSH(spec.DSN, spec)
}
