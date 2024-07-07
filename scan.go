package main

import (
	"fmt"
	"log"
	"net"
)

func scanPort(host string, p int64) {
	// assuming localhost for now
	conn, err := net.Dial("tcp", fmt.Sprintf(":%d", p))
	if err != nil {
		log.Println("error connecting to server: ", err)
		return
	}

	log.Printf("connection to %s:%d succeeded\n", host, port)
	conn.Close()
}

func scanPorts(host, ports string) {
	var portStart, portEnd int64
	n, err := fmt.Sscanf(ports, "%d-%d", &portStart, &portEnd)
	if err != nil || n != 2 {
		log.Println("Usage: -z <hostname> <port-start>-<port-end>", err)
		return
	}

	for i := portStart; i <= portEnd; i++ {
		scanPort(host, i)
	}
}
