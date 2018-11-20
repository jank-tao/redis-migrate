package internal

import (
	"io/ioutil"
	"net"

	"golang.org/x/crypto/ssh"
	"gopkg.in/redis.v5"
)

type (
	sshConfig struct {
		Addr           string `json:"addr"`
		PrivateKeyPath string `json:"private_key_path"`
		Username       string `json:"username"`
	}

	RedisConfig struct {
		Addr      string     `json:"addr"`
		Password  string     `json:"password"`
		DB        int        `json:"db"`
		SSHConfig *sshConfig `json:"ssh_config"`
	}

	TaskConfig struct {
		From     string   `json:"from"`
		To       string   `json:"to"`
		Patterns []string `json:"patterns"`
	}

	MigrateConfig struct {
		Redis map[string]RedisConfig `json:"redis"`
		Tasks map[string]TaskConfig  `json:"tasks"`
	}
)

func (c RedisConfig) Client() (*redis.Client, error) {
	if c.SSHConfig == nil {
		return c.newClient()
	}
	return c.newSSHClient()
}

func (c RedisConfig) newClient() (*redis.Client, error) {
	rc := redis.NewClient(
		&redis.Options{
			Addr:     c.Addr,
			Password: c.Password,
			DB:       c.DB,
		},
	)
	if err := rc.Ping().Err(); err != nil {
		return nil, err
	}
	return rc, nil
}

func (c RedisConfig) newSSHClient() (*redis.Client, error) {
	if c.SSHConfig.PrivateKeyPath == "" {
	}
	privateKey, err := ioutil.ReadFile(c.SSHConfig.PrivateKeyPath)
	if err != nil {
		return nil, err
	}
	signer, err := ssh.ParsePrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	rc := redis.NewClient(&redis.Options{
		Password: c.Password,
		DB:       c.DB,
		Dialer: func() (net.Conn, error) {
			conn, err := net.Dial("tcp", c.SSHConfig.Addr)
			if err != nil {
				return nil, err
			}
			sshConn, chans, reqs, err := ssh.NewClientConn(conn, c.SSHConfig.Addr, &ssh.ClientConfig{
				User: c.SSHConfig.Username,
				Auth: []ssh.AuthMethod{
					ssh.PublicKeys(signer),
				},
				HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			})
			if err != nil {
				return nil, err
			}
			client := ssh.NewClient(sshConn, chans, reqs)
			return client.Dial("tcp", c.Addr)
		},
	})
	if err := rc.Ping().Err(); err != nil {
		return nil, err
	}

	return rc, nil
}
