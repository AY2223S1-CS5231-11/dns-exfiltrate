package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	fmt.Println("Starting server...")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("Stopping server...")
		os.Exit(0)
	}()

	udpServer, err := net.ListenPacket("udp", ":53")
	if err != nil {
		log.Fatalln(err)
	}
	defer udpServer.Close()

	buf := make([]byte, 1024)
	for {
		_, _, err := udpServer.ReadFrom(buf)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println(buf)
	}
}
