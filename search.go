package onvif

import "github.com/golang/glog"

func (device Device) GetRecordingSummary() ([]RecordingSummary, error) {
	// create soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body:     `<GetRecordingSummary xmlns="http://www.onvif.org/ver10/search/wsdl"/>`,
	}

	result := make([]RecordingSummary, 0)
	// send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return result, err
	}

	// parse response
	summaries, err := response.ValuesForPath("Envelope.Body.GetRecordingSummaryResponse.Summary")
	if err != nil {
		return result, err
	}

	for _, summaryIf := range summaries {
		if summary, ok := summaryIf.(map[string]interface{}); ok {
			recordingSummary := RecordingSummary{}
			recordingSummary.DataFrom = interfaceToString(summary["DataFrom"])
			recordingSummary.DataUntil = interfaceToString(summary["DataUntil"])
			recordingSummary.NumberRecordings = interfaceToInt(summary["NumberRecordings"])
			result = append(result, recordingSummary)
		}
	}

	return result, nil
}

func (device Device) GetMediaAttributes(time string) ([]MediaAttributes, error) {
	// create soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body: `<GetMediaAttributes xmlns="http://www.onvif.org/ver10/search/wsdl">
					<Time>` + time + `</Time>					
			   </GetMediaAttributes>`,
	}

	result := make([]MediaAttributes, 0)
	// send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return result, err
	}

	// parse response
	mediaAttributesIfs, err := response.ValuesForPath("Envelope.Body.GetMediaAttributesResponse.MediaAttributes")
	if err != nil {
		return result, err
	}
	for _, mediaAttributesIf := range mediaAttributesIfs {
		if mediaAttributes, ok := mediaAttributesIf.(map[string]interface{}); ok {
			ma := MediaAttributes{}
			ma.RecordingToken = interfaceToString(mediaAttributes["RecordingToken"])
			ma.From = interfaceToString(mediaAttributes["From"])
			ma.Until = interfaceToString(mediaAttributes["Until"])
			result = append(result, ma)
		}
	}
	return result, nil
}

func (device Device) FindRecordings() (string, error) {
	// create soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body: `<FindRecordings xmlns="http://www.onvif.org/ver10/search/wsdl">
					</FindRecordings>`,
	}

	var result = ""
	// send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return result, err
	}

	// parse response
	data, err := response.ValueForPath("Envelope.Body.FindRecordingsResponse.SearchToken")
	if err != nil {
		return result, err
	}
	return interfaceToString(data), nil
}

func (device Device) GetRecordingSearchResults(searchToken string) (interface{}, error) {
	// create soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body: `<GetRecordingSearchResults xmlns="http://www.onvif.org/ver10/search/wsdl">
					<SearchToken>` + searchToken + `</SearchToken>
				</GetRecordingSearchResults>`,
	}

	var result interface{}
	// send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return result, err
	}

	// parse response
	data, err := response.ValueForPath("Envelope.Body.GetRecordingSearchResultsResponse")
	if err != nil {
		return result, err
	}
	glog.Infof("Data %v", data)
	return result, nil
}

func (device Device) FindEvents(startPoint string) (interface{}, error) {
	// create soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body: `<FindEvents xmlns="http://www.onvif.org/ver10/search/wsdl">
					<StartPoint>` + startPoint + `</StartPoint>
			   </FindEvents>`,
	}

	var result interface{}
	// send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return result, err
	}

	// parse response
	data, err := response.ValueForPath("Envelope.Body.FindEventsResponse")
	if err != nil {
		return result, err
	}
	glog.Infof("Data %v", data)
	return result, nil
}

func (device Device) GetEventSearchResults(searchToken string) (interface{}, error) {
	// create soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body: `<GetEventSearchResults xmlns="http://www.onvif.org/ver10/search/wsdl">
					<SearchToken>` + searchToken + `</SearchToken>
				</GetEventSearchResults>`,
	}

	var result interface{}
	// send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return result, err
	}

	// parse response
	data, err := response.ValueForPath("Envelope.Body.GetEventSearchResultsResponse")
	if err != nil {
		return result, err
	}
	glog.Infof("Data %v", data)
	return result, nil
}
