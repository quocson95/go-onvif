package onvif


// return url for unsubscribe
func (device Device) Subscribe(address string) (string, error) {
	// create soap
	soap := SOAP{
		User: device.User,
		Password: device.Password,
		Body: `<Subscribe xmlns="http://docs.oasis-open.org/wsn/b-2">
					<ConsumerReference>
						<Address>` + address + `</Address>
					</ConsumerReference>
				</Subscribe>`,
	}

	var result string = ""
	// send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil{
		return result, err
	}

	ifaceResult, err := response.ValueForPath("Envelope.Body.SubscribeResponse.SubscriptionReference.Address")
	if err != nil{
		return result, err
	}

	result,_ = ifaceResult.(string)
	return result, nil
}


func (device Device) UnSubscribe(address string) error {
	// create soap
	soap := SOAP{
		User: device.User,
		Password: device.Password,
		Body: `<Unsubscribe xmlns="http://docs.oasis-open.org/wsn/b-2"/>`,
	}
	// send request
	response, err := soap.SendRequest(address)
	if err != nil{
		return err
	}
	_, err = response.ValueForPath("Envelope.Body.UnsubscribeResponse")
	if err != nil{
		return err
	}

	return nil
}
