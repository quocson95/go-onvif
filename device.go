package onvif

import (
	"github.com/golang/glog"
	"strings"
)

var deviceXMLNs = []string{
	`xmlns:tds="http://www.onvif.org/ver10/device/wsdl"`,
	`xmlns:tt="http://www.onvif.org/ver10/schema"`,
}

// GetInformation fetch information of ONVIF camera
func (device Device) GetInformation() (DeviceInformation, error) {
	// Create SOAP
	soap := SOAP{
		Body:     "<tds:GetDeviceInformation/>",
		XMLNs:    deviceXMLNs,
		User:     device.User,
		Password: device.Password,
	}

	// Send SOAP request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return DeviceInformation{}, err
	}

	// Parse response to interface
	deviceInfo, err := response.ValueForPath("Envelope.Body.GetDeviceInformationResponse")
	if err != nil {
		return DeviceInformation{}, err
	}

	// Parse interface to struct
	result := DeviceInformation{}
	if mapInfo, ok := deviceInfo.(map[string]interface{}); ok {
		result.Manufacturer = interfaceToString(mapInfo["Manufacturer"])
		result.Model = interfaceToString(mapInfo["Model"])
		result.FirmwareVersion = interfaceToString(mapInfo["FirmwareVersion"])
		result.SerialNumber = interfaceToString(mapInfo["SerialNumber"])
		result.HardwareID = interfaceToString(mapInfo["HardwareId"])
	}

	return result, nil
}

// GetInformation fetch information of ONVIF camera
func (device Device) GetNetworkInterfaces() ([]NetworkInterface, error) {
	// Create SOAP
	soap := SOAP{
		Body:     "<tds:GetNetworkInterfaces/>",
		XMLNs:    deviceXMLNs,
		User:     device.User,
		Password: device.Password,
	}

	// Send SOAP request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return nil, err
	}

	// Parse response to interface
	networkInterfacesInfos, err := response.ValuesForPath("Envelope.Body.GetNetworkInterfacesResponse.NetworkInterfaces")
	if err != nil {
		return nil, err
	}

	// Parse interface to struct
	result := make([]NetworkInterface, 0)
	for _, networkInterfacesInfo := range networkInterfacesInfos {
		if mapNetworkInterfacesInfo, ok := networkInterfacesInfo.(map[string]interface{}); ok {
			glog.Infof("mapNetworkInterfacesInfo %v", mapNetworkInterfacesInfo)
			networkInterface := NetworkInterface{}
			networkInterface.Token = interfaceToString(mapNetworkInterfacesInfo["-token"])
			networkInterface.Enabled = interfaceToBool(mapNetworkInterfacesInfo["Enabled"])
			if mapIPv4Info, ok := mapNetworkInterfacesInfo["IPv4"].(map[string]interface{}); ok {
				networkInterface.IPv4.Enabled = interfaceToBool(mapIPv4Info["Enabled"])
				if mapIPv4Config, ok := mapIPv4Info["Config"].(map[string]interface{}); ok {
					networkInterface.IPv4.Config.DHCP = interfaceToBool(mapIPv4Config["DHCP"])
					if mapIPv4FromDHCPConfig, ok := mapIPv4Config["FromDHCP"].(map[string]interface{}); ok {
						networkInterface.IPv4.Config.FromDHCP.Address = interfaceToString(mapIPv4FromDHCPConfig["Address"])
						networkInterface.IPv4.Config.FromDHCP.PrefixLength = interfaceToInt(mapIPv4FromDHCPConfig["PrefixLength"])
					}
					if mapIPv4ManualConfig, ok := mapIPv4Config["Manual"].(map[string]interface{}); ok {
						networkInterface.IPv4.Config.Manual.Address = interfaceToString(mapIPv4ManualConfig["Address"])
						networkInterface.IPv4.Config.Manual.PrefixLength = interfaceToInt(mapIPv4ManualConfig["PrefixLength"])
					}
					if mapIPv4LinkLocalConfig, ok := mapIPv4Config["LinkLocal"].(map[string]interface{}); ok {
						networkInterface.IPv4.Config.LinkLocal.Address = interfaceToString(mapIPv4LinkLocalConfig["Address"])
						networkInterface.IPv4.Config.LinkLocal.PrefixLength = interfaceToInt(mapIPv4LinkLocalConfig["PrefixLength"])
					}
				}
			}
			if mapInfo, ok := mapNetworkInterfacesInfo["Info"].(map[string]interface{}); ok {
				networkInterface.Info.Name = interfaceToString(mapInfo["Name"])
				networkInterface.Info.MTU = interfaceToInt(mapInfo["MTU"])
				networkInterface.Info.HwAddress = interfaceToString(mapInfo["HwAddress"])
			}
			result = append(result, networkInterface)
		}
	}

	return result, nil
}

