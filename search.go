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
		Body: `<GetMediaAttributes xmlns="http://www.onvif.org/ver10/search/wsdl">
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

func (device Device) FindRecordings() (interface{}, error) {
	// create soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body: `<FindRecordings xmlns="http://www.onvif.org/ver10/search/wsdl">
					</FindRecordings>`,
	}

	var result interface{}
	// send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return result, err
	}

	// parse response
	data, err := response.ValueForPath("Envelope.Body.FindRecordingsResponse")
	if err != nil {
		return result, err
	}
	glog.Infof("Data %v", data)
	return result, nil
}
