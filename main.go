package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net"
	"strconv"

	"gopkg.in/ldap.v2"
)

var hostname = ""
var port = 389
var useTLS = false
var bindUser = ""
var bindPass = ""

func main() {
	flag.StringVar(&hostname, "h", "", "LDAP Hostname")
	flag.IntVar(&port, "p", 389, "LDAP Port")
	flag.BoolVar(&useTLS, "tls", useTLS, "LDAP use tls")
	flag.StringVar(&bindUser, "bindUser", bindUser, "LDAP Bind Username")
	flag.StringVar(&bindPass, "bindPass", bindPass, "LDAP Bind Password")
	flag.Parse()

	addr := net.JoinHostPort(hostname, strconv.Itoa(port))
	l, err := ldap.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer l.Close()

	fmt.Println("✓ Connected to LDAP Server")

	// Reconnect with TLS
	if useTLS {
		err = l.StartTLS(&tls.Config{InsecureSkipVerify: true})
		if err != nil {
			panic(err)
		}
		fmt.Println("✓ Start TLS Complete")
	}

	if err = l.Bind(bindUser, bindPass); err != nil {
		panic(err)
	}
	fmt.Println("✓ Bind Complete")

	fmt.Println("Finished.")
}
