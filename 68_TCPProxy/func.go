package main

import (
	"errors"
	"flag"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

var localhost, remotehost string

func init() {
	flag.StringVar(&localhost, "l", "", "local host addr")
	flag.StringVar(&remotehost, "r", "", "remote host addr")
	flag.Parse()

	if localhost == "" || remotehost == "" {
		flag.Usage()
		os.Exit(-1)
	}
}

func parseIP(hostPort string, IP *string, port *string) error {
	if !strings.Contains(hostPort, ":") {
		return errors.New("not in format (172.17.0.1:80)")
	}
	// split hostPort into host and port
	idx := strings.LastIndex(hostPort, ":")
	host := hostPort[0:idx]
	*port = hostPort[idx + 1:]

	IPArray, err := net.LookupHost(host)
	if err != nil {
		log.Fatal(err)
	}
	// host has more than one ip, take only one, other will be implement later
	*IP = IPArray[0]

	return nil
}

func TCPProxy(localIP, localPort, remoteIP, remotePort string) error {
	remoteConn, err := net.Dial("tcp", remoteIP + ":" + remotePort)
	if err != nil {
		return errors.New("remote host or port closed")
	}

	localListener, err := net.Listen("tcp", localIP + ":" + localPort)
	if err != nil {
		log.Fatal(err)
		// return errors.New("could not listen on this ports")
	}

	for {
		localConn, err := localListener.Accept()
		if err != nil {
			return errors.New("localListener.Accept()")
		}

		// copy data to remote
		go func() {
			if _, err := io.Copy(remoteConn, localConn); err != nil {
				log.Fatal(err)
			}
		}()

		// receive data to remote
		if _, err := io.Copy(localConn, remoteConn); err != nil {
			log.Fatal(err)
		}
	}
}