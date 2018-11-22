package onvif

import (
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
		Body:  "<tds:GetDeviceInformation/>",
		XMLNs: deviceXMLNs,
		User: device.User,
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
func (device Device) GetNetworkInterfaces() (NetworkInterfaces, error) {
	// Create SOAP
	soap := SOAP{
		Body:  "<tds:GetNetworkInterfaces/>",
		XMLNs: deviceXMLNs,
		User: device.User,
		Password: device.Password,
	}

	// Send SOAP request
	_, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return NetworkInterfaces{}, err
	}

	return NetworkInterfaces{}, nil
}

// GetCapabilities fetch info of ONVIF camera's capabilities
func (device Device) GetCapabilities() (DeviceCapabilities, error) {
	// Create SOAP
	soap := SOAP{
		XMLNs: deviceXMLNs,
		Body: `<tds:GetCapabilities>
			<tds:Category>All</tds:Category>
		</tds:GetCapabilities>`,
		User: device.User,
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
	if err != nil {
		return DeviceCapabilities{}, err
	}

	eventsCap := make(map[string]bool)
	if mapEventsCap, ok := ifaceEventsCap.(map[string]interface{}); ok {
		for key, value := range mapEventsCap {
			if strings.ToLower(key) == "xaddr" {
				continue
			}

			key = strings.Replace(key, "WS", "", 1)
			eventsCap[key] = interfaceToBool(value)
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

	// Create final result
	deviceCapabilities := DeviceCapabilities{
		Network:   netCap,
		Media:	   mediaCap,
		Events:    eventsCap,
		Streaming: streamingCap,
	}

	return deviceCapabilities, nil
}

// GetDiscoveryMode fetch network discovery mode of an ONVIF camera
func (device Device) GetDiscoveryMode() (string, error) {
	// Create SOAP
	soap := SOAP{
		Body:  "<tds:GetDiscoveryMode/>",
		XMLNs: deviceXMLNs,
		User: device.User,
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
		Body:  "<tds:GetScopes/>",
		XMLNs: deviceXMLNs,
		User:device.User,
		Password:device.Password,
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
		Body:  "<tds:GetHostname/>",
		XMLNs: deviceXMLNs,
		User: device.User,
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
