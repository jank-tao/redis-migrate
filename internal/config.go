package internal

import (
	"io/ioutil"
	"net"

	"github.com/gomodule/redigo/redis"
	"golang.org/x/crypto/ssh"
)

type (
	sshConfig struct {
		Addr           string `json:"addr"`
		Username       string `json:"username"`
		PrivateKeyPath string `json:"private_key_path"`
		Password       string `json:"password"`
	}

	redisConfig struct {
		Addr      string     `json:"addr"`
		Password  string     `json:"password"`
		DB        int        `json:"db"`
		SSHConfig *sshConfig `json:"ssh_config"`
	}

	migrateOption struct {
		IgnoreTTL bool
	}

	taskConfig struct {
		From     string   `json:"from"`
		To       string   `json:"to"`
		Patterns []string `json:"patterns"`
		migrateOption
	}

	MigrateConfig struct {
		Redis map[string]redisConfig `json:"redis"`
		Tasks map[string]taskConfig  `json:"tasks"`
	}
)

func (c redisConfig) Client() (redis.Conn, error) {
	if c.SSHConfig == nil {
		return c.newClient()
	}
	return c.newSSHClient()
}

func (c redisConfig) newClient() (redis.Conn, error) {
	return redis.Dial("tcp", c.Addr, redis.DialPassword(c.Password), redis.DialDatabase(c.DB))
}

func (c redisConfig) newSSHClient() (redis.Conn, error) {
	privateKey, err := ioutil.ReadFile(c.SSHConfig.PrivateKeyPath)
	if err != nil {
		return nil, err
	}
	signer, err := ssh.ParsePrivateKey(privateKey)
	if err != nil {
		return nil, err
	}

	return redis.Dial("tcp", c.Addr,
		redis.DialPassword(c.Password),
		redis.DialDatabase(c.DB),
		redis.DialNetDial(func(network, addr string) (net.Conn, error) {
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
			return client.Dial(network, addr)
		}))
}

