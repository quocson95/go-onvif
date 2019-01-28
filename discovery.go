package onvif

import (
	"errors"
	"net"
	"regexp"
	"strings"
	"time"
	"github.com/golang/glog"
	"github.com/clbanning/mxj"
	"github.com/satori/go.uuid"
)

var errWrongDiscoveryResponse = errors.New("Response is not related to discovery request ")

// StartDiscovery send a WS-Discovery message and wait for all matching device to respond
func StartDiscoveryOn(interfaceName string, duration time.Duration) ([]Device, error) {
	itf, err := net.InterfaceByName(interfaceName) //here your interface

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

// StartDiscovery send a WS-Discovery message and wait for all matching device to respond
func StartDiscovery(interfaceName string, duration time.Duration) ([]Device, error) {
	// Get list of interface address
	if interfaceName != "" {
		return StartDiscoveryOn(interfaceName, duration)
	}

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return []Device{}, err
	}

	// Fetch IPv4 address
	var ipAddrs []string
	for _, addr := range addrs {
		ipAddr, ok := addr.(*net.IPNet)
		if ok && !ipAddr.IP.IsLoopback() && ipAddr.IP.To4() != nil {
			ipAddrs = append(ipAddrs, ipAddr.IP.String())
		}
	}

	// Create initial discovery results
	var discoveryResults []Device

	// Discover device on each interface's network
	for _, ipAddr := range ipAddrs {
		devices, err := discoverDevices(ipAddr, duration)
		if err != nil {
			return []Device{}, err
		}

		discoveryResults = append(discoveryResults, devices...)
	}

	return discoveryResults, nil
}

