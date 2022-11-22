package exfiltrator

import (
	"dns-exfiltration-server/fileutils"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/miekg/dns"
)

const (
	// Messages carried by UDP are restricted to 512 bytes (not counting the IP
	// or UDP headers).
	// - https://www.rfc-editor.org/rfc/rfc1035#section-4.2.1
	MAX_UDP_PACKET_SIZE = 512
	EXFILTRATION_DIR    = "./exfiltrated-data"
)

type dnsExfiltrator struct {
	NameServer string
}

func NewDnsExfiltrator(nameServer string) *dnsExfiltrator {
	absoluteNameServer := nameServer
	if nameServer[len(nameServer)-1] != '.' {
		absoluteNameServer += "."
	}

	fileutils.CreateDirIfNotExists(EXFILTRATION_DIR)

	return &dnsExfiltrator{
		NameServer: absoluteNameServer,
	}
}

func (ex *dnsExfiltrator) HandleDnsRequests(udpServer *net.UDPConn, nameServer string) {
	buf := make([]byte, MAX_UDP_PACKET_SIZE)
	for {
		_, clientAddr, err := udpServer.ReadFromUDP(buf)
		if err != nil {
			log.Fatalln(err)
		}

		var request dns.Msg
		request.Unpack(buf)
		name := request.Question[0].Name
		subdomains := strings.Split(name, nameServer)[0]
		data := strings.ReplaceAll(subdomains, ".", "")
		fmt.Println(data)

		var reply dns.Msg
		reply.SetReply(&request)
		rr, err := dns.NewRR(fmt.Sprintf("%s 300 IN A 8.8.8.8", name))
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
