package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"time"
)

func ConnectSSHDirectly(targetHost SSHHost) error {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(targetHost.Host, targetHost.Port), 5*time.Second)
	if err != nil {
		return fmt.Errorf("failed to dial directly: %w", err)
	}
	defer conn.Close()

	errChan := make(chan error, 2)

	go func() {
		_, err := io.Copy(conn, os.Stdin)
		errChan <- err
	}()
	go func() {
		_, err := io.Copy(os.Stdout, conn)
		errChan <- err
	}()

	if err := <-errChan; err != nil && err != io.EOF {
		return fmt.Errorf("error during direct connection (stdin->conn): %w", err)
	}
	if err := <-errChan; err != nil && err != io.EOF {
		return fmt.Errorf("error during direct connection (conn->stdout): %w", err)
	}

	return nil
}

func ConnectSSHWithCommand(jumpHost SSHHost, targetHost SSHHost) error {
	sshArgs := []string{
		"-T",
		"-W", fmt.Sprintf("%s:%s", targetHost.Host, targetHost.Port),
		"-p", jumpHost.Port,
		fmt.Sprintf("%s@%s", jumpHost.User, jumpHost.Host),
	}

	cmd := exec.Command("ssh", sshArgs...)

	log.Printf("Connecting to %s...\n", targetHost.Host)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Printf("Error executing ssh command: %v\n", err)
		os.Exit(1)
	}
	return nil
}
