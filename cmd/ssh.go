package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

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
			log.Printf("Failed to close client: %v", err)
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

	restore := func() error {
		return term.Restore(fd, oldState)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	// Important: os.Exit (including log.Fatal) skips defers — we must Restore before any Exit after MakeRaw.
	defer func() {
		if err := restore(); err != nil {
			log.Printf("Failed to restore terminal: %v", err)
		}
	}()

	// Ctrl+C delivered to ssh-list restores the local TTY instead of terminating without restoring.
	go func() {
		for range sigCh {
			_ = restore()
			fmt.Fprintln(os.Stderr)
			os.Exit(130)
		}
	}()

	if err := session.RequestPty("xterm", 80, 40, ssh.TerminalModes{}); err != nil {
		log.Printf("Failed to request pseudo-terminal: %v", err)
		return
	}

	session.Stdin = os.Stdin
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	if err := session.Shell(); err != nil {
		log.Printf("Failed to start shell: %v", err)
		return
	}

	if err := session.Wait(); err != nil {
		fmt.Println()
		log.Printf("SSH session ended: %v", err)
		return
	}

	fmt.Println()
}
