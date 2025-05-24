package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"time"

	"golang.org/x/net/proxy"
)

var httpClient *http.Client
var sshCmd *exec.Cmd

func initSOCKSProxiedHttpClient(jumpHost SSHHost) (*http.Client, *exec.Cmd, error) {
	if httpClient != nil {
		return httpClient, sshCmd, nil
	}
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, nil, fmt.Errorf("could not start SOCKS proxy listener: %w", err)
	}
	localSocksPort := listener.Addr().(*net.TCPAddr).Port
	listener.Close()

	localSocksAddr := net.JoinHostPort("127.0.0.1", fmt.Sprintf("%d", localSocksPort))

	sshCmd = exec.Command("ssh",
		"-D", localSocksAddr,
		"-N",
		"-f",
		"-o", "StrictHostKeyChecking=no",
		"-o", "ExitOnForwardFailure=yes",
		fmt.Sprintf("%s@%s", jumpHost.User, jumpHost.Host),
		"-p", jumpHost.Port,
	)

	sshCmd.Stdin = nil
	sshCmd.Stdout = os.Stderr
	sshCmd.Stderr = os.Stderr

	if err := sshCmd.Start(); err != nil {
		return nil, nil, fmt.Errorf("could not start SSH command: %w", err)
	}

	processDone := make(chan error, 1)
	go func() {
		processDone <- sshCmd.Wait()
	}()

	if err := CheckPortAvailability(localSocksAddr, 100*time.Millisecond, 10*time.Second); err != nil {
		select {
		case sshErr := <-processDone:
			return nil, nil, fmt.Errorf("SSH process exited before port became available: %w", sshErr)
		default:
			return nil, nil, fmt.Errorf("SOCKS proxy port did not become available within timeout: %v", err)
		}
	}

	socksDialer, err := proxy.SOCKS5("tcp", localSocksAddr, nil, proxy.Direct)
	if err != nil {
		return nil, nil, fmt.Errorf("could not create SOCKS5 dialer: %w", err)
	}

	httpClient = &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return socksDialer.Dial(network, addr)
			},
		},
	}

	return httpClient, sshCmd, nil
}

type SOCKSProxy struct {
	http *http.Client
	ssh  *exec.Cmd
}

func NewSOCKSProxy(jumpHost SSHHost) (*SOCKSProxy, error) {
	httpClient, sshCmd, err := initSOCKSProxiedHttpClient(jumpHost)
	if err != nil {
		return nil, fmt.Errorf("could not initialize SOCKS proxy: %w", err)
	}
	return &SOCKSProxy{
		http: httpClient,
		ssh:  sshCmd,
	}, nil
}

func (p *SOCKSProxy) Close() error {
	if p.ssh != nil && p.ssh.Process != nil {
		if err := p.ssh.Process.Kill(); err != nil {
			return fmt.Errorf("could not kill SSH process: %w", err)
		}
	}
	return nil
}
