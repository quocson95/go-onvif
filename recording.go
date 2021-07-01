package onvif

import "github.com/golang/glog"

func (device Device) GetRecordingSummary() (interface{}, error) {
	// create soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body:     `<GetRecordingSummary xmlns="http://www.onvif.org/ver10/search/wsdl"/>`,
	}

	var result interface{}
	// send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return result, err
	}

	// parse response
	data, err := response.ValueForPath("Envelope.Body.GetRecordingSummaryResponse.Summary")
	if err != nil {
		return result, err
	}
	glog.Infof("Data %v", data)
	return result, nil
}

func (device Device) GetMediaAttributes(time string) (interface{}, error) {
	// create soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body:     `<GetMediaAttributes xmlns="http://www.onvif.org/ver10/search/wsdl">
						<Time>` + time + `</Time>					
					</GetMediaAttributes>`,
	}

	var result interface{}
	// send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return result, err
	}

	// parse response
	data, err := response.ValueForPath("Envelope.Body.GetMediaAttributesResponse.MediaAttributes")
	if err != nil {
		return result, err
	}
	glog.Infof("Data %v", data)
	return result, nil
}

func (device Device) GetReplayUri(recordingToken string) (interface{}, error) {
	// create soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body:     `<GetReplayUri xmlns="http://www.onvif.org/ver10/replay/wsdl">
						<RecordingToken>` + recordingToken + `</RecordingToken>		
						<StreamSetup>
							<Stream xmlns="http://www.onvif.org/ver10/schema">RTP-Unicast</Stream>
							<Transport xmlns="http://www.onvif.org/ver10/schema">
								<Protocol xmlns="http://www.onvif.org/ver10/schema">TCP</Protocol>
							</Transport>
						</StreamSetup>
					</GetReplayUri>`,
	}

	var result interface{}
	// send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return result, err
	}

	// parse response
	data, err := response.ValueForPath("Envelope.Body.GetReplayUriResponse")
	if err != nil {
		return result, err
	}
	glog.Infof("Data %v", data)
	return result, nil
}

func (device Device) GetReplayConfiguration() (interface{}, error) {
	// create soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body:     `<GetReplayConfiguration xmlns="http://www.onvif.org/ver10/replay/wsdl"/>`,
	}

	var result interface{}
	// send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return result, err
	}

	// parse response
	data, err := response.ValueForPath("Envelope.Body.GetReplayConfigurationResponse")
	if err != nil {
		return result, err
	}
	glog.Infof("Data %v", data)
	return result, nil
}

func (device Device) GetReplayServiceCapabilities() (interface{}, error) {
	// create soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body:     `<GetServiceCapabilities xmlns="http://www.onvif.org/ver10/replay/wsdl"/>`,
	}

	var result interface{}
	// send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return result, err
	}

	// parse response
	data, err := response.ValueForPath("Envelope.Body.GetServiceCapabilitiesResponse")
	if err != nil {
		return result, err
	}
	glog.Infof("Data %v", data)
	return result, nil
}
