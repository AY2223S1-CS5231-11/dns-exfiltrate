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

	udpAddr := net.UDPAddr{
		Port: 53,
	}
	udpServer, err := net.ListenUDP("udp", &udpAddr)
	if err != nil {
		log.Fatalln(err)
	}
	defer udpServer.Close()

	buf := make([]byte, MAX_UDP_PACKET_SIZE)
	for {
		_, clientAddr, err := udpServer.ReadFromUDP(buf)
		if err != nil {
			log.Fatalln(err)
		}

		var request dns.Msg
		request.Unpack(buf)
		name := request.Question[0].Name
		fmt.Println(name)

		var reply dns.Msg
		reply.SetReply(&request)
		rr, err := dns.NewRR(fmt.Sprintf("%s A 8.8.8.8", name))
		if err != nil {
			log.Fatalln(err)
		}
		reply.Answer = append(reply.Answer, rr)

		response, err := reply.Pack()
		if err != nil {
			log.Fatalln(err)
		}
		udpServer.WriteToUDP(response, clientAddr)
	}
}
