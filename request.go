package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"slices"
	"time"
)

var needSSHProxyList = []string{}

func execHttpRequestWithFallback(jumpHost SSHHost, req http.Request) (string, error) {
	if slices.Contains(needSSHProxyList, req.URL.Host) {
		return execHttpRequestOverSSH(jumpHost, req)
	}

	httpClient := http.Client{
		Timeout: 3 * time.Second,
	}
	res, err := httpClient.Do(&req)
	if err != nil {
		needSSHProxyList = append(needSSHProxyList, req.URL.Host)
		return execHttpRequestOverSSH(jumpHost, req)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		needSSHProxyList = append(needSSHProxyList, req.URL.Host)
		return execHttpRequestOverSSH(jumpHost, req)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}
	return string(body), nil
}

var socksProxy *SOCKSProxy

func execHttpRequestOverSSH(jumpHost SSHHost, req http.Request) (string, error) {
	if socksProxy == nil {
		var err error
		socksProxy, err = NewSOCKSProxy(jumpHost)
		if err != nil {
			log.Fatalf("error creating SOCKS proxy: %v", err)
			os.Exit(1)
		}
	}

	res, err := socksProxy.http.Do(&req)
	if err != nil {
		return "", fmt.Errorf("error executing request over SOCKS proxy: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("request failed with status code: %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}
	return string(body), nil
}
