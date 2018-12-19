package onvif

func (device Device) Subscribe(address string, timeout string) (error) {
	// create soap
	soap := SOAP{
		User: device.User,
		Password: device.Password,
		Body: `<Subscribe xmlns="http://docs.oasis-open.org/wsn/b-2">
					<ConsumerReference>
						<Address>` + address + `</Address>
					</ConsumerReference>
					<InitialTerminationTime>` + timeout +`</InitialTerminationTime>
			   </Subscribe>`,
	}

	// send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil{
		return err
	}

	_, err = response.ValueForPath("Envelope.Body.SubscribeResponse")
	if err != nil{
		return err
	}

	return nil
}