func (device Device) SetNetworkInterfaces(networkInterface NetworkInterface) error {
	//create soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body: `<SetNetworkInterfaces xmlns="http://www.onvif.org/ver10/device/wsdl">
					<InterfaceToken>` + networkInterface.Token + `</InterfaceToken>
					<NetworkInterface>
						<Enabled xmlns="http://www.onvif.org/ver10/schema">` + boolToString(networkInterface.Enabled) + `</Enabled>
						<MTU xmlns="http://www.onvif.org/ver10/schema">` + intToString(networkInterface.Info.MTU) + `</MTU>
						<IPv4 xmlns="http://www.onvif.org/ver10/schema">
							<Enabled>true</Enabled>
							<Manual>
								<Address>` + networkInterface.IPv4.Config.Manual.Address + `</Address>
								<PrefixLength>` + intToString(networkInterface.IPv4.Config.Manual.PrefixLength) + `</PrefixLength>
							</Manual>
							<DHCP>` + boolToString(networkInterface.IPv4.Config.DHCP) + `</DHCP>
						</IPv4>
					</NetworkInterface>
 			  </SetNetworkInterfaces>`,
	}
	// send request
	response, err := soap.SendRequest(device.XAddr)

	if err != nil {
		return err
	}

	_, err = response.ValueForPath("Envelope.Body.SetNetworkInterfacesResponse")
	if err != nil {
		return err
	}
	return nil
}

// GetCapabilities fetch info of ONVIF camera's capabilities
func (device Device) GetCapabilities() (DeviceCapabilities, error) {
	// Create SOAP
	soap := SOAP{
		XMLNs: deviceXMLNs,
		Body: `<tds:GetCapabilities>
			<tds:Category>All</tds:Category>
		</tds:GetCapabilities>`,
		User:     device.User,
		Password: device.Password,
	}

	// Send SOAP request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return DeviceCapabilities{}, err
	}

	// Get network capabilities
	envelopeBodyPath := "Envelope.Body.GetCapabilitiesResponse.Capabilities"
	ifaceNetCap, err := response.ValueForPath(envelopeBodyPath + ".Device.Network")
	if err != nil {
		return DeviceCapabilities{}, err
	}

	netCap := NetworkCapabilities{}
	if mapNetCap, ok := ifaceNetCap.(map[string]interface{}); ok {
		netCap.DynDNS = interfaceToBool(mapNetCap["DynDNS"])
		netCap.IPFilter = interfaceToBool(mapNetCap["IPFilter"])
		netCap.IPVersion6 = interfaceToBool(mapNetCap["IPVersion6"])
		netCap.ZeroConfig = interfaceToBool(mapNetCap["ZeroConfiguration"])
	}

	// Get events capabilities
	ifaceEventsCap, err := response.ValueForPath(envelopeBodyPath + ".Events")
	eventsCap := EventsCapabilities{}
	if err == nil {
		eventsCap.Events = map[string]bool{}
		if mapEventsCap, ok := ifaceEventsCap.(map[string]interface{}); ok {
			for key, value := range mapEventsCap {
				if strings.ToLower(key) == "xaddr" {
					eventsCap.XAddr = interfaceToString(value)
					continue
				}

				key = strings.Replace(key, "WS", "", 1)
				eventsCap.Events[key] = interfaceToBool(value)
			}
		}
	}

	// Get streaming capabilities
	ifaceStreamingCap, err := response.ValueForPath(envelopeBodyPath + ".Media.StreamingCapabilities")
	if err != nil {
		return DeviceCapabilities{}, err
	}

	streamingCap := make(map[string]bool)
	if mapStreamingCap, ok := ifaceStreamingCap.(map[string]interface{}); ok {
		for key, value := range mapStreamingCap {
			key = strings.Replace(key, "_", " ", -1)
			streamingCap[key] = interfaceToBool(value)
		}
	}

	// Get media capabilities
	ifaceMediaCap, err := response.ValueForPath(envelopeBodyPath + ".Media")
	if err != nil {
		return DeviceCapabilities{}, err
	}

	mediaCap := MediaCapabilities{}
	if mapMediaCap, ok := ifaceMediaCap.(map[string]interface{}); ok {
		mediaCap.XAddr = interfaceToString(mapMediaCap["XAddr"])
	}

	// Get Ptz capabilities
	ifacePtzCap, err := response.ValueForPath(envelopeBodyPath + ".PTZ")
	ptzCap := PTZCapabilities{}
	if err == nil {
		if mapPtzCap, ok := ifacePtzCap.(map[string]interface{}); ok {
			ptzCap.XAddr = interfaceToString(mapPtzCap["XAddr"])
		}
	}

	// Create final result
	deviceCapabilities := DeviceCapabilities{
		Network:   netCap,
		Media:     mediaCap,
		Ptz:       ptzCap,
		EventsCap: eventsCap,
		Streaming: streamingCap,
	}

	return deviceCapabilities, nil
}

// GetDiscoveryMode fetch network discovery mode of an ONVIF camera
func (device Device) GetDiscoveryMode() (string, error) {
	// Create SOAP
	soap := SOAP{
		Body:     "<tds:GetDiscoveryMode/>",
		XMLNs:    deviceXMLNs,
		User:     device.User,
		Password: device.Password,
	}

	// Send SOAP request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return "", err
	}

	// Parse response
	discoveryMode, _ := response.ValueForPathString("Envelope.Body.GetDiscoveryModeResponse.DiscoveryMode")
	return discoveryMode, nil
}

// GetScopes fetch scopes of an ONVIF camera
func (device Device) GetScopes() ([]string, error) {
	// Create SOAP
	soap := SOAP{
		Body:     "<tds:GetScopes/>",
		XMLNs:    deviceXMLNs,
		User:     device.User,
		Password: device.Password,
	}

	// Send SOAP request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return nil, err
	}

	// Parse response to interface
	ifaceScopes, err := response.ValuesForPath("Envelope.Body.GetScopesResponse.Scopes")
	if err != nil {
		return nil, err
	}

	// Convert interface to array of scope
	scopes := []string{}
	for _, ifaceScope := range ifaceScopes {
		if mapScope, ok := ifaceScope.(map[string]interface{}); ok {
			scope := interfaceToString(mapScope["ScopeItem"])
			scopes = append(scopes, scope)
		}
	}

	return scopes, nil
}

// GetHostname fetch hostname of an ONVIF camera
func (device Device) GetHostname() (HostnameInformation, error) {
	// Create SOAP
	soap := SOAP{
		Body:     "<tds:GetHostname/>",
		XMLNs:    deviceXMLNs,
		User:     device.User,
		Password: device.Password,
	}

	// Send SOAP request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return HostnameInformation{}, err
	}

	// Parse response to interface
	ifaceHostInfo, err := response.ValueForPath("Envelope.Body.GetHostnameResponse.HostnameInformation")
	if err != nil {
		return HostnameInformation{}, err
	}

	// Parse interface to struct
	hostnameInfo := HostnameInformation{}
	if mapHostInfo, ok := ifaceHostInfo.(map[string]interface{}); ok {
		hostnameInfo.Name = interfaceToString(mapHostInfo["Name"])
		hostnameInfo.FromDHCP = interfaceToBool(mapHostInfo["FromDHCP"])
	}

	return hostnameInfo, nil
}

func (device Device) GetSystemDateAndTime() (SystemDateAndTime, error) {
	// Create SOAP
	soap := SOAP{
		XMLNs: deviceXMLNs,
		Body:  `<GetSystemDateAndTime xmlns="http://www.onvif.org/ver10/device/wsdl"/>`,
	}

	systemDT := SystemDateAndTime{}

	// send SOAP request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return systemDT, err
	}

	// Parse response to interface
	ifaceSystemDateTime, err := response.ValueForPath("Envelope.Body.GetSystemDateAndTimeResponse.SystemDateAndTime")
	if err != nil {
		return systemDT, err
	}

	// Parse response to struct
	if mapSystemDateTime, ok := ifaceSystemDateTime.(map[string]interface{}); ok {
		systemDT.DateTimeType = interfaceToString(mapSystemDateTime["DateTimeType"])
		systemDT.DaylightSavings = interfaceToBool(mapSystemDateTime["DaylightSavings"])

		if ifaceTZ, ok := mapSystemDateTime["TimeZone"].(map[string]interface{}); ok {
			systemDT.TZ = interfaceToString(ifaceTZ["TZ"])
		}

		if ifaceUCTDT, ok := mapSystemDateTime["UTCDateTime"].(map[string]interface{}); ok {
			time := ifaceUCTDT["Time"]
			if mapTime, ok := time.(map[string]interface{}); ok {
				systemDT.Hour = interfaceToInt(mapTime["Hour"])
				systemDT.Minute = interfaceToInt(mapTime["Minute"])
				systemDT.Second = interfaceToInt(mapTime["Second"])
			}

			date := ifaceUCTDT["Date"]

			if mapDate, ok := date.(map[string]interface{}); ok {
				systemDT.Day = interfaceToInt(mapDate["Day"])
				systemDT.Month = interfaceToInt(mapDate["Month"])
				systemDT.Year = interfaceToInt(mapDate["Year"])
			}
		}
	}

	return systemDT, nil
}

func (device Device) SetSystemDateAndTime(systemDT SystemDateAndTime) error {
	// create Body request
	var body string
	if systemDT.DateTimeType == "Manual" { // Manual mode
		body = `<SetSystemDateAndTime xmlns="http://www.onvif.org/ver10/device/wsdl">
					<DateTimeType>` + systemDT.DateTimeType + `</DateTimeType>
					<DaylightSavings>` + boolToString(systemDT.DaylightSavings) + `</DaylightSavings>
					<TimeZone>
						<TZ xmlns="http://www.onvif.org/ver10/schema">` + systemDT.TZ + `</TZ>
					</TimeZone>
					<UTCDateTime>
						<Time xmlns="http://www.onvif.org/ver10/schema">
							<Hour>` + intToString(systemDT.Hour) + `</Hour>
							<Minute>` + intToString(systemDT.Minute) + `</Minute>
							<Second>` + intToString(systemDT.Second) + `</Second>
						</Time>
						<Date xmlns="http://www.onvif.org/ver10/schema">
							<Year>` + intToString(systemDT.Year) + `</Year>
							<Month>` + intToString(systemDT.Month) + `</Month>
							<Day>` + intToString(systemDT.Day) + `</Day>
						</Date>
					</UTCDateTime>
				</SetSystemDateAndTime>`
	} else { // NTP mode
		body = `<SetSystemDateAndTime xmlns="http://www.onvif.org/ver10/device/wsdl">
					<DateTimeType>` + systemDT.DateTimeType + `</DateTimeType>
					<DaylightSavings>` + boolToString(systemDT.DaylightSavings) + `</DaylightSavings>`

		if systemDT.TZ != "" {
			body += `<TimeZone><TZ xmlns="http://www.onvif.org/ver10/schema">` + systemDT.TZ + `</TZ></TimeZone>`
		}

		body += `</SetSystemDateAndTime>`
	}

	// Create SOAP
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body:     body,
	}

	// send soap request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return err
	}

	_, err = response.ValueForPath("Envelope.Body.SetSystemDateAndTimeResponse")
	if err != nil {
		return err
	}

	return nil
}

func (device Device) GetNTP() (NTPInformation, error) {
	//create SOAP
	soap := SOAP{
		XMLNs:    deviceXMLNs,
		User:     device.User,
		Password: device.Password,
		Body:     `<GetNTP xmlns="http://www.onvif.org/ver10/device/wsdl"/>`,
	}

	ntpInformation := NTPInformation{}

	// send request
	response, err := soap.SendRequest(device.XAddr)

	if err != nil {
		return ntpInformation, err
	}

	// parse response into interface
	ifaceNTPInformation, err := response.ValueForPath("Envelope.Body.GetNTPResponse.NTPInformation")
	if err != nil {
		return ntpInformation, err
	}

	// parse interface to struct
	if mapNTPInformation, ok := ifaceNTPInformation.(map[string]interface{}); ok {
		ntpInformation.FromDHCP = interfaceToBool(mapNTPInformation["FromDHCP"])
		ntpInformation.NTPNetworkHost = make([]NetworkHost, 0)
		if ntpInformation.FromDHCP {
			if ifaceManuals, ok := mapNTPInformation["NTPFromDHCP"].([]interface{}); ok {
				for _, ifaceManual := range ifaceManuals {
					if mapNTPManual, ok := ifaceManual.(map[string]interface{}); ok {
						NTPNetworkHost := NetworkHost{}

						NTPNetworkHost.Type = interfaceToString(mapNTPManual["Type"])
						NTPNetworkHost.IPv4Address = interfaceToString(mapNTPManual["IPv4Address"])
						NTPNetworkHost.DNSname = interfaceToString(mapNTPManual["DNSname"])

						ntpInformation.NTPNetworkHost = append(ntpInformation.NTPNetworkHost, NTPNetworkHost)
					}
				}
			} else {
				if mapNTPManual, ok := mapNTPInformation["NTPFromDHCP"].(map[string]interface{}); ok {
					NTPNetworkHost := NetworkHost{}

					NTPNetworkHost.Type = interfaceToString(mapNTPManual["Type"])
					NTPNetworkHost.IPv4Address = interfaceToString(mapNTPManual["IPv4Address"])
					NTPNetworkHost.DNSname = interfaceToString(mapNTPManual["DNSname"])

					ntpInformation.NTPNetworkHost = append(ntpInformation.NTPNetworkHost, NTPNetworkHost)
				}
			}

		} else {
			if ifaceManuals, ok := mapNTPInformation["NTPManual"].([]interface{}); ok {
				for _, ifaceManual := range ifaceManuals {
					if mapNTPManual, ok := ifaceManual.(map[string]interface{}); ok {
						NTPNetworkHost := NetworkHost{}

						NTPNetworkHost.Type = interfaceToString(mapNTPManual["Type"])
						NTPNetworkHost.IPv4Address = interfaceToString(mapNTPManual["IPv4Address"])
						NTPNetworkHost.DNSname = interfaceToString(mapNTPManual["DNSname"])

						ntpInformation.NTPNetworkHost = append(ntpInformation.NTPNetworkHost, NTPNetworkHost)
					}
				}
			} else {
				if mapNTPManual, ok := mapNTPInformation["NTPManual"].(map[string]interface{}); ok {
					NTPNetworkHost := NetworkHost{}

					NTPNetworkHost.Type = interfaceToString(mapNTPManual["Type"])
					NTPNetworkHost.IPv4Address = interfaceToString(mapNTPManual["IPv4Address"])
					NTPNetworkHost.DNSname = interfaceToString(mapNTPManual["DNSname"])

					ntpInformation.NTPNetworkHost = append(ntpInformation.NTPNetworkHost, NTPNetworkHost)
				}
			}
		}
	}

	return ntpInformation, nil
}

func (device Device) SetNTP(ntpInformation NTPInformation) error {
	// create soap
	soap := SOAP{
		XMLNs:    deviceXMLNs,
		User:     device.User,
		Password: device.Password,
		Body: `<SetNTP xmlns="http://www.onvif.org/ver10/device/wsdl">
					<FromDHCP>` + boolToString(ntpInformation.FromDHCP) + `</FromDHCP>
					<NTPManual>
						<Type xmlns="http://www.onvif.org/ver10/schema">` + ntpInformation.NTPNetworkHost[0].Type + `</Type>
						<IPv4Address xmlns="http://www.onvif.org/ver10/schema">` + ntpInformation.NTPNetworkHost[0].IPv4Address + `</IPv4Address>
						<DNSname xmlns="http://www.onvif.org/ver10/schema">` + ntpInformation.NTPNetworkHost[0].DNSname + `</DNSname>
					</NTPManual>
				</SetNTP>`,
	}

	// send soap request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return err
	}

	_, err = response.ValueForPath("Envelope.Body.SetNTPResponse")
	if err != nil {
		return err
	}

	return nil
}

func (device Device) SystemReboot() (string, error) {
	// create SOAP
	soap := SOAP{
		XMLNs:    deviceXMLNs,
		User:     device.User,
		Password: device.Password,
		Body:     `<SystemReboot xmlns="http://www.onvif.org/ver10/device/wsdl"/>`,
	}

	var message string

	// send request
	response, err := soap.SendRequest(device.XAddr)

	if err != nil {
		return message, err
	}

	// parse response into interface
	ifaceMessage, err := response.ValueForPath("Envelope.Body.SystemRebootResponse")
	if err != nil {
		return message, err
	}

	// parse message reboot
	if mapMessage, ok := ifaceMessage.(map[string]interface{}); ok {
		message = interfaceToString(mapMessage["Message"])
	}

	return message, nil
}

func (device Device) GetDNS() (DNSInformation, error) {
	soap := SOAP{
		XMLNs:    deviceXMLNs,
		User:     device.User,
		Password: device.Password,
		Body:     `<GetDNS xmlns="http://www.onvif.org/ver10/device/wsdl"/>`,
	}

	dnsInformation := DNSInformation{}

	// send soap request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return dnsInformation, err
	}

	// parse response to interface
	ifaceDNSInformation, err := response.ValueForPath("Envelope.Body.GetDNSResponse.DNSInformation")

	if err != nil {
		return dnsInformation, err
	}

	// parse interface to map
	if mapDNSInformation, ok := ifaceDNSInformation.(map[string]interface{}); ok {
		dnsInformation.FromDHCP = interfaceToBool(mapDNSInformation["FromDHCP"])
		dnsInformation.SearchDomain = interfaceToString(mapDNSInformation["SearchDomain"])

		if dnsInformation.FromDHCP {
			if mapListDNSFromDHCP, ok := mapDNSInformation["DNSFromDHCP"].([]interface{}); ok {
				for _, mapDNSFromDHCPIf := range mapListDNSFromDHCP {
					DNSAddress := IPAddress{}
					if mapDNSFromDHCP, ok := mapDNSFromDHCPIf.(map[string]interface{}); ok {
						DNSAddress.Type = interfaceToString(mapDNSFromDHCP["Type"])
						DNSAddress.IPv4Address = interfaceToString(mapDNSFromDHCP["IPv4Address"])
					}
					dnsInformation.DNSAddress = append(dnsInformation.DNSAddress, DNSAddress)
				}
			}
		} else {
			if mapListDNSManual, ok := mapDNSInformation["DNSManual"].([]interface{}); ok {
				for _, mapDNSManualIf := range mapListDNSManual {
					DNSAddress := IPAddress{}
					if mapDNSManual, ok := mapDNSManualIf.(map[string]interface{}); ok {
						DNSAddress.Type = interfaceToString(mapDNSManual["Type"])
						DNSAddress.IPv4Address = interfaceToString(mapDNSManual["IPv4Address"])
					}
					dnsInformation.DNSAddress = append(dnsInformation.DNSAddress, DNSAddress)
				}
			}
		}
	}

	return dnsInformation, nil
}

func (device Device) SetDNS(dnsInformation DNSInformation) error {
	// create soap
	soap := SOAP{
		XMLNs:    deviceXMLNs,
		User:     device.User,
		Password: device.Password,
		Body: `<SetDNS xmlns="http://www.onvif.org/ver10/device/wsdl">
				<FromDHCP>` + boolToString(dnsInformation.FromDHCP) + `</FromDHCP>
				<SearchDomain>` + dnsInformation.SearchDomain + `</SearchDomain>
				<DNSManual>
					<Type xmlns="http://www.onvif.org/ver10/schema">` + dnsInformation.DNSAddress[0].Type + `</Type>
					<IPv4Address xmlns="http://www.onvif.org/ver10/schema">` + dnsInformation.DNSAddress[0].IPv4Address + `</IPv4Address>
				</DNSManual>
			  </SetDNS>`,
	}

	// send soap request
	response, err := soap.SendRequest(device.XAddr)

	if err != nil {
		return err
	}

	_, err = response.ValueForPath("Envelope.Body.SetDNSResponse")
	if err != nil {
		return err
	}

	return nil
}

func (device Device) GetDynamicDNS() (DynamicDNSInformation, error) {
	soap := SOAP{
		XMLNs:    deviceXMLNs,
		User:     device.User,
		Password: device.Password,
		Body:     `<GetDynamicDNS xmlns="http://www.onvif.org/ver10/device/wsdl"/>`,
	}
	result := DynamicDNSInformation{}

	// send resquest
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return result, err
	}

	// parse response
	ifaceDynamicDNS, err := response.ValueForPath("Envelope.Body.GetDynamicDNSResponse")
	if err != nil {
		return result, err
	}

	// parse interface
	if mapDynamicDNS, ok := ifaceDynamicDNS.(map[string]interface{}); ok {
		if mapDynamicDNSInformation, ok := mapDynamicDNS["DynamicDNSInformation"].(map[string]interface{}); ok {
			result.Type = interfaceToString(mapDynamicDNSInformation["Type"])
			result.Name = interfaceToString(mapDynamicDNSInformation["Name"])
			result.TTL = interfaceToString(mapDynamicDNSInformation["TTL"])
		}
	}

	return result, nil
}

func (device Device) SetHostName(nameToken string) error {
	// create soap
	soap := SOAP{
		XMLNs:    deviceXMLNs,
		User:     device.User,
		Password: device.Password,
		Body: `<SetHostname xmlns="http://www.onvif.org/ver10/device/wsdl">
				<Name xmlns="http://www.onvif.org/ver10/schema">` + nameToken + `</Name>
			   </SetHostname>`,
	}

	// send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return err
	}

	_, err = response.ValueForPath("Envelope.Body.SetHostnameResponse")
	if err != nil {
		return err
	}

	return nil
}

func (device Device) GetNetworkProtocols() ([]NetworkProtocol, error) {
	// create soap
	soap := SOAP{
		XMLNs:    deviceXMLNs,
		User:     device.User,
		Password: device.Password,
		Body:     `<GetNetworkProtocols xmlns:="http://www.onvif.org/ver10/device/wsdl"/>`,
	}

	result := []NetworkProtocol{}
	// send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return result, err
	}

	// parse response to interface
	ifaceNetworkProtocols, err := response.ValuesForPath("Envelope.Body.GetNetworkProtocolsResponse.NetworkProtocols")
	if err != nil {
		return result, err
	}

	// parse iface
	for _, ifaceNetworkProtocol := range ifaceNetworkProtocols {
		if mapNetworkProtocol, ok := ifaceNetworkProtocol.(map[string]interface{}); ok {
			networkProtocol := NetworkProtocol{}

			networkProtocol.Name = interfaceToString(mapNetworkProtocol["Name"])
			networkProtocol.Enabled = interfaceToBool(mapNetworkProtocol["Enabled"])
			networkProtocol.Port = interfaceToInt(mapNetworkProtocol["Port"])

			result = append(result, networkProtocol)
		}
	}

	return result, nil
}

func (device Device) SetNetworkProtocols(protocols []NetworkProtocol) error {
	// create body for array protocols
	var protocolsBody string = ``
	for _, protocol := range protocols {
		protocolsBody += `<NetworkProtocols>
							<Name xmlns="http://www.onvif.org/ver10/schema">` + protocol.Name + `</Name>
							<Enabled xmlns="http://www.onvif.org/ver10/schema">` + boolToString(protocol.Enabled) + `</Enable>
							<Port xmlns="http://www.onvif.org/ver10/schema">` + intToString(protocol.Port) + `</Port>
						  </NetworkProtocols>`
	}

	// create soap
	soap := SOAP{
		XMLNs:    deviceXMLNs,
		User:     device.User,
		Password: device.Password,
		Body:     `<SetNetworkProtocols xmlns="http://www.onvif.org/ver10/device/wsdl">` + protocolsBody + `</SetNetworkProtocols>`,
	}
	// send soap request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return err
	}

	_, err = response.ValueForPath("Envelope.Body.SetNetworkProtocolsResponse")
	if err != nil {
		return err
	}

	return nil
}

func (device Device) SetScopes(listScopes []string) error {
	// create scopes body
	var scopesBody string
	for _, scope := range listScopes {
		scopesBody += `<Scopes>` + scope + `</Scopes>`
	}

	// create soap
	soap := SOAP{
		XMLNs:    deviceXMLNs,
		User:     device.User,
		Password: device.Password,
		Body:     `<SetScopes xmlns="http://www.onvif.org/ver10/device/wsdl">` + scopesBody + `</SetScopes>`,
	}

	// send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return err
	}

	_, err = response.ValueForPath("Envelope.Body.SetScopesResponse")
	if err != nil {
		return err
	}

	return nil
}

func (device Device) AddScopes(listScopes []string) error {
	// create scopes body
	var scopesBody string
	for _, scope := range listScopes {
		scopesBody += `<ScopeItem>` + scope + `</ScopeItem>`
	}

	// create soap
	soap := SOAP{
		XMLNs:    deviceXMLNs,
		User:     device.User,
		Password: device.Password,
		Body:     `<AddScopes xmlns="http://www.onvif.org/ver10/device/wsdl">` + scopesBody + `</AddScopes>`,
	}

	// send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return err
	}

	_, err = response.ValueForPath("Envelope.Body.AddScopesResponse")
	if err != nil {
		return err
	}

	return nil
}

func (device Device) RemoveScopes(listScopes []string) ([]string, error) {
	// create scopes body
	var scopesBody string
	for _, scope := range listScopes {
		scopesBody += `<ScopeItem>` + scope + `</ScopeItem>`
	}
	// create soap
	soap := SOAP{
		XMLNs:    deviceXMLNs,
		User:     device.User,
		Password: device.Password,
		Body:     `<RemoveScopes xmlns="http://www.onvif.org/ver10/device/wsdl">` + scopesBody + `</RemoveScopes>`,
	}

	var result []string
	// send request
	respone, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return result, err
	}

	// parse response
	ifaceScopeItems, err := respone.ValuesForPath("Envelope.Body.RemoveScopesResponse.ScopeItem")
	for _, scopeItem := range ifaceScopeItems {
		result = append(result, interfaceToString(scopeItem))
	}

	return result, nil
}

func (device Device) GetNetworkDefaultGateway() (NetworkGateway, error) {
	//create soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body:     `<GetNetworkDefaultGateway xmlns="http://www.onvif.org/ver10/device/wsdl"/>`,
	}

	result := NetworkGateway{}
	// send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return result, err
	}

	// parse response
	ifaceNetworkGateway, err := response.ValueForPath("Envelope.Body.GetNetworkDefaultGatewayResponse.NetworkGateway")
	if err != nil {
		return result, err
	}

	if mapNetworkGateway, ok := ifaceNetworkGateway.(map[string]interface{}); ok {
		result.IPv4Address = interfaceToString(mapNetworkGateway["IPv4Address"])
	}

	return result, nil
}

func (device Device) SetNetworkDefaultGateway(defaultGateway NetworkGateway) error {
	//create soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body: `<SetNetworkDefaultGateway xmlns="http://www.onvif.org/ver10/device/wsdl">
					<IPv4Address>` + defaultGateway.IPv4Address + `</IPv4Address>
 			  </SetNetworkDefaultGateway>`,
	}
	// send request
	response, err := soap.SendRequest(device.XAddr)

	if err != nil {
		return err
	}

	_, err = response.ValueForPath("Envelope.Body.SetNetworkDefaultGatewayResponse")
	if err != nil {
		return err
	}
	return nil
}

func (device Device) GetUsers() ([]User, error) {
	//create soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body:     `<GetUsers xmlns="http://www.onvif.org/ver10/device/wsdl"/>`,
	}

	result := []User{}

	// send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return result, err
	}

	// parse response into interface
	ifaceUsers, err := response.ValuesForPath("Envelope.Body.GetUsersResponse.User")
	if err != nil {
		return result, err
	}

	for _, ifaceUser := range ifaceUsers {
		if mapUser, ok := ifaceUser.(map[string]interface{}); ok {
			user := User{}

			user.Username = interfaceToString(mapUser["Username"])
			user.Password = interfaceToString(mapUser["Password"])
			user.UserLevel = interfaceToString(mapUser["UserLevel"])

			result = append(result, user)
		}
	}

	return result, nil
}

func (device Device) SetUser(user User) error {
	// create soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body: `<SetUser xmlns="http://www.onvif.org/ver10/device/wsdl"><User>
					<Username xmlns="http://www.onvif.org/ver10/schema">` + user.Username + `</Username>
					<Password xmlns="http://www.onvif.org/ver10/schema">` + user.Password + `</Password>
					<UserLevel xmlns="http://www.onvif.org/ver10/schema">` + user.UserLevel + `</UserLevel>
				</User></SetUser>`,
	}
	// send soap request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return err
	}

	_, err = response.ValueForPath("Envelope.Body.SetUserResponse")
	if err != nil {
		return err
	}

	return nil
}

func (device Device) DeleteUsers(usernames []string) error {
	// create usernamebody
	var usernameBody = ``
	for _, username := range usernames {
		usernameBody += `<Username>` + username + `</Username>`
	}

	// create soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body:     `<DeleteUsers xmlns="http://www.onvif.org/ver10/device/wsdl">` + usernameBody + `</DeleteUsers>`,
	}

	// send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return err
	}

	_, err = response.ValueForPath("Envelope.Body.DeleteUsersResponse")
	if err != nil {
		return err
	}

	return nil
}

func (device Device) CreateUsers(users []User) error {
	// create UserBody
	var userBody = ``
	for _, user := range users {
		userBody += `<User>
						<Username xmlns="http://www.onvif.org/ver10/schema">` + user.Username + `</Username>
						<Password xmlns="http://www.onvif.org/ver10/schema">` + user.Password + `</Password>
						<UserLevel xmlns="http://www.onvif.org/ver10/schema">` + user.UserLevel + `</UserLevel>
					 </User>`
	}

	//create soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body:     `<CreateUsers xmlns="http://www.onvif.org/ver10/device/wsdl">` + userBody + `</CreateUsers>`,
	}

	// send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return err
	}

	_, err = response.ValueForPath("Envelope.Body.CreateUsersResponse")
	if err != nil {
		return err
	}

	return nil
}

func (device Device) GetRelayOutputs() (RelayOutput, error) {
	// create soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body:     `<GetRelayOutputs xmlns="http://www.onvif.org/ver10/device/wsdl"/>`,
	}

	result := RelayOutput{}

	// send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return result, err
	}

	// parse response
	ifaceRelayOutput, err := response.ValueForPath("Envelope.Body.GetRelayOutputsResponse.RelayOutputs")
	if err != nil {
		return result, err
	}

	// parse into result
	if mapRelayOutput, ok := ifaceRelayOutput.(map[string]interface{}); ok {
		result.Token = interfaceToString(mapRelayOutput["-token"])
		// parse properties
		if mapProperties, ok := mapRelayOutput["Properties"].(map[string]interface{}); ok {
			result.Properties.Mode = interfaceToString(mapProperties["Mode"])
			result.Properties.DelayTime = interfaceToString(mapProperties["DelayTime"])
			result.Properties.IdleState = interfaceToString(mapProperties["IdleState"])
		}
	}

	return result, nil
}

func (device Device) GetZeroConfiguration() (NetworkZeroConfiguration, error) {
	// create soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body:     `<GetZeroConfiguration xmlns="http://www.onvif.org/ver10/device/wsdl"/>`,
	}

	result := NetworkZeroConfiguration{}

	// send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return result, err
	}

	//parse response
	ifaceNetworkZeroConfiguration, err := response.ValueForPath("Envelope.Body.GetZeroConfigurationResponse")
	if err != nil {
		return result, err
	}

	// parse interface
	if mapNetworkZeroConfiguration, ok := ifaceNetworkZeroConfiguration.(map[string]interface{}); ok {
		if mapZeroConfiguration, ok := mapNetworkZeroConfiguration["ZeroConfiguration"].(map[string]interface{}); ok {
			result.InterfaceToken = interfaceToString(mapZeroConfiguration["InterfaceToken"])
			result.Enabled = interfaceToBool(mapZeroConfiguration["Enabled"])
			// parse addresses
			for _, address := range mapZeroConfiguration["Addresses"].([]interface{}) {
				result.Addresses = append(result.Addresses, interfaceToString(address))
			}
		}
	}

	glog.Info(result)
	return result, nil
}

func (device Device) GetServices() ([]Service, error) {
	// create soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body:     `<GetServices xmlns="http://www.onvif.org/ver10/device/wsdl"><IncludeCapability>false</IncludeCapability></GetServices>`,
	}

	result := []Service{}

	//send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return result, err
	}

	// parse response into interface
	ifaceServices, err := response.ValuesForPath("Envelope.Body.GetServicesResponse.Service")
	if err != nil {
		return result, err
	}

	// parse interface into struct
	for _, ifaceService := range ifaceServices {
		if mapService, ok := ifaceService.(map[string]interface{}); ok {
			service := Service{}

			service.Namespace = interfaceToString(mapService["Namespace"])
			service.XAddr = interfaceToString(mapService["XAddr"])

			// parse version
			if mapVersion, ok := mapService["Version"].(map[string]interface{}); ok {
				service.Version.Major = interfaceToInt(mapVersion["Major"])
				service.Version.Minor = interfaceToInt(mapVersion["Minor"])
			}
			result = append(result, service)
		}
	}

	return result, nil
}

func (device Device) GetServiceCapabilities() ([]Service, error) {
	// create soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body:     `<GetServices xmlns="http://www.onvif.org/ver10/device/wsdl"><IncludeCapability>true</IncludeCapability></GetServices>`,
	}

	result := []Service{}

	//send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return result, err
	}

	// parse response into interface
	ifaceServices, err := response.ValuesForPath("Envelope.Body.GetServicesResponse.Service")
	if err != nil {
		return result, err
	}

	// parse interface into struct
	for _, ifaceService := range ifaceServices {
		if mapService, ok := ifaceService.(map[string]interface{}); ok {
			service := Service{}

			service.Namespace = interfaceToString(mapService["Namespace"])
			service.XAddr = interfaceToString(mapService["XAddr"])

			// parse version
			if mapVersion, ok := mapService["Version"].(map[string]interface{}); ok {
				service.Version.Major = interfaceToInt(mapVersion["Major"])
				service.Version.Minor = interfaceToInt(mapVersion["Minor"])
			}

			// parse capabilities
			if mapCapabilities, ok := mapService["Capabilities"].(map[string]interface{}); ok {
				if mapCapabilitiesAtrr, ok := mapCapabilities["Capabilities"].(map[string]interface{}); ok {
					if strings.Index(service.Namespace, "device") > -1 { // Device Capabilities
						deviceCapabilities := DeviceCapabilitiesService{}
						// parse network
						network := DeviceNetworkCapabilitiesService{}
						if mapNetwork, ok := mapCapabilitiesAtrr["Network"].(map[string]interface{}); ok {
							network.DynDNS = interfaceToBool(mapNetwork["-DynDNS"])
							network.IPFilter = interfaceToBool(mapNetwork["-IPFilter"])
							network.IPVersion6 = interfaceToBool(mapNetwork["-IPVersion6"])
							network.NTP = interfaceToInt(mapNetwork["-NTP"])
							network.ZeroConfiguration = interfaceToBool(mapNetwork["-ZeroConfiguration"])
						}
						deviceCapabilities.Network = network

						// parse System
						system := DeviceSystemCapabilitiesService{}
						if mapSystem, ok := mapCapabilitiesAtrr["System"].(map[string]interface{}); ok {
							system.DiscoveryBye = interfaceToBool(mapSystem["-DiscoveryBye"])
							system.DiscoveryResolve = interfaceToBool(mapSystem["-DiscoveryResolve"])
							system.FirmwareUpgrade = interfaceToBool(mapSystem["-FirmwareUpgrade"])
							system.RemoteDiscovery = interfaceToBool(mapSystem["-RemoteDiscovery"])
							system.SystemBackup = interfaceToBool(mapSystem["-SystemBackup"])
							system.SystemLogging = interfaceToBool(mapSystem["-SystemLogging"])
						}
						deviceCapabilities.System = system

						// parse Security
						security := DeviceSecurityCapabilitiesService{}
						if mapSecurity, ok := mapCapabilitiesAtrr["Security"].(map[string]interface{}); ok {
							security.AccesssPolicyConfig = interfaceToBool(mapSecurity["-AccessPolicyConfig"])
							security.DefaultAccessPolicy = interfaceToBool(mapSecurity["-DefaultAccessPolicy"])
							security.Dot1X = interfaceToBool(mapSecurity["-Dot1X"])
							security.HttpDigest = interfaceToBool(mapSecurity["-HttpDigest"])
							security.KerberosToken = interfaceToBool(mapSecurity["-KerberosToken"])
							security.OnboardKeyGeneration = interfaceToBool(mapSecurity["-OnboardKeyGeneration"])
							security.RELToken = interfaceToBool(mapSecurity["-RELToken"])
							security.RemoteUserHandling = interfaceToBool(mapSecurity["-RemoteUserHandling"])
							security.SAMLToken = interfaceToBool(mapSecurity["-SAMLToken"])
							security.TLS10 = interfaceToBool(mapSecurity["-TLS1.0"])
							security.TLS11 = interfaceToBool(mapSecurity["-TLS1.1"])
							security.TLS12 = interfaceToBool(mapSecurity["-TLS1.2"])
							security.UsernameToken = interfaceToBool(mapSecurity["-UsernameToken"])
							security.X509Token = interfaceToBool(mapSecurity["-X.509Token"])
						}
						deviceCapabilities.Security = security

						// add to service
						capabilities := CapabilitiesService{}
						capabilities.Name = "device"
						capabilities.Capabilities = deviceCapabilities
						service.Capabilities = capabilities
					} else if strings.Index(service.Namespace, "media") > -1 { // Media Capabilities
						mediaCapabilities := MediaCapabilitiesService{}

						mediaCapabilities.OSD = interfaceToBool(mapCapabilitiesAtrr["-OSD"])
						mediaCapabilities.Rotation = interfaceToBool(mapCapabilitiesAtrr["-Rotation"])
						mediaCapabilities.SnapshotUri = interfaceToBool(mapCapabilitiesAtrr["-SnapshotUri"])
						mediaCapabilities.VideoSourceMode = interfaceToBool(mapCapabilitiesAtrr["-VideoSourceMode"])

						// parse ProfileCapabilities
						if mapProfileCapabilities, ok := mapCapabilitiesAtrr["ProfileCapabilities"].(map[string]interface{}); ok {
							mediaCapabilities.ProfileCapabilities.MaximumNumberOfProfiles = interfaceToInt(mapProfileCapabilities["MaximumNumberOfProfiles"])
						}

						// parse Streaming Capabilities
						if mapStreamingCapabilities, ok := mapCapabilitiesAtrr["StreamingCapabilities"].(map[string]interface{}); ok {
							streamingCapabilities := MediaStreamingCapabilitiesService{}

							streamingCapabilities.NonAggregateControl = interfaceToBool(mapStreamingCapabilities["-NonAggregateControl"])
							streamingCapabilities.NoRTSPStreaming = interfaceToBool(mapStreamingCapabilities["-NoRTSPStreaming"])
							streamingCapabilities.RTP_RTSP_TCP = interfaceToBool(mapStreamingCapabilities["-RTP_RTSP_TCP"])
							streamingCapabilities.RTP_TCP = interfaceToBool(mapStreamingCapabilities["-RTP_TCP"])
							streamingCapabilities.RTPMulticast = interfaceToBool(mapStreamingCapabilities["-RTPMulticast"])

							mediaCapabilities.StreamingCapabilities = streamingCapabilities
						}

						// add to service
						capabilities := CapabilitiesService{}
						capabilities.Name = "media"
						capabilities.Capabilities = mediaCapabilities
						service.Capabilities = capabilities
					} else if strings.Index(service.Namespace, "events") > -1 { // Events Capabilities
						eventsCapabilities := EventsCapabilitiesService{}

						eventsCapabilities.MaxNotificationProducers = interfaceToInt(mapCapabilitiesAtrr["-MaxNotificationProducers"])
						eventsCapabilities.MaxPullPoints = interfaceToInt(mapCapabilitiesAtrr["-MaxPullPoints"])
						eventsCapabilities.PersistentNotificationStorage = interfaceToBool(mapCapabilitiesAtrr["-PersistentNotificationStorage"])
						eventsCapabilities.WSPausableSubscriptionManagerInterfaceSupport = interfaceToBool(mapCapabilitiesAtrr["-WSPausableSubscriptionManagerInterfaceSupport"])
						eventsCapabilities.WSPullPointSupport = interfaceToBool(mapCapabilitiesAtrr["-WSPullPointSupport"])
						eventsCapabilities.WSSubscriptionPolicySupport = interfaceToBool(mapCapabilitiesAtrr["-WSSubscriptionPolicySupport"])

						// add to service
						capabilities := CapabilitiesService{}
						capabilities.Name = "events"
						capabilities.Capabilities = eventsCapabilities
						service.Capabilities = capabilities
					} else if strings.Index(service.Namespace, "imaging") > -1 { // Imaging Capabilities
						imagingCapabilities := ImagingCapabilitiesService{}

						imagingCapabilities.ImageStabilization = interfaceToBool(mapCapabilitiesAtrr["-ImageStabilization"])

						// add to service
						capabilities := CapabilitiesService{}
						capabilities.Name = "imaging"
						capabilities.Capabilities = imagingCapabilities
						service.Capabilities = capabilities
					} else if strings.Index(service.Namespace, "ptz") > -1 { // PTZ Capabilities
						ptzCapabilities := PTZCapabilitiesService{}

						ptzCapabilities.EFlip = interfaceToBool(mapCapabilitiesAtrr["-EFlip"])
						ptzCapabilities.GetCompatibleConfigurations = interfaceToBool(mapCapabilitiesAtrr["-GetCompatibleConfigurations"])
						ptzCapabilities.Reverse = interfaceToBool(mapCapabilitiesAtrr["-Reverse"])

						// add to service
						capabilities := CapabilitiesService{}
						capabilities.Name = "PTZ"
						capabilities.Capabilities = ptzCapabilities
						service.Capabilities = capabilities
					}
				}
			}
			result = append(result, service)
		}
	}

	return result, nil
}
