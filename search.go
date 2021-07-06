package onvif

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

func (device Device) GetRecordingSearchResults(searchToken string) (ResultList, error) {
	// create soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body: `<GetRecordingSearchResults xmlns="http://www.onvif.org/ver10/search/wsdl">
					<SearchToken>` + searchToken + `</SearchToken>
				</GetRecordingSearchResults>`,
	}

	result := ResultList{}
	// send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return result, err
	}

	// parse response
	resultListIf, err := response.ValueForPath("Envelope.Body.GetRecordingSearchResultsResponse.ResultList")
	if err != nil {
		return result, err
	}
	if resultList, ok := resultListIf.(map[string]interface{}); ok {
		result.SearchState = interfaceToString(resultList["SearchState"])
		result.RecordingInformation = make([]RecordingInformation, 0)
		if recordingInformationList, ok := resultList["RecordingInformation"].(map[string]interface{}); ok {
			recordingInformation := RecordingInformation{}
			recordingInformation.RecordingToken = interfaceToString(recordingInformationList["RecordingToken"])
			recordingInformation.EarliestRecording = interfaceToString(recordingInformationList["EarliestRecording"])
			recordingInformation.LatestRecording = interfaceToString(recordingInformationList["LatestRecording"])
			recordingInformation.RecordingStatus = interfaceToString(recordingInformationList["RecordingStatus"])
			recordingInformation.Track = make([]Track, 0)
			if listTrack, ok := recordingInformationList["Track"].([]interface{}); ok {
				for _, t := range listTrack {
					if track, ok := t.(map[string]interface{}); ok {
						recordingInformation.Track = append(recordingInformation.Track, Track{
							TrackToken:  interfaceToString(track["TrackToken"]),
							TrackType:   interfaceToString(track["TrackType"]),
							Description: interfaceToString(track["Description"]),
							DataFrom:    interfaceToString(track["DataFrom"]),
							DataTo:      interfaceToString(track["DataTo"]),
						})
					}
				}
			} else {
				if track, ok := recordingInformationList["Track"].(map[string]interface{}); ok {
					recordingInformation.Track = append(recordingInformation.Track, Track{
						TrackToken:  interfaceToString(track["TrackToken"]),
						TrackType:   interfaceToString(track["TrackType"]),
						Description: interfaceToString(track["Description"]),
						DataFrom:    interfaceToString(track["DataFrom"]),
						DataTo:      interfaceToString(track["DataTo"]),
					})
				}
			}
			result.RecordingInformation = append(result.RecordingInformation, recordingInformation)
		}
	}

	return result, nil
}

func (device Device) FindEvents(startPoint string) (string, error) {
	// create soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body: `<FindEvents xmlns="http://www.onvif.org/ver10/search/wsdl">
					<StartPoint>` + startPoint + `</StartPoint>
			   </FindEvents>`,
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

func (device Device) GetEventSearchResults(searchToken string) (ResultList, error) {
	// create soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body: `<GetEventSearchResults xmlns="http://www.onvif.org/ver10/search/wsdl">
					<SearchToken>` + searchToken + `</SearchToken>
				</GetEventSearchResults>`,
	}

	result := ResultList{}
	// send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return result, err
	}

	// parse response
	resultListIf, err := response.ValueForPath("Envelope.Body.GetEventSearchResultsResponse.ResultList")
	if err != nil {
		return result, err
	}
	if resultList, ok := resultListIf.(map[string]interface{}); ok {
		result.SearchState = interfaceToString(resultList["SearchState"])
		result.RecordingInformation = make([]RecordingInformation, 0)
		if recordingInformationList, ok := resultList["RecordingInformation"].(map[string]interface{}); ok {
			recordingInformation := RecordingInformation{}
			recordingInformation.RecordingToken = interfaceToString(recordingInformationList["RecordingToken"])
			recordingInformation.EarliestRecording = interfaceToString(recordingInformationList["EarliestRecording"])
			recordingInformation.LatestRecording = interfaceToString(recordingInformationList["LatestRecording"])
			recordingInformation.RecordingStatus = interfaceToString(recordingInformationList["RecordingStatus"])
			recordingInformation.Track = make([]Track, 0)
			if listTrack, ok := recordingInformationList["Track"].([]interface{}); ok {
				for _, t := range listTrack {
					if track, ok := t.(map[string]interface{}); ok {
						recordingInformation.Track = append(recordingInformation.Track, Track{
							TrackToken:  interfaceToString(track["TrackToken"]),
							TrackType:   interfaceToString(track["TrackType"]),
							Description: interfaceToString(track["Description"]),
							DataFrom:    interfaceToString(track["DataFrom"]),
							DataTo:      interfaceToString(track["DataTo"]),
						})
					}
				}
			} else {
				if track, ok := recordingInformationList["Track"].(map[string]interface{}); ok {
					recordingInformation.Track = append(recordingInformation.Track, Track{
						TrackToken:  interfaceToString(track["TrackToken"]),
						TrackType:   interfaceToString(track["TrackType"]),
						Description: interfaceToString(track["Description"]),
						DataFrom:    interfaceToString(track["DataFrom"]),
						DataTo:      interfaceToString(track["DataTo"]),
					})
				}
			}
			result.RecordingInformation = append(result.RecordingInformation, recordingInformation)
		}
	}

	return result, nil
}
