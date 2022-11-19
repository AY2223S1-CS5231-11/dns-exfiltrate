package main

import (
	"dns-exfiltration-client/parser"
	"fmt"
)

func main() {
	args := parser.ParseArgs()
	fmt.Println(args.Filename)
}
