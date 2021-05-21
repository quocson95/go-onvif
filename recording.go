package onvif

import "github.com/golang/glog"

func (device Device) GetRecordingSummary(nodeToken string) (interface{}, error) {
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
