package main

import (
	"log"
)

func main() {
	var localIP, remoteIP string
	var localPort, remotePort string

	err := parseIP(localhost, &localIP, &localPort)
	if err != nil {
		log.Fatal(err)
	}
	err = parseIP(remotehost, &remoteIP, &remotePort)
	if err != nil {
		log.Fatal(err)
	}

	err = TCPProxy(localIP, localPort, remoteIP, remotePort)
	if err != nil {
		log.Fatal(err)
	}
}
