package main

import (
	"dns-exfiltration-server/exfiltrator"
	"dns-exfiltration-server/parser"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
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

	// We can drop privileges after binding to port 53, which is protected.
	dropPrivileges()

	dnsExfiltrator := exfiltrator.NewDnsExfiltrator(args.NameServer)
	dnsExfiltrator.HandleDnsRequests(udpServer, args.NameServer)
}

func dropPrivileges() {
	if os.Geteuid() != 0 {
		return
	}

	// GID needs to be set before we lose privileges when setting UID.
	originalGidEnvVar := os.Getenv("SUDO_GID")
	originalGid, err := strconv.Atoi(originalGidEnvVar)
	if err != nil {
		log.Fatalln(err)
	}
	syscall.Setgid(originalGid)

	originalUidEnvVar := os.Getenv("SUDO_UID")
	originalUid, err := strconv.Atoi(originalUidEnvVar)
	if err != nil {
		log.Fatalln(err)
	}
	syscall.Setuid(originalUid)
}
