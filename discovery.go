package onvif

import (
	"errors"
	"kakacam-hub/config"
	"net"
	"regexp"
	"strings"
	"time"
	//"fmt"
	"github.com/clbanning/mxj"
	"github.com/satori/go.uuid"
)

var errWrongDiscoveryResponse = errors.New("Response is not related to discovery request ")

// StartDiscovery send a WS-Discovery message and wait for all matching device to respond
func StartDiscovery(duration time.Duration) ([]Device, error) {
	itf, err := net.InterfaceByName(config.GetAppConfig().IntFace) //here your interface

	if err != nil {
		return []Device{}, err
	}

	item, _ := itf.Addrs()
	var ip net.IP
	for _, addr := range item {
		switch v := addr.(type) {
		case *net.IPNet:
			if !v.IP.IsLoopback() {
				if v.IP.To4() != nil { //Verify if IP is IPV4
					ip = v.IP
				}
			}
		}
	}

	if ip == nil {
		return []Device{}, err
	}

	// Discover device on interface's network
	devices, err := discoverDevices(ip.String(), duration)

	return devices, nil
}

func discoverDevices(ipAddr string, duration time.Duration) ([]Device, error) {
	// Create WS-Discovery request
	u, _ := uuid.NewV4()
	requestID := "uuid:" + u.String()
	request := `		
		<?xml version="1.0" encoding="UTF-8"?>
		<e:Envelope
		    xmlns:e="http://www.w3.org/2003/05/soap-envelope"
		    xmlns:w="http://schemas.xmlsoap.org/ws/2004/08/addressing"
		    xmlns:d="http://schemas.xmlsoap.org/ws/2005/04/discovery"
		    xmlns:dn="http://www.onvif.org/ver10/network/wsdl">
		    <e:Header>
		        <w:MessageID>` + requestID + `</w:MessageID>
		        <w:To e:mustUnderstand="true">urn:schemas-xmlsoap-org:ws:2005:04:discovery</w:To>
		        <w:Action a:mustUnderstand="true">http://schemas.xmlsoap.org/ws/2005/04/discovery/Probe
		        </w:Action>
		    </e:Header>
		    <e:Body>
		        <d:Probe>
		            <d:Types>dn:NetworkVideoTransmitter</d:Types>
		        </d:Probe>
		    </e:Body>
		</e:Envelope>`

	// Clean WS-Discovery message
	request = regexp.MustCompile(`>\s+<`).ReplaceAllString(request, "><")
	request = regexp.MustCompile(`\s+`).ReplaceAllString(request, " ")

	// Create UDP address for local and multicast address
	localAddress, err := net.ResolveUDPAddr("udp4", ipAddr+":0")
	if err != nil {
		return []Device{}, err
	}

	multicastAddress, err := net.ResolveUDPAddr("udp4", "239.255.255.250:3702")
	if err != nil {
		return []Device{}, err
	}

	// Create UDP connection to listen for respond from matching device
	conn, err := net.ListenUDP("udp", localAddress)
	if err != nil {
		return []Device{}, err
	}
	defer conn.Close()

	// Set connection's timeout
	err = conn.SetDeadline(time.Now().Add(duration))
	if err != nil {
		return []Device{}, err
	}

	// Send WS-Discovery request to multicast address
	_, err = conn.WriteToUDP([]byte(request), multicastAddress)
	if err != nil {
		return []Device{}, err
	}

	// Create initial discovery results
	var discoveryResults []Device

	// Keep reading UDP message until timeout
	for {
		// Create buffer and receive UDP response
		buffer := make([]byte, 10*1024)
		_, _, err = conn.ReadFromUDP(buffer)

		// Check if connection timeout
		if err != nil {
			if udpErr, ok := err.(net.Error); ok && udpErr.Timeout() {
				break
			} else {
				return discoveryResults, err
			}
		}

		//fmt.Println(string(buffer))
		// Read and parse WS-Discovery response
		device, err := readDiscoveryResponse(requestID, buffer)
		if err != nil && err != errWrongDiscoveryResponse {
			return discoveryResults, err
		}

		// Push device to results
		discoveryResults = append(discoveryResults, device)
	}

	return discoveryResults, nil
}

// readDiscoveryResponse reads and parses WS-Discovery response
func readDiscoveryResponse(messageID string, buffer []byte) (Device, error) {
	// Inital result
	result := Device{}

	// Parse XML to map
	mapXML, err := mxj.NewMapXml(buffer)
	if err != nil {
		return result, err
	}

	// Check if this response is for our request
	responseMessageID, _ := mapXML.ValueForPathString("Envelope.Header.RelatesTo")
	if responseMessageID != messageID {
		return result, errWrongDiscoveryResponse
	}

	// Get device's ID and clean it
	deviceID, _ := mapXML.ValueForPathString("Envelope.Body.ProbeMatches.ProbeMatch.EndpointReference.Address")
	deviceID = strings.Replace(deviceID, "urn:uuid:", "", 1)

	// Get device's name
	deviceName := ""
	scopes, _ := mapXML.ValueForPathString("Envelope.Body.ProbeMatches.ProbeMatch.Scopes")
	for _, scope := range strings.Split(scopes, " ") {
		if strings.HasPrefix(scope, "onvif://www.onvif.org/name/") {
			deviceName = strings.Replace(scope, "onvif://www.onvif.org/name/", "", 1)
			deviceName = strings.Replace(deviceName, "_", " ", -1)
			break
		}
	}

	// Get device's xAddrs
	xAddrs, _ := mapXML.ValueForPathString("Envelope.Body.ProbeMatches.ProbeMatch.XAddrs")
	listXAddr := strings.Split(xAddrs, " ")
	if len(listXAddr) == 0 {
		return result, errors.New("Device does not have any xAddr ")
	}

	// Finalize result
	result.ID = deviceID
	result.Name = deviceName
	result.XAddr = listXAddr[0]

	return result, nil
}
