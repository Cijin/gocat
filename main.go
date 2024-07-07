package main

import (
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
)

var (
	isUdp  bool
	listen bool
	port   int
	z      string
	e      bool
	x      bool
)

func init() {
	flag.BoolVar(&isUdp, "u", false, "Use UDP")
	flag.BoolVar(&listen, "l", true, "Listen")
	flag.BoolVar(&x, "x", false, "Hex dump")
	flag.BoolVar(&e, "e", false, "Turn a process into a server")
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

	listenTcp(port, e, x)
}
