package main

import (
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func listenTcp(port int, e, x bool) {
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
				if x {
					buf := make([]byte, 4096)

					n, err := conn.Read(buf)
					if err != nil {
						log.Println("Error reading from connection:", err)
						return
					}

					log.Println(hex.Dump(buf[:n]))
				} else {
					io.Copy(os.Stdout, c)
				}

				c.Close()
			}(conn)
		}
	}
}
