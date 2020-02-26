package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"gopkg.in/ldap.v2"
)

var (
	hostname string
	port = 389

	useTLS bool

	bindUser string
	bindPass string
)

func main() {
	flag.StringVar(&hostname, "h", hostname, "LDAP Hostname")
	flag.IntVar(&port, "p", port, "LDAP Port")
	flag.BoolVar(&useTLS, "tls", useTLS, "LDAP use tls")
	flag.StringVar(&bindUser, "bindUser", bindUser, "LDAP Bind Username")
	flag.StringVar(&bindPass, "bindPass", bindPass, "LDAP Bind Password")
	flag.Parse()

	addr := net.JoinHostPort(hostname, strconv.Itoa(port))

	fmt.Printf("\nGot options:\n\tAddr: %s\n\tUse TLS: %t\n\tBind User: %s\n\tBind Pass: %s\n\n", addr, useTLS, bindUser, bindPass)

	fmt.Println("Attempting to dial the LDAP server...")
	l, err := ldap.Dial("tcp", addr)
	if err != nil {
		fmt.Printf("Failed to dial the LDAP server: %v\n", err)
		os.Exit(1)
	}
	defer l.Close()

	fmt.Println("✓ Connected to LDAP Server")

	// Reconnect with TLS
	if useTLS {
		fmt.Println("Attempting to start TLS on the LDAP connection...")
		if err = l.StartTLS(&tls.Config{InsecureSkipVerify: true}); err != nil {
			fmt.Printf("Could not start TLS on the LDAP server connection: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✓ Start TLS Complete")
	}

	if err = l.Bind(bindUser, bindPass); err != nil {
		fmt.Printf("Failed to bind on the LDAP server connection: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✓ Bind Complete")

	fmt.Println("Starting message listening loop for a few seconds...")
	l.Start()
	time.Sleep(time.Second * 2)
	fmt.Println("No issues on message loop, closing connection")

	l.Close()

	fmt.Println("Finished.")
}
