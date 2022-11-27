package parser

import (
	"fmt"
	"os"

	"github.com/akamensky/argparse"
)

type Arguments struct {
	NameServer string
	Filename   string
	Delay      int
}

func ParseArgs() *Arguments {
	parser := argparse.NewParser("DNS Data Exfiltration Client", "Client for exfiltrating data over DNS")

	nameServer := parser.String("n", "nameserver", &argparse.Options{Required: true, Help: "Address of the nameserver to exfiltrate to"})
	filename := parser.String("f", "filename", &argparse.Options{Required: true, Help: "Name of file to exfiltrate"})
	delay := parser.Int("d", "delay", &argparse.Options{Default: 0, Required: false, Help: "Delay in milliseconds between each DNS request"})

	err := parser.Parse(os.Args)
	if err != nil {
		// If there is a parse error, print usage.
		fmt.Print(parser.Usage(err))
		os.Exit(0)
	}

	return &Arguments{
		NameServer: *nameServer,
		Filename:   *filename,
		Delay:      *delay,
	}
}
