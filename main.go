package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var (
	isUdp  bool
	listen bool
	port   int
	z      string
	e      bool
)

func init() {
	flag.BoolVar(&isUdp, "u", false, "Use UDP")
	flag.BoolVar(&listen, "l", true, "Listen")
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

	listenTcp(port, e)
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

func listenTcp(port int, e bool) {
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

		if e {
			go execCommand(conn)
		} else {
			go func(c net.Conn) {
				io.Copy(os.Stdout, c)

				c.Close()
			}(conn)
		}
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

		_, err = packetConn.WriteTo([]byte(fmt.Sprintf("recieved %d bytes\n", n)), addr)
		if err != nil {
			log.Println("error sending data to client:", err)
		}
	}
}

func execCommand(conn net.Conn) {
	fmt.Println("Executing shell...")
	cmd := exec.Command("/bin/sh")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Println("Error creating stdin pipe for cmd", err)
		return
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println("Error creating stdout pipe for cmd", err)
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Println("Error creating stderr pipe for cmd", err)
		return
	}

	defer func() {
		err := stdin.Close()
		if err != nil {
			log.Println("Error closing stdin pipe:", err)
		}

		err = stdout.Close()
		if err != nil {
			log.Println("Error closing stdout pipe:", err)
		}

		err = stderr.Close()
		if err != nil {
			log.Println("Error closing stdout pipe:", err)
		}
	}()

	if err := cmd.Start(); err != nil {
		log.Fatal("Error starting command:", err)
		return
	}

	go func() {
		// what could go wrong :D
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			fmt.Fprintln(stdin, scanner.Text())
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			fmt.Fprintln(conn, scanner.Text())
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			fmt.Fprintln(conn, scanner.Text())
		}
	}()

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
}
