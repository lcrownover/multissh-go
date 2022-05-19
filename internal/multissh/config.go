package multissh

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

// Configure : Returns a validated SSH Config struct
func newClientConfig(p CommandParameters) *ssh.ClientConfig {
	// if password was passed, let's use that.
	if p.Password != "" {
		config, err := getSSHPasswordConfig(p)
		if err != nil {
			logger("config").Error(err.Error())
			os.Exit(1)
		}
		return config
	}

	// otherwise hope the agent has the key
	config, err := getSSHKeyConfig(p)
	if err != nil {
		logger("config").Error(err.Error())
		os.Exit(1)
	}
	return config
}

// getSSHPasswordConfig : if ssh key fails, fall back to password auth
// returns ssh.ClientConfig to be passed around to ssh sessions
func getSSHPasswordConfig(p CommandParameters) (*ssh.ClientConfig, error) {
	var config *ssh.ClientConfig
	var err error

	password := p.Password

	if password == "" {
		prompt := fmt.Sprintf("Password for user %s: ", p.Username)
		password, err = getSecurePassword(prompt)
		if err != nil {
			return config, err
		}
	}

	config = &ssh.ClientConfig{
		User: p.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	return config, nil
}

// getSSHKeyConfig : use private key to authentication without password
// returns ssh.ClientConfig to be passed around to ssh sessions
func getSSHKeyConfig(p CommandParameters) (*ssh.ClientConfig, error) {
	var config *ssh.ClientConfig
	var authMethod ssh.AuthMethod
	var conn net.Conn
	var err error

	// dial the ssh agent, conn is nil when no agent
	socket := os.Getenv("SSH_AUTH_SOCK")
	conn, err = net.Dial("unix", socket)
	if err != nil {
		logger("config").Error("failed to dial unix ssh socket")
		return config, err
	}

	// if ssh agent is detected, we don't need a private key password
	// otherwise prompt for one
	if conn != nil {
		agentClient := agent.NewClient(conn)
		authMethod = ssh.PublicKeysCallback(agentClient.Signers)
	} else {
		prompt := fmt.Sprintf("SSH key passphrase for user %s:", p.Username)
		passphrase, err := getSecurePassword(prompt)
		if err != nil {
			return config, err
		}
		authMethod, err = privateKey(p.PrivateKeyPath, passphrase)
		if err != nil {
			return config, err
		}
	}

	config = &ssh.ClientConfig{
		User:            p.Username,
		Auth:            []ssh.AuthMethod{authMethod},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	return config, nil
}

// privateKey : unlock private key and return PublicKeys object
func privateKey(path, passphrase string) (ssh.AuthMethod, error) {
	var authMethod ssh.AuthMethod
	logger("config").Debug("unlocking private key")

	key, err := ioutil.ReadFile(path)
	if err != nil {
		logger("config").Error(fmt.Sprintf("failed to read key path: %s", path))
		return authMethod, err
	}
	signer, err := ssh.ParsePrivateKeyWithPassphrase(key, []byte(passphrase))
	if err != nil {
		logger("config").Error("failed to unlock private key")
		return authMethod, err
	}
	return ssh.PublicKeys(signer), nil
}
