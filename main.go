package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

var (
	isUdp  bool
	listen bool
	port   int
	z      string
)

func init() {
	flag.BoolVar(&isUdp, "u", true, "Use UDP")
	flag.BoolVar(&listen, "l", true, "Listen")
	flag.IntVar(&port, "p", 8080, "Port to listen on")
	flag.StringVar(&z, "z", "", "Connect to port without sending data")
}

func main() {
	flag.Parse()

	if z != "" {
		if len(os.Args) != 4 {
			log.Println("Usage: -z <hostname> <port>")
			return
		}

		if strings.Contains(os.Args[3], "-") {
			scanPorts(os.Args[2], os.Args[3])
			return
		}

		p, err := strconv.ParseInt(os.Args[3], 10, 64)
		if err != nil {
			log.Fatal("invalid port:", err)
			return
		}
		scanPort(os.Args[2], p)
		return
	}

	if isUdp {
		listenUdp(port)
		return
	}

	listenTcp(port)
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

func listenTcp(port int) {
	fmt.Println("listening on port:", port)

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Println("listen err:", err)
	}

	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal("accept err:", err)
		}

		go func(c net.Conn) {
			io.Copy(os.Stdout, c)

			c.Close()
		}(conn)

	}
}

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

		// write data back
		_, err = packetConn.WriteTo([]byte(fmt.Sprintf("recieved %d bytes\n", n)), addr)
		if err != nil {
			log.Println("error sending data to client:", err)
		}
	}
}
