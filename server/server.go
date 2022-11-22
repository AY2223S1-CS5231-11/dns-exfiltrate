package main

import (
	"dns-exfiltration-server/exfiltrator"
	"dns-exfiltration-server/parser"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	args := parser.ParseArgs()

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

	dnsExfiltrator := exfiltrator.NewDnsExfiltrator(args.NameServer)
	dnsExfiltrator.HandleDnsRequests(udpServer, args.NameServer)
}