func discoverDevices(ipAddr string, duration time.Duration) ([]Device, error) {
	// Create WS-Discovery request
	u, _ := uuid.NewV4()
	requestID := "uuid:" + u.String()
	//request := `
	//	<?xml version="1.0" encoding="UTF-8"?>
	//	<e:Envelope
	//	    xmlns:e="http://www.w3.org/2003/05/soap-envelope"
	//	    xmlns:w="http://schemas.xmlsoap.org/ws/2004/08/addressing"
	//	    xmlns:d="http://schemas.xmlsoap.org/ws/2005/04/discovery"
	//	    xmlns:dn="http://www.onvif.org/ver10/network/wsdl">
	//	    <e:Header>
	//	        <w:MessageID>` + requestID + `</w:MessageID>
	//	        <w:To e:mustUnderstand="true">urn:schemas-xmlsoap-org:ws:2005:04:discovery</w:To>
	//	        <w:Action a:mustUnderstand="true">http://schemas.xmlsoap.org/ws/2005/04/discovery/Probe
	//	        </w:Action>
	//	    </e:Header>
	//	    <e:Body>
	//	        <d:Probe>
	//	            <d:Types>dn:NetworkVideoTransmitter</d:Types>
	//	        </d:Probe>
	//	    </e:Body>
	//	</e:Envelope>`

	request := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope
					xmlns:s="http://www.w3.org/2003/05/soap-envelope"
					xmlns:a="http://schemas.xmlsoap.org/ws/2004/08/addressing">
					<s:Header>
						<a:Action s:mustUnderstand="1">http://schemas.xmlsoap.org/ws/2005/04/discovery/Probe</a:Action>
						<a:MessageID>` + requestID + `</a:MessageID>
						<a:ReplyTo><a:Address>http://schemas.xmlsoap.org/ws/2004/08/addressing/role/anonymous</a:Address></a:ReplyTo>
						<a:To s:mustUnderstand="1">urn:schemas-xmlsoap-org:ws:2005:04:discovery</a:To>
					</s:Header>
					<s:Body>
						<Probe xmlns="http://schemas.xmlsoap.org/ws/2005/04/discovery">
							<d:Types xmlns:d="http://schemas.xmlsoap.org/ws/2005/04/discovery" xmlns:dp0="http://www.onvif.org/ver10/network/wsdl">dp0:NetworkVideoTransmitter</d:Types>
						</Probe>
					</s:Body>
				</s:Envelope>`

	// Clean WS-Discovery message
	request = regexp.MustCompile(`>\s+<`).ReplaceAllString(request, "><")
	request = regexp.MustCompile(`\s+`).ReplaceAllString(request, " ")

	// Create UDP address for local and multicast address
	localAddress, err := net.ResolveUDPAddr("udp4", ipAddr + ":0")
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
	glog.Infof("Discover response: %s", string(buffer))

	// Inital result
	result := Device{}

	buffer = []byte(`<?xml version="1.0" encoding="utf-8"?>
<SOAP-ENV:Envelope
   xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope"
   xmlns:SOAP-ENC="http://www.w3.org/2003/05/soap-encoding"
   xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
   xmlns:xsd="http://www.w3.org/2001/XMLSchema"
   xmlns:wsa="http://schemas.xmlsoap.org/ws/2004/08/addressing"
   xmlns:wsd="http://schemas.xmlsoap.org/ws/2005/04/discovery"
   xmlns:dn="http://www.onvif.org/ver10/network/wsdl">
   <SOAP-ENV:Header>
       <wsa:MessageID SOAP-ENV:mustUnderstand="true">urn:uuid:D5F76F58-99EF-444A-9995-cc43127dcb5c</wsa:MessageID>
       <wsa:RelatesTo SOAP-ENV:mustUnderstandard="true">uuid:c396c249-a721-493c-b34c-29cb9028a31a</wsa:RelatesTo>
       <wsa:To SOAP-ENV:mustUnderstand="true">http://schemas.xmlsoap.org/ws/2004/08/addressing/role/anonymous</wsa:To>
       <wsa:Action SOAP-ENV:mustUnderstand="true">http://schemas.xmlsoap.org/ws/2005/04/discovery/ProbeMatches</wsa:Action>
       <wsd:AppSequence MessageNumber="9" InstanceId="1543855428"></wsd:AppSequence>
   </SOAP-ENV:Header>
   <SOAP-ENV:Body>
       <wsd:ProbeMatches>
           <wsd:ProbeMatch>
               <wsa:EndpointReference>
                   <wsa:Address>urn:uuid:5A946080-C4AE-4A5F-9238-7682c8a9f37e</wsa:Address>
               </wsa:EndpointReference>
               <wsd:Types>dn:NetworkVideoTransmitter</wsd:Types>
               <wsd:Scopes>onvif://www.onvif.org/type/video_encoder onvif://www.onvif.org/type/audio_encoder onvif://www.onvif.org/Profile/Streaming onvif://www.onvif.org/hardware/LNB5100 onvif://www.onvif.org/location/GMT onvif://www.onvif.org/name/LNB5100-C04DEE</wsd:Scopes>
               <wsd:XAddrs>http://192.168.88.237/onvif/device_service</wsd:XAddrs>
               <wsd:MetadataVersion>1543855428</wsd:MetadataVersion>
           </wsd:ProbeMatch>
       </wsd:ProbeMatches>
   </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`)

	// Parse XML to map
	mapXML, err := mxj.NewMapXml(buffer)
	if err != nil {
		glog.Warningf("Parse response error %v", err)
		return result, err
	}

	// Check if this response is for our request
	responseMessageID, err := mapXML.ValueForPath("Envelope.Header.RelatesTo")
	if err != nil {
		glog.Warningf("Parse message id error %v", err)
		return result, err
	}

	if responseMessageMap, ok := responseMessageID.(map[string] interface{}); ok{
		responseMessage := responseMessageMap["#text"].(string)
		if responseMessage != messageID {
			glog.Info(responseMessage)
			glog.Info(messageID)
			return result, errWrongDiscoveryResponse
		}
	}

	// Get device's ID and clean it
	deviceID, _ := mapXML.ValueForPathString("Envelope.Body.ProbeMatches.ProbeMatch.EndpointReference.Address")
	deviceID = strings.Replace(deviceID, "urn:uuid:", "", 1)
	glog.Infof("Discover device id: %s", deviceID)

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
	glog.Infof("Discover address: %s", xAddrs)
	if len(listXAddr) == 0 {
		glog.Warning("Discover address len 0")
		return result, errors.New("Device does not have any xAddr ")
	}

	// Finalize result
	result.ID = deviceID
	result.Name = deviceName
	result.XAddr = listXAddr[0]

	return result, nil
}
