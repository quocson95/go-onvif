package onvif

// return url for unsubscribe
func (device Device) Subscribe(address string) (string, error) {
	// create soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		XMLNs: []string{
			`xmlns:wsnt="http://docs.oasis-open.org/wsn/b-2"`,
			`xmlns:wsa="http://www.w3.org/2005/08/addressing"`,
		},
		Body: `<wsnt:Subscribe xmlns="http://docs.oasis-open.org/wsn/b-2.xsd">
					<wsnt:ConsumerReference>
						<wsa:Address xmlns="http://www.w3.org/2005/08/addressing">` + address + `</wsa:Address>
					</wsnt:ConsumerReference>
					<wsnt:InitialTerminationTime>PT3600S</wsnt:InitialTerminationTime>
				</wsnt:Subscribe>`,
	}

	var result string = ""
	// send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return result, err
	}

	ifaceResult, err := response.ValueForPath("Envelope.Body.SubscribeResponse.SubscriptionReference.Address")
	if err != nil {
		return result, err
	}

	result, _ = ifaceResult.(string)
	return result, nil
}

func (device Device) CreatePullPointSubscription() (CreatePullPointSubscriptionResponse, error) {
	// create soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Action:   "http://www.onvif.org/ver10/events/wsdl/EventPortType/CreatePullPointSubscriptionRequest",
		Body: `<CreatePullPointSubscription xmlns="http://www.onvif.org/ver10/events/wsdl">
					<InitialTerminationTime>PT3600S</InitialTerminationTime>
				</CreatePullPointSubscription>`,
	}

	result := CreatePullPointSubscriptionResponse{}
	// send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return result, err
	}

	ifaceResult, err := response.ValueForPath("Envelope.Body.CreatePullPointSubscriptionResponse")
	if err != nil {
		return result, err
	}

	if mapMetadata, ok := ifaceResult.(map[string]interface{}); ok {
		result.CurrentTime = interfaceToString(mapMetadata["CurrentTime"])
		result.TerminationTime = interfaceToString(mapMetadata["TerminationTime"])
		result.SubscriptionReference = SubscriptionReference{}
		if mapSubscriptionReference, ok := mapMetadata["SubscriptionReference"].(map[string]interface{}); ok {
			result.SubscriptionReference.Address = interfaceToString(mapSubscriptionReference["Address"])
		}
	}

	return result, nil
}

// return url for unsubscribe
func (device Device) GetEventProperties() (interface{}, error) {
	// create soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body:     `<GetEventProperties xmlns="http://www.onvif.org/ver10/events/wsdl"/>`,
	}

	// send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return "", err
	}

	ifaceResult, err := response.ValueForPath("Envelope.Body.GetEventPropertiesResponse")
	if err != nil {
		return "", err
	}

	return ifaceResult, nil
}

// return url for unsubscribe
func (device Device) PullMessages(address string) ([]NotificationMessage, error) {
	// create soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Action:   "http://www.onvif.org/ver10/events/wsdl/PullPointSubscription/PullMessagesRequest",
		Body: `<PullMessages xmlns="http://www.onvif.org/ver10/events/wsdl">
					<Timeout>PT3S</Timeout>
					<MessageLimit>100</MessageLimit>
				</PullMessages>`,
	}

	var result = make([]NotificationMessage, 0)
	// send request
	response, err := soap.SendRequest(address)
	if err != nil {
		return result, err
	}

	ifaceResult, err := response.ValuesForPath("Envelope.Body.PullMessagesResponse.NotificationMessage")
	if err != nil {
		return result, err
	}

	for _, notificationMessage := range ifaceResult {
		if mapNotiMsg, ok := notificationMessage.(map[string]interface{}); ok {
			msg := NotificationMessage{}
			if mapTopic, ok := mapNotiMsg["Topic"].(map[string]interface{}); ok {
				msg.Topic = interfaceToString(mapTopic["#text"])
			}
			if mapMsg, ok := mapNotiMsg["Message"].(map[string]interface{}); ok {
				if mapMsg, ok := mapMsg["Message"].(map[string]interface{}); ok {
					msg.UtcTime = interfaceToString(mapMsg["-UtcTime"])
					if mapData, ok := mapMsg["Data"].(map[string]interface{}); ok {
						if mapSimpleItems, ok := mapData["SimpleItem"].([]interface{}); ok {
							for _, item := range mapSimpleItems {
								if mapItem, ok := item.(map[string]interface{}); ok {
									msg.Data = append(msg.Data, MessageData{
										Name:  interfaceToString(mapItem["-Name"]),
										Value: interfaceToString(mapItem["-Value"]),
									})
								}
							}
						} else if mapSimpleItem, ok := mapData["SimpleItem"].(map[string]interface{}); ok {
							msg.Data = append(msg.Data, MessageData{
								Name:  interfaceToString(mapSimpleItem["-Name"]),
								Value: interfaceToString(mapSimpleItem["-Value"]),
							})
						}
					}
					if mapSource, ok := mapMsg["Source"].(map[string]interface{}); ok {
						if mapSimpleItems, ok := mapSource["SimpleItem"].([]interface{}); ok {
							for _, item := range mapSimpleItems {
								if mapItem, ok := item.(map[string]interface{}); ok {
									msg.Source = append(msg.Source, MessageData{
										Name:  interfaceToString(mapItem["-Name"]),
										Value: interfaceToString(mapItem["-Value"]),
									})
								}
							}
						} else if mapSimpleItem, ok := mapSource["SimpleItem"].(map[string]interface{}); ok {
							msg.Source = append(msg.Source, MessageData{
								Name:  interfaceToString(mapSimpleItem["-Name"]),
								Value: interfaceToString(mapSimpleItem["-Value"]),
							})
						}
					}
				}
			}
			result = append(result, msg)
		}
	}
	return result, nil
}

func (device Device) UnSubscribe(address string) error {
	// create soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body:     `<Unsubscribe xmlns="http://docs.oasis-open.org/wsn/b-2"/>`,
	}
	// send request
	response, err := soap.SendRequest(address)
	if err != nil {
		return err
	}
	_, err = response.ValueForPath("Envelope.Body.UnsubscribeResponse")
	if err != nil {
		return err
	}

	return nil
}

func (device Device) ReNew(address string) (CreatePullPointSubscriptionResponse, error) {
	// create soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body: `<Renew xmlns="http://docs.oasis-open.org/wsn/b-2">
						<TerminationTime>PT3600S</TerminationTime>
					</Renew>`,
	}
	result := CreatePullPointSubscriptionResponse{}

	// send request
	response, err := soap.SendRequest(address)
	if err != nil {
		return result, err
	}

	ifaceResult, err := response.ValueForPath("Envelope.Body.RenewResponse")
	if err != nil {
		return result, err
	}

	if mapMetadata, ok := ifaceResult.(map[string]interface{}); ok {
		result.CurrentTime = interfaceToString(mapMetadata["CurrentTime"])
		result.TerminationTime = interfaceToString(mapMetadata["TerminationTime"])
	}
	return result, nil
}
