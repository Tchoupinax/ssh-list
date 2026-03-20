package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	homedir "github.com/mitchellh/go-homedir"
	"golang.org/x/crypto/ssh"
)

// dialSSHClient opens an SSH connection using the same rules as createSSH (key file, default port 22).
func dialSSHClient(config Config) (*ssh.Client, error) {
	home, err := homedir.Dir()
	if err != nil {
		return nil, err
	}

	keyPath := strings.Replace(config.IdentityFile, "~", home, 1)
	key, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("read identity file: %w", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("parse private key: %w", err)
	}

	port := config.Port
	if port == 0 {
		port = 22
	}

	sshConfig := &ssh.ClientConfig{
		User:            config.User,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // TODO: improve host key verification
		Timeout:         12 * time.Second,
	}

	addr := fmt.Sprintf("%s:%d", config.Hostname, port)
	return ssh.Dial("tcp", addr, sshConfig)
}
