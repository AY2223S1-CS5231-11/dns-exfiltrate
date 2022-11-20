package parser

import (
	"fmt"
	"os"

	"github.com/akamensky/argparse"
)

type Arguments struct {
	NameServer string
	Filename   string
}

func ParseArgs() *Arguments {
	parser := argparse.NewParser("DNS Data Exfiltration Client", "Client for exfiltrating data over DNS")

	nameServer := parser.String("n", "nameserver", &argparse.Options{Required: true, Help: "Address of the nameserver to exfiltrate to"})
	filename := parser.String("f", "filename", &argparse.Options{Required: true, Help: "Name of file to exfiltrate"})

	err := parser.Parse(os.Args)
	if err != nil {
		// If there is a parse error, print usage.
		fmt.Print(parser.Usage(err))
	}

	return &Arguments{
		NameServer: *nameServer,
		Filename:   *filename,
	}
}
