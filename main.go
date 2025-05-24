package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	log.SetPrefix("[UEC-SSH-PROXY] ")
	log.SetFlags(log.Lmsgprefix)

	host := flag.String("host", "", "Host to connect to")
	user := flag.String("user", "", "User to connect as")
	port := flag.String("port", "22", "Port to connect to (default: 22)")

	flag.Parse()
	if *host == "" || *user == "" {
		flag.Usage()
		os.Exit(1)
	}

	jumpHost := SSHHost{
		Host: "ssh.cc.uec.ac.jp",
		User: *user,
		Port: "22",
	}
	targetHost := SSHHost{
		Host: *host,
		User: *user,
		Port: *port,
	}

	var actualHost string
	var err error

	switch {
	case strings.Contains(strings.ToLower(targetHost.Host), "ced"):
		log.Println("Checking for available hosts in CED...")
		actualHost, err = GetCEDHost(jumpHost) // CED IDを取得する関数を呼び出す
	case strings.Contains(strings.ToLower(targetHost.Host), "ied"):
		log.Println("Checking for available hosts in IED...")
		actualHost, err = GetIEDHost(jumpHost) // PC IDを取得する関数を呼び出す
	case strings.Contains(strings.ToLower(targetHost.Host), "sol"):
		actualHost = "sol.edu.cc.uec.ac.jp"
	case strings.Contains(strings.ToLower(targetHost.Host), "ssh"):
		actualHost = "ssh.cc.uec.ac.jp"
	}
	if err != nil {
		log.Printf("Error checking available hosts: %v\n", err)
		os.Exit(1)
	}

	if socksProxy != nil {
		socksProxy.Close()
	}

	targetHost.Host = actualHost

	err = CheckPortAvailability(fmt.Sprintf("%s:%s", targetHost.Host, targetHost.Port), 100*time.Millisecond, 500*time.Millisecond)
	requireSSHProxy := false
	if err != nil {
		requireSSHProxy = true
	}

	if requireSSHProxy {
		log.Printf("Connecting to %s via %s...\n", targetHost.Host, jumpHost.Host)
		err := ConnectSSHWithCommand(jumpHost, targetHost)
		if err != nil {
			log.Printf("Error connecting to host: %v\n", err)
			os.Exit(1)
		}
	} else {
		log.Printf("Connecting directly to %s...\n", targetHost.Host)
		err := ConnectSSHDirectly(targetHost)
		if err != nil {
			log.Printf("Error connecting to host: %v\n", err)
			os.Exit(1)
		}
	}
}
