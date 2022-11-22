package parser

import (
	"fmt"
	"os"

	"github.com/akamensky/argparse"
)

type Arguments struct {
	NameServer string
}

func ParseArgs() *Arguments {
	parser := argparse.NewParser("DNS Data Exfiltration Server", "Server for exfiltrating data over DNS")

	nameServer := parser.String("n", "nameserver", &argparse.Options{Required: true, Help: "Address of the nameserver that is being exfiltrated to"})

	err := parser.Parse(os.Args)
	if err != nil {
		// If there is a parse error, print usage.
		fmt.Print(parser.Usage(err))
		os.Exit(0)
	}

	return &Arguments{
		NameServer: *nameServer,
	}
}
