package cmd

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

func createSSH(config Config) {
	fmt.Println("Connection to " + config.Alias + "...")

	client, err := dialSSHClient(config)
	if err != nil {
		log.Fatalf("Failed to dial: %s", err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			log.Printf("Failed to close session: %v", err)
		}
	}()

	session, err := client.NewSession()
	if err != nil {
		log.Fatalf("Failed to create session: %s", err)
	}
	defer func() {
		if err := session.Close(); err != nil {
			log.Printf("Failed to close session: %v", err)
		}
	}()

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
