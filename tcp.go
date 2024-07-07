package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

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
