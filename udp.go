package main

import (
	"fmt"
	"log"
	"net"
)

func listenUdp(port int) {
	fmt.Println("listening on port:", port)

	packetConn, err := net.ListenPacket("udp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal("udp listen error:", err)
	}

	defer packetConn.Close()
	buf := make([]byte, 1024)

	for {
		n, addr, err := packetConn.ReadFrom(buf)
		if err != nil {
			log.Println("udp read error:", err)
		}

		log.Printf("recieved data %d bytes from client at addr: %s\n Data: %s", n, addr.String(), buf)

		_, err = packetConn.WriteTo([]byte(fmt.Sprintf("recieved %d bytes\n", n)), addr)
		if err != nil {
			log.Println("error sending data to client:", err)
		}
	}
}
