package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

func createSSH(config Config) {
	fmt.Println("Connection to " + config.Alias + "...")

	home, err := homedir.Dir()
	check(err)
	key, err := os.ReadFile(strings.Replace(config.IdentityFile, "~", home, 1))
	if err != nil {
		log.Fatalf("Unable to read private key: %v", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("Unable to parse private key: %v", err)
	}

	sshConfig := &ssh.ClientConfig{
		User:            config.User,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // TODO improve it
	}

	if config.Port == 0 {
		config.Port = 22
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", config.Hostname, config.Port), sshConfig)
	if err != nil {
		log.Fatalf("Failed to dial: %s", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		log.Fatalf("Failed to create session: %s", err)
	}
	defer session.Close()

	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		log.Fatalf("Failed to set terminal raw mode: %s", err)
	}
	defer func() {
		if err := term.Restore(fd, oldState); err != nil {
			log.Printf("Failed to restore terminal: %v", err)
		}
	}()

	if err := session.RequestPty("xterm", 80, 40, ssh.TerminalModes{}); err != nil {
		log.Fatalf("Failed to request pseudo-terminal: %s", err)
	}

	session.Stdin = os.Stdin
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	if err := session.Shell(); err != nil {
		log.Fatalf("Failed to start shell: %s", err)
	}

	if err := session.Wait(); err != nil {
		log.Fatalf("Failed to wait shell: %s", err)
	}

	fmt.Println()
}
