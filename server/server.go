package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/miekg/dns"
)

const (
	// Messages carried by UDP are restricted to 512 bytes (not counting the IP
	// or UDP headers).
	// - https://www.rfc-editor.org/rfc/rfc1035#section-4.2.1
	MAX_UDP_PACKET_SIZE = 512
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

	buf := make([]byte, MAX_UDP_PACKET_SIZE)
	for {
		_, _, err := udpServer.ReadFrom(buf)
		if err != nil {
			log.Fatalln(err)
		}

		var request dns.Msg
		request.Unpack(buf)
		fmt.Println(request.Question[0].Name)
	}
}
