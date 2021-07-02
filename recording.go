package onvif

import "github.com/golang/glog"

func (device Device) GetRecordingConfiguration(recordingToken string) (interface{}, error) {
	// create soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body: `<GetRecordingConfiguration xmlns="http://www.onvif.org/ver10/recording/wsdl">
						<RecordingToken>` + recordingToken + `</RecordingToken>
					</GetRecordingConfiguration>`,
	}

	var result interface{}
	// send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return result, err
	}

	// parse response
	data, err := response.ValueForPath("Envelope.Body.GetRecordingConfigurationResponse")
	if err != nil {
		return result, err
	}
	glog.Infof("Data %v", data)
	return result, nil
}
