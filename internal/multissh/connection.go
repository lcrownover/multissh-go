package multissh

import (
	"fmt"
	"golang.org/x/crypto/ssh"
)

type connection struct {
	host       string
	connection *ssh.Client
}

func Connect(host string, config *ssh.ClientConfig) (connection, error) {
	logger("connection").Debug(fmt.Sprintf("%s: connecting...", host))
	connectHost := fmt.Sprintf("%s:22", host)
	conn, err := ssh.Dial("tcp", connectHost, config)
	if err != nil {
		logger("connection").Error(fmt.Sprintf("%s: failed to connect", host))
        return connection{}, err
	}
	logger("connection").Debug(fmt.Sprintf("%s: connected", host))
    c := connection{
        host: host,
        connection: conn,
    }
	return c, nil
}
