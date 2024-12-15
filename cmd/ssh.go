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

	// Parse the private key
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("Unable to parse private key: %v", err)
	}

	// Create the SSH client configuration
	sshConfig := &ssh.ClientConfig{
		User: config.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer), // Authenticate using the private key
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Disable host key check (not recommended for production)
	}

	// If there is no port specified consider it's default
	if config.Port == 0 {
		config.Port = 22
	}

	// Establish connection to the SSH server
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", config.Hostname, config.Port), sshConfig)
	if err != nil {
		log.Fatalf("Failed to dial: %s", err)
	}
	defer client.Close()

	// Open a new session
	session, err := client.NewSession()
	if err != nil {
		log.Fatalf("Failed to create session: %s", err)
	}
	defer session.Close()

	// Request a pseudo-terminal
	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		log.Fatalf("Failed to set terminal raw mode: %s", err)
	}
	defer term.Restore(fd, oldState)

	if err := session.RequestPty("xterm", 80, 40, ssh.TerminalModes{}); err != nil {
		log.Fatalf("Failed to request pseudo-terminal: %s", err)
	}

	// Set up input/output redirection
	session.Stdin = os.Stdin
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	if err := session.Shell(); err != nil {
		log.Fatalf("Failed to start shell: %s", err)
	}

	if err := session.Wait(); err != nil {
		log.Fatalf("Session ended with error: %s", err)
	}
}
