package onvif

var ptzXMLNs = []string{
	`xmlns:i="http://www.w3.org/2001/XMLSchema-instance"`,
	`xmlns:d="http://www.w3.org/2001/XMLSchema"`,
	`xmlns:c="http://www.w3.org/2003/05/soap-encoding"`,
}

func (device Device) GetNodes() ([]PTZNode, error) {
	//create soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body:     `<GetNodes xmlns="http://www.onvif.org/ver20/ptz/wsdl"/>`,
	}
	result := []PTZNode{}

	// send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return result, err
	}

	// parse response
	ifacePTZNodes, err := response.ValuesForPath("Envelope.Body.GetNodesResponse.PTZNode")
	if err != nil {
		return result, err
	}

	// parse interface
	for _, ifacePTZNode := range ifacePTZNodes {
		if mapPTZNode, ok := ifacePTZNode.(map[string]interface{}); ok {
			PTZNode := PTZNode{}

			PTZNode.Token = interfaceToString(mapPTZNode["-token"])
			PTZNode.Name = interfaceToString(mapPTZNode["Name"])
			PTZNode.FixedHomePosition = interfaceToBool(mapPTZNode["FixedHomePosition"])
			PTZNode.GeoMove = interfaceToBool(mapPTZNode["GeoMove"])
			PTZNode.MaximumNumberOfPresets = interfaceToInt(mapPTZNode["MaximumNumberOfPresets"])
			PTZNode.HomeSupported = interfaceToBool(mapPTZNode["HomeSupported"])

			// parse SupportedPTZSpaces
			if mapSupportedPTZSpaces, ok := mapPTZNode["SupportedPTZSpaces"].(map[string]interface{}); ok {
				SupportedPTZSpaces := PTZSpaces{}

				// parse AbsolutePanTiltPositionSpace
				if mapAbsolutePanTiltPositionSpace, ok := mapSupportedPTZSpaces["AbsolutePanTiltPositionSpace"].(map[string]interface{}); ok {
					AbsolutePanTiltPositionSpace := Space2DDescription{}

					AbsolutePanTiltPositionSpace.URI = interfaceToString(mapAbsolutePanTiltPositionSpace["URI"])
					// parse XRange
					if mapXRange, ok := mapAbsolutePanTiltPositionSpace["XRange"].(map[string]interface{}); ok {
						AbsolutePanTiltPositionSpace.XRange.Min = interfaceToFloat64(mapXRange["Min"])
						AbsolutePanTiltPositionSpace.XRange.Max = interfaceToFloat64(mapXRange["Max"])
					}
					// parse YRange
					if mapYRange, ok := mapAbsolutePanTiltPositionSpace["YRange"].(map[string]interface{}); ok {
						AbsolutePanTiltPositionSpace.YRange.Min = interfaceToFloat64(mapYRange["Min"])
						AbsolutePanTiltPositionSpace.YRange.Max = interfaceToFloat64(mapYRange["Max"])
					}

					SupportedPTZSpaces.AbsolutePanTiltPositionSpace = AbsolutePanTiltPositionSpace
				}

				// parse AbsoluteZoomPositionSpace
				if mapAbsoluteZoomPositionSpace, ok := mapSupportedPTZSpaces["AbsoluteZoomPositionSpace"].(map[string]interface{}); ok {
					AbsoluteZoomPositionSpace := Space1DDescription{}

					AbsoluteZoomPositionSpace.URI = interfaceToString(mapAbsoluteZoomPositionSpace["URI"])
					// parse XRange
					if mapXRange, ok := mapAbsoluteZoomPositionSpace["XRange"].(map[string]interface{}); ok {
						AbsoluteZoomPositionSpace.XRange.Min = interfaceToFloat64(mapXRange["Min"])
						AbsoluteZoomPositionSpace.XRange.Max = interfaceToFloat64(mapXRange["Max"])
					}

					SupportedPTZSpaces.AbsoluteZoomPositionSpace = AbsoluteZoomPositionSpace
				}

				// parse RelativePanTiltTranslationSpace
				if mapRelativePanTiltTranslationSpace, ok := mapSupportedPTZSpaces["RelativePanTiltTranslationSpace"].(map[string]interface{}); ok {
					RelativePanTiltTranslationSpace := Space2DDescription{}

					RelativePanTiltTranslationSpace.URI = interfaceToString(mapRelativePanTiltTranslationSpace["URI"])
					// parse XRange
					if mapXRange, ok := mapRelativePanTiltTranslationSpace["XRange"].(map[string]interface{}); ok {
						RelativePanTiltTranslationSpace.XRange.Min = interfaceToFloat64(mapXRange["Min"])
						RelativePanTiltTranslationSpace.XRange.Max = interfaceToFloat64(mapXRange["Max"])
					}
					// parse YRange
					if mapYRange, ok := mapRelativePanTiltTranslationSpace["YRange"].(map[string]interface{}); ok {
						RelativePanTiltTranslationSpace.YRange.Min = interfaceToFloat64(mapYRange["Min"])
						RelativePanTiltTranslationSpace.YRange.Max = interfaceToFloat64(mapYRange["Max"])
					}

					SupportedPTZSpaces.RelativePanTiltTranslationSpace = RelativePanTiltTranslationSpace
				}

				// parse RelativeZoomTranslationSpace
				if mapRelativeZoomTranslationSpace, ok := mapSupportedPTZSpaces["RelativeZoomTranslationSpace"].(map[string]interface{}); ok {
					RelativeZoomTranslationSpace := Space1DDescription{}

					RelativeZoomTranslationSpace.URI = interfaceToString(mapRelativeZoomTranslationSpace["URI"])
					// parse XRange
					if mapXRange, ok := mapRelativeZoomTranslationSpace["XRange"].(map[string]interface{}); ok {
						RelativeZoomTranslationSpace.XRange.Min = interfaceToFloat64(mapXRange["Min"])
						RelativeZoomTranslationSpace.XRange.Max = interfaceToFloat64(mapXRange["Max"])
					}

					SupportedPTZSpaces.RelativeZoomTranslationSpace = RelativeZoomTranslationSpace
				}

				// parse ContinuousPanTiltVelocitySpace
				if mapContinuousPanTiltVelocitySpace, ok := mapSupportedPTZSpaces["ContinuousPanTiltVelocitySpace"].(map[string]interface{}); ok {
					ContinuousPanTiltVelocitySpace := Space2DDescription{}

					ContinuousPanTiltVelocitySpace.URI = interfaceToString(mapContinuousPanTiltVelocitySpace["URI"])
					// parse XRange
					if mapXRange, ok := mapContinuousPanTiltVelocitySpace["XRange"].(map[string]interface{}); ok {
						ContinuousPanTiltVelocitySpace.XRange.Min = interfaceToFloat64(mapXRange["Min"])
						ContinuousPanTiltVelocitySpace.XRange.Max = interfaceToFloat64(mapXRange["Max"])
					}
					// parse YRange
					if mapYRange, ok := mapContinuousPanTiltVelocitySpace["YRange"].(map[string]interface{}); ok {
						ContinuousPanTiltVelocitySpace.YRange.Min = interfaceToFloat64(mapYRange["Min"])
						ContinuousPanTiltVelocitySpace.YRange.Max = interfaceToFloat64(mapYRange["Max"])
					}

					SupportedPTZSpaces.ContinuousPanTiltVelocitySpace = ContinuousPanTiltVelocitySpace
				}

				// parse ContinuousZoomVelocitySpace
				if mapContinuousZoomVelocitySpace, ok := mapSupportedPTZSpaces["ContinuousZoomVelocitySpace"].(map[string]interface{}); ok {
					ContinuousZoomVelocitySpace := Space1DDescription{}

					ContinuousZoomVelocitySpace.URI = interfaceToString(mapContinuousZoomVelocitySpace["URI"])
					// parse XRange
					if mapXRange, ok := mapContinuousZoomVelocitySpace["XRange"].(map[string]interface{}); ok {
						ContinuousZoomVelocitySpace.XRange.Min = interfaceToFloat64(mapXRange["Min"])
						ContinuousZoomVelocitySpace.XRange.Max = interfaceToFloat64(mapXRange["Max"])
					}

					SupportedPTZSpaces.ContinuousZoomVelocitySpace = ContinuousZoomVelocitySpace
				}

				// parse PanTiltSpeedSpace
				if mapPanTiltSpeedSpace, ok := mapSupportedPTZSpaces["PanTiltSpeedSpace"].(map[string]interface{}); ok {
					PanTiltSpeedSpace := Space1DDescription{}

					PanTiltSpeedSpace.URI = interfaceToString(mapPanTiltSpeedSpace["URI"])
					// parse XRange
					if mapXRange, ok := mapPanTiltSpeedSpace["XRange"].(map[string]interface{}); ok {
						PanTiltSpeedSpace.XRange.Min = interfaceToFloat64(mapXRange["Min"])
						PanTiltSpeedSpace.XRange.Max = interfaceToFloat64(mapXRange["Max"])
					}

					SupportedPTZSpaces.PanTiltSpeedSpace = PanTiltSpeedSpace
				}

				// parse ZoomSpeedSpace
				if mapZoomSpeedSpace, ok := mapSupportedPTZSpaces["ZoomSpeedSpace"].(map[string]interface{}); ok {
					ZoomSpeedSpace := Space1DDescription{}

					ZoomSpeedSpace.URI = interfaceToString(mapZoomSpeedSpace["URI"])
					// parse XRange
					if mapXRange, ok := mapZoomSpeedSpace["XRange"].(map[string]interface{}); ok {
						ZoomSpeedSpace.XRange.Min = interfaceToFloat64(mapXRange["Min"])
						ZoomSpeedSpace.XRange.Max = interfaceToFloat64(mapXRange["Max"])
					}

					SupportedPTZSpaces.ZoomSpeedSpace = ZoomSpeedSpace
				}

				PTZNode.SupportedPTZSpaces = SupportedPTZSpaces
			}

			// push into result
			result = append(result, PTZNode)
		}
	}

	return result, nil
}

func (device Device) GetNode(nodeToken string) (PTZNode, error) {
	// create soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body: `<GetNode xmlns="http://www.onvif.org/ver20/ptz/wsdl">
					<NodeToken>` + nodeToken + `</NodeToken>
				</GetNode>`,
	}

	result := PTZNode{}

	// send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return result, err
	}

	// parse response
	ifacePTZNode, err := response.ValueForPath("Envelope.Body.GetNodeResponse.PTZNode")
	if err != nil {
		return result, err
	}

	if mapPTZNode, ok := ifacePTZNode.(map[string]interface{}); ok {
		result.Token = interfaceToString(mapPTZNode["-token"])
		result.Name = interfaceToString(mapPTZNode["Name"])
		result.FixedHomePosition = interfaceToBool(mapPTZNode["FixedHomePosition"])
		result.GeoMove = interfaceToBool(mapPTZNode["GeoMove"])
		result.MaximumNumberOfPresets = interfaceToInt(mapPTZNode["MaximumNumberOfPresets"])
		result.HomeSupported = interfaceToBool(mapPTZNode["HomeSupported"])

		// parse SupportedPTZSpaces
		if mapSupportedPTZSpaces, ok := mapPTZNode["SupportedPTZSpaces"].(map[string]interface{}); ok {
			SupportedPTZSpaces := PTZSpaces{}

			// parse AbsolutePanTiltPositionSpace
			if mapAbsolutePanTiltPositionSpace, ok := mapSupportedPTZSpaces["AbsolutePanTiltPositionSpace"].(map[string]interface{}); ok {
				AbsolutePanTiltPositionSpace := Space2DDescription{}

				AbsolutePanTiltPositionSpace.URI = interfaceToString(mapAbsolutePanTiltPositionSpace["URI"])
				// parse XRange
				if mapXRange, ok := mapAbsolutePanTiltPositionSpace["XRange"].(map[string]interface{}); ok {
					AbsolutePanTiltPositionSpace.XRange.Min = interfaceToFloat64(mapXRange["Min"])
					AbsolutePanTiltPositionSpace.XRange.Max = interfaceToFloat64(mapXRange["Max"])
				}
				// parse YRange
				if mapYRange, ok := mapAbsolutePanTiltPositionSpace["YRange"].(map[string]interface{}); ok {
					AbsolutePanTiltPositionSpace.YRange.Min = interfaceToFloat64(mapYRange["Min"])
					AbsolutePanTiltPositionSpace.YRange.Max = interfaceToFloat64(mapYRange["Max"])
				}

				SupportedPTZSpaces.AbsolutePanTiltPositionSpace = AbsolutePanTiltPositionSpace
			}

			// parse AbsoluteZoomPositionSpace
			if mapAbsoluteZoomPositionSpace, ok := mapSupportedPTZSpaces["AbsoluteZoomPositionSpace"].(map[string]interface{}); ok {
				AbsoluteZoomPositionSpace := Space1DDescription{}

				AbsoluteZoomPositionSpace.URI = interfaceToString(mapAbsoluteZoomPositionSpace["URI"])
				// parse XRange
				if mapXRange, ok := mapAbsoluteZoomPositionSpace["XRange"].(map[string]interface{}); ok {
					AbsoluteZoomPositionSpace.XRange.Min = interfaceToFloat64(mapXRange["Min"])
					AbsoluteZoomPositionSpace.XRange.Max = interfaceToFloat64(mapXRange["Max"])
				}

				SupportedPTZSpaces.AbsoluteZoomPositionSpace = AbsoluteZoomPositionSpace
			}

			// parse RelativePanTiltTranslationSpace
			if mapRelativePanTiltTranslationSpace, ok := mapSupportedPTZSpaces["RelativePanTiltTranslationSpace"].(map[string]interface{}); ok {
				RelativePanTiltTranslationSpace := Space2DDescription{}

				RelativePanTiltTranslationSpace.URI = interfaceToString(mapRelativePanTiltTranslationSpace["URI"])
				// parse XRange
				if mapXRange, ok := mapRelativePanTiltTranslationSpace["XRange"].(map[string]interface{}); ok {
					RelativePanTiltTranslationSpace.XRange.Min = interfaceToFloat64(mapXRange["Min"])
					RelativePanTiltTranslationSpace.XRange.Max = interfaceToFloat64(mapXRange["Max"])
				}
				// parse YRange
				if mapYRange, ok := mapRelativePanTiltTranslationSpace["YRange"].(map[string]interface{}); ok {
					RelativePanTiltTranslationSpace.YRange.Min = interfaceToFloat64(mapYRange["Min"])
					RelativePanTiltTranslationSpace.YRange.Max = interfaceToFloat64(mapYRange["Max"])
				}

				SupportedPTZSpaces.RelativePanTiltTranslationSpace = RelativePanTiltTranslationSpace
			}

			// parse RelativeZoomTranslationSpace
			if mapRelativeZoomTranslationSpace, ok := mapSupportedPTZSpaces["RelativeZoomTranslationSpace"].(map[string]interface{}); ok {
				RelativeZoomTranslationSpace := Space1DDescription{}

				RelativeZoomTranslationSpace.URI = interfaceToString(mapRelativeZoomTranslationSpace["URI"])
				// parse XRange
				if mapXRange, ok := mapRelativeZoomTranslationSpace["XRange"].(map[string]interface{}); ok {
					RelativeZoomTranslationSpace.XRange.Min = interfaceToFloat64(mapXRange["Min"])
					RelativeZoomTranslationSpace.XRange.Max = interfaceToFloat64(mapXRange["Max"])
				}

				SupportedPTZSpaces.RelativeZoomTranslationSpace = RelativeZoomTranslationSpace
			}

			// parse ContinuousPanTiltVelocitySpace
			if mapContinuousPanTiltVelocitySpace, ok := mapSupportedPTZSpaces["ContinuousPanTiltVelocitySpace"].(map[string]interface{}); ok {
				ContinuousPanTiltVelocitySpace := Space2DDescription{}

				ContinuousPanTiltVelocitySpace.URI = interfaceToString(mapContinuousPanTiltVelocitySpace["URI"])
				// parse XRange
				if mapXRange, ok := mapContinuousPanTiltVelocitySpace["XRange"].(map[string]interface{}); ok {
					ContinuousPanTiltVelocitySpace.XRange.Min = interfaceToFloat64(mapXRange["Min"])
					ContinuousPanTiltVelocitySpace.XRange.Max = interfaceToFloat64(mapXRange["Max"])
				}
				// parse YRange
				if mapYRange, ok := mapContinuousPanTiltVelocitySpace["YRange"].(map[string]interface{}); ok {
					ContinuousPanTiltVelocitySpace.YRange.Min = interfaceToFloat64(mapYRange["Min"])
					ContinuousPanTiltVelocitySpace.YRange.Max = interfaceToFloat64(mapYRange["Max"])
				}

				SupportedPTZSpaces.ContinuousPanTiltVelocitySpace = ContinuousPanTiltVelocitySpace
			}

			// parse ContinuousZoomVelocitySpace
			if mapContinuousZoomVelocitySpace, ok := mapSupportedPTZSpaces["ContinuousZoomVelocitySpace"].(map[string]interface{}); ok {
				ContinuousZoomVelocitySpace := Space1DDescription{}

				ContinuousZoomVelocitySpace.URI = interfaceToString(mapContinuousZoomVelocitySpace["URI"])
				// parse XRange
				if mapXRange, ok := mapContinuousZoomVelocitySpace["XRange"].(map[string]interface{}); ok {
					ContinuousZoomVelocitySpace.XRange.Min = interfaceToFloat64(mapXRange["Min"])
					ContinuousZoomVelocitySpace.XRange.Max = interfaceToFloat64(mapXRange["Max"])
				}

				SupportedPTZSpaces.ContinuousZoomVelocitySpace = ContinuousZoomVelocitySpace
			}

			// parse PanTiltSpeedSpace
			if mapPanTiltSpeedSpace, ok := mapSupportedPTZSpaces["PanTiltSpeedSpace"].(map[string]interface{}); ok {
				PanTiltSpeedSpace := Space1DDescription{}

				PanTiltSpeedSpace.URI = interfaceToString(mapPanTiltSpeedSpace["URI"])
				// parse XRange
				if mapXRange, ok := mapPanTiltSpeedSpace["XRange"].(map[string]interface{}); ok {
					PanTiltSpeedSpace.XRange.Min = interfaceToFloat64(mapXRange["Min"])
					PanTiltSpeedSpace.XRange.Max = interfaceToFloat64(mapXRange["Max"])
				}

				SupportedPTZSpaces.PanTiltSpeedSpace = PanTiltSpeedSpace
			}

			// parse ZoomSpeedSpace
			if mapZoomSpeedSpace, ok := mapSupportedPTZSpaces["ZoomSpeedSpace"].(map[string]interface{}); ok {
				ZoomSpeedSpace := Space1DDescription{}

				ZoomSpeedSpace.URI = interfaceToString(mapZoomSpeedSpace["URI"])
				// parse XRange
				if mapXRange, ok := mapZoomSpeedSpace["XRange"].(map[string]interface{}); ok {
					ZoomSpeedSpace.XRange.Min = interfaceToFloat64(mapXRange["Min"])
					ZoomSpeedSpace.XRange.Max = interfaceToFloat64(mapXRange["Max"])
				}

				SupportedPTZSpaces.ZoomSpeedSpace = ZoomSpeedSpace
			}

			result.SupportedPTZSpaces = SupportedPTZSpaces
		}
	}

	return result, nil
}

func (device Device) GetConfigurations() ([]PTZConfiguration, error) {
	// create soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body:     `<GetConfigurations xmlns="http://www.onvif.org/ver20/ptz/wsdl"/>`,
	}

	result := []PTZConfiguration{}

	//send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return result, err
	}

	// parse response
	ifacePTZConfigurations, err := response.ValuesForPath("Envelope.Body.GetConfigurationsResponse.PTZConfiguration")
	if err != nil {
		return result, err
	}

	// parse interface
	for _, ifacePTZConfiguration := range ifacePTZConfigurations {
		if mapPTZConfiguration, ok := ifacePTZConfiguration.(map[string]interface{}); ok {
			PTZConfiguration := PTZConfiguration{}

			PTZConfiguration.Token = interfaceToString(mapPTZConfiguration["-token"])
			PTZConfiguration.Name = interfaceToString(mapPTZConfiguration["Name"])
			PTZConfiguration.MoveRamp = interfaceToInt(mapPTZConfiguration["MoveRamp"])
			PTZConfiguration.PresetRamp = interfaceToInt(mapPTZConfiguration["PresetRamp"])
			PTZConfiguration.PresetTourRamp = interfaceToInt(mapPTZConfiguration["PresetTourRamp"])
			PTZConfiguration.NodeToken = interfaceToString(mapPTZConfiguration["NodeToken"])

			PTZConfiguration.DefaultAbsolutePantTiltPositionSpace = interfaceToString(mapPTZConfiguration["DefaultAbsolutePantTiltPositionSpace"])
			PTZConfiguration.DefaultAbsoluteZoomPositionSpace = interfaceToString(mapPTZConfiguration["DefaultAbsoluteZoomPositionSpace"])
			PTZConfiguration.DefaultRelativePanTiltTranslationSpace = interfaceToString(mapPTZConfiguration["DefaultRelativePanTiltTranslationSpace"])
			PTZConfiguration.DefaultRelativeZoomTranslationSpace = interfaceToString(mapPTZConfiguration["DefaultRelativeZoomTranslationSpace"])
			PTZConfiguration.DefaultContinuousPanTiltVelocitySpace = interfaceToString(mapPTZConfiguration["DefaultContinuousPanTiltVelocitySpace"])
			PTZConfiguration.DefaultContinuousZoomVelocitySpace = interfaceToString(mapPTZConfiguration["DefaultContinuousZoomVelocitySpace"])

			// parse DefaultPTZSpeed
			if mapDefaultPTZSpeed, ok := mapPTZConfiguration["DefaultPTZSpeed"].(map[string]interface{}); ok {
				DefaultPTZSpeed := PTZVector{}
				// parse PanTilt
				if mapPanTilt, ok := mapDefaultPTZSpeed["PanTilt"].(map[string]interface{}); ok {
					DefaultPTZSpeed.PanTilt.Space = interfaceToString(mapPanTilt["-space"])
					DefaultPTZSpeed.PanTilt.X = interfaceToFloat64(mapPanTilt["-x"])
					DefaultPTZSpeed.PanTilt.Y = interfaceToFloat64(mapPanTilt["-y"])
				}
				// parse Zoom
				if mapZoom, ok := mapDefaultPTZSpeed["Zoom"].(map[string]interface{}); ok {
					DefaultPTZSpeed.Zoom.Space = interfaceToString(mapZoom["-space"])
					DefaultPTZSpeed.Zoom.X = interfaceToFloat64(mapZoom["-x"])
				}
				PTZConfiguration.DefaultPTZSpeed = DefaultPTZSpeed
			}

			PTZConfiguration.DefaultPTZTimeout = interfaceToString(mapPTZConfiguration["DefaultPTZTimeout"])
			// parse PanTiltLimits
			if mapPanTiltLimits, ok := mapPTZConfiguration["PanTiltLimits"].(map[string]interface{}); ok {
				if mapRange, ok := mapPanTiltLimits["Range"].(map[string]interface{}); ok {
					PanTiltLimits := PanTiltLimits{}

					PanTiltLimits.Range.URI = interfaceToString(mapRange["URI"])
					// parse XRange
					if mapXRange, ok := mapRange["XRange"].(map[string]interface{}); ok {
						PanTiltLimits.Range.XRange.Min = interfaceToFloat64(mapXRange["Min"])
						PanTiltLimits.Range.XRange.Max = interfaceToFloat64(mapXRange["Max"])
					}
					// parse YRange
					if mapYRange, ok := mapRange["YRange"].(map[string]interface{}); ok {
						PanTiltLimits.Range.YRange.Min = interfaceToFloat64(mapYRange["Min"])
						PanTiltLimits.Range.YRange.Max = interfaceToFloat64(mapYRange["Max"])
					}

					PTZConfiguration.PanTiltLimits = PanTiltLimits
				}
			}

			// parse ZoomLimits
			if mapZoomLimits, ok := mapPTZConfiguration["ZoomLimits"].(map[string]interface{}); ok {
				if mapRange, ok := mapZoomLimits["Range"].(map[string]interface{}); ok {
					ZoomLimits := ZoomLimits{}

					ZoomLimits.Range.URI = interfaceToString(mapRange["URI"])
					// parse XRange
					if mapXRange, ok := mapRange["XRange"].(map[string]interface{}); ok {
						ZoomLimits.Range.XRange.Min = interfaceToFloat64(mapXRange["Min"])
						ZoomLimits.Range.XRange.Max = interfaceToFloat64(mapXRange["Max"])
					}

					PTZConfiguration.ZoomLimits = ZoomLimits
				}
			}

			// push into result
			result = append(result, PTZConfiguration)
		}
	}
	return result, nil
}

func (device Device) GetConfiguration(ptzConfigurationToken string) (PTZConfiguration, error) {
	// create soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body: `<GetConfiguration xmlns="http://www.onvif.org/ver20/ptz/wsdl">
					<PTZConfigurationToken>` + ptzConfigurationToken + `</PTZConfigurationToken>
				</GetConfiguration>`,
	}

	result := PTZConfiguration{}

	// send response
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return result, err
	}

	// parse response
	ifacePTZConfiguration, err := response.ValueForPath("Envelope.Body.GetConfigurationResponse.PTZConfiguration")
	if err != nil {
		return result, err
	}

	// parse interface
	if mapPTZConfiguration, ok := ifacePTZConfiguration.(map[string]interface{}); ok {
		result.Token = interfaceToString(mapPTZConfiguration["-token"])
		result.Name = interfaceToString(mapPTZConfiguration["Name"])
		result.MoveRamp = interfaceToInt(mapPTZConfiguration["MoveRamp"])
		result.PresetRamp = interfaceToInt(mapPTZConfiguration["PresetRamp"])
		result.PresetTourRamp = interfaceToInt(mapPTZConfiguration["PresetTourRamp"])
		result.NodeToken = interfaceToString(mapPTZConfiguration["NodeToken"])

		result.DefaultAbsolutePantTiltPositionSpace = interfaceToString(mapPTZConfiguration["DefaultAbsolutePantTiltPositionSpace"])
		result.DefaultAbsoluteZoomPositionSpace = interfaceToString(mapPTZConfiguration["DefaultAbsoluteZoomPositionSpace"])
		result.DefaultRelativePanTiltTranslationSpace = interfaceToString(mapPTZConfiguration["DefaultRelativePanTiltTranslationSpace"])
		result.DefaultRelativeZoomTranslationSpace = interfaceToString(mapPTZConfiguration["DefaultRelativeZoomTranslationSpace"])
		result.DefaultContinuousPanTiltVelocitySpace = interfaceToString(mapPTZConfiguration["DefaultContinuousPanTiltVelocitySpace"])
		result.DefaultContinuousZoomVelocitySpace = interfaceToString(mapPTZConfiguration["DefaultContinuousZoomVelocitySpace"])

		// parse DefaultPTZSpeed
		if mapDefaultPTZSpeed, ok := mapPTZConfiguration["DefaultPTZSpeed"].(map[string]interface{}); ok {
			DefaultPTZSpeed := PTZVector{}
			// parse PanTilt
			if mapPanTilt, ok := mapDefaultPTZSpeed["PanTilt"].(map[string]interface{}); ok {
				DefaultPTZSpeed.PanTilt.Space = interfaceToString(mapPanTilt["-space"])
				DefaultPTZSpeed.PanTilt.X = interfaceToFloat64(mapPanTilt["-x"])
				DefaultPTZSpeed.PanTilt.Y = interfaceToFloat64(mapPanTilt["-y"])
			}
			// parse Zoom
			if mapZoom, ok := mapDefaultPTZSpeed["Zoom"].(map[string]interface{}); ok {
				DefaultPTZSpeed.Zoom.Space = interfaceToString(mapZoom["-space"])
				DefaultPTZSpeed.Zoom.X = interfaceToFloat64(mapZoom["-x"])
			}
			result.DefaultPTZSpeed = DefaultPTZSpeed
		}

		result.DefaultPTZTimeout = interfaceToString(mapPTZConfiguration["DefaultPTZTimeout"])
		// parse PanTiltLimits
		if mapPanTiltLimits, ok := mapPTZConfiguration["PanTiltLimits"].(map[string]interface{}); ok {
			if mapRange, ok := mapPanTiltLimits["Range"].(map[string]interface{}); ok {
				PanTiltLimits := PanTiltLimits{}

				PanTiltLimits.Range.URI = interfaceToString(mapRange["URI"])
				// parse XRange
				if mapXRange, ok := mapRange["XRange"].(map[string]interface{}); ok {
					PanTiltLimits.Range.XRange.Min = interfaceToFloat64(mapXRange["Min"])
					PanTiltLimits.Range.XRange.Max = interfaceToFloat64(mapXRange["Max"])
				}
				// parse YRange
				if mapYRange, ok := mapRange["YRange"].(map[string]interface{}); ok {
					PanTiltLimits.Range.YRange.Min = interfaceToFloat64(mapYRange["Min"])
					PanTiltLimits.Range.YRange.Max = interfaceToFloat64(mapYRange["Max"])
				}

				result.PanTiltLimits = PanTiltLimits
			}
		}

		// parse ZoomLimits
		if mapZoomLimits, ok := mapPTZConfiguration["ZoomLimits"].(map[string]interface{}); ok {
			if mapRange, ok := mapZoomLimits["Range"].(map[string]interface{}); ok {
				ZoomLimits := ZoomLimits{}

				ZoomLimits.Range.URI = interfaceToString(mapRange["URI"])
				// parse XRange
				if mapXRange, ok := mapRange["XRange"].(map[string]interface{}); ok {
					ZoomLimits.Range.XRange.Min = interfaceToFloat64(mapXRange["Min"])
					ZoomLimits.Range.XRange.Max = interfaceToFloat64(mapXRange["Max"])
				}

				result.ZoomLimits = ZoomLimits
			}
		}
	}

	return result, nil
}

func (device Device) GetConfigurationOptions(configurationToken string) (PTZConfigurationOptions, error) {
	// create soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body: `<GetConfigurationOptions xmlns="http://www.onvif.org/ver20/ptz/wsdl">
					<ConfigurationToken>` + configurationToken + `</ConfigurationToken>
				</GetConfigurationOptions>`,
	}

	result := PTZConfigurationOptions{}

	// send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return result, err
	}

	// parse response
	ifacePTZConfigurationOptions, err := response.ValueForPath("Envelope.Body.GetConfigurationOptionsResponse.PTZConfigurationOptions")
	if err != nil {
		return result, err
	}

	// parse interface
	if mapPTZConfigurationOptions, ok := ifacePTZConfigurationOptions.(map[string]interface{}); ok {
		// parse Spaces
		if mapSpaces, ok := mapPTZConfigurationOptions["Spaces"].(map[string]interface{}); ok {
			Spaces := PTZSpaces{}

			// parse AbsolutePanTiltPositionSpace
			if mapAbsolutePanTiltPositionSpace, ok := mapSpaces["AbsolutePanTiltPositionSpace"].(map[string]interface{}); ok {
				AbsolutePanTiltPositionSpace := Space2DDescription{}

				AbsolutePanTiltPositionSpace.URI = interfaceToString(mapAbsolutePanTiltPositionSpace["URI"])
				// parse XRange
				if mapXRange, ok := mapAbsolutePanTiltPositionSpace["XRange"].(map[string]interface{}); ok {
					AbsolutePanTiltPositionSpace.XRange.Min = interfaceToFloat64(mapXRange["Min"])
					AbsolutePanTiltPositionSpace.XRange.Max = interfaceToFloat64(mapXRange["Max"])
				}
				// parse YRange
				if mapYRange, ok := mapAbsolutePanTiltPositionSpace["YRange"].(map[string]interface{}); ok {
					AbsolutePanTiltPositionSpace.YRange.Min = interfaceToFloat64(mapYRange["Min"])
					AbsolutePanTiltPositionSpace.YRange.Max = interfaceToFloat64(mapYRange["Max"])
				}

				Spaces.AbsolutePanTiltPositionSpace = AbsolutePanTiltPositionSpace
			}

			// parse AbsoluteZoomPositionSpace
			if mapAbsoluteZoomPositionSpace, ok := mapSpaces["AbsoluteZoomPositionSpace"].(map[string]interface{}); ok {
				AbsoluteZoomPositionSpace := Space1DDescription{}

				AbsoluteZoomPositionSpace.URI = interfaceToString(mapAbsoluteZoomPositionSpace["URI"])
				// parse XRange
				if mapXRange, ok := mapAbsoluteZoomPositionSpace["XRange"].(map[string]interface{}); ok {
					AbsoluteZoomPositionSpace.XRange.Min = interfaceToFloat64(mapXRange["Min"])
					AbsoluteZoomPositionSpace.XRange.Max = interfaceToFloat64(mapXRange["Max"])
				}

				Spaces.AbsoluteZoomPositionSpace = AbsoluteZoomPositionSpace
			}

			// parse RelativePanTiltTranslationSpace
			if mapRelativePanTiltTranslationSpace, ok := mapSpaces["RelativePanTiltTranslationSpace"].(map[string]interface{}); ok {
				RelativePanTiltTranslationSpace := Space2DDescription{}

				RelativePanTiltTranslationSpace.URI = interfaceToString(mapRelativePanTiltTranslationSpace["URI"])
				// parse XRange
				if mapXRange, ok := mapRelativePanTiltTranslationSpace["XRange"].(map[string]interface{}); ok {
					RelativePanTiltTranslationSpace.XRange.Min = interfaceToFloat64(mapXRange["Min"])
					RelativePanTiltTranslationSpace.XRange.Max = interfaceToFloat64(mapXRange["Max"])
				}
				// parse YRange
				if mapYRange, ok := mapRelativePanTiltTranslationSpace["YRange"].(map[string]interface{}); ok {
					RelativePanTiltTranslationSpace.YRange.Min = interfaceToFloat64(mapYRange["Min"])
					RelativePanTiltTranslationSpace.YRange.Max = interfaceToFloat64(mapYRange["Max"])
				}

				Spaces.RelativePanTiltTranslationSpace = RelativePanTiltTranslationSpace
			}

			// parse RelativeZoomTranslationSpace
			if mapRelativeZoomTranslationSpace, ok := mapSpaces["RelativeZoomTranslationSpace"].(map[string]interface{}); ok {
				RelativeZoomTranslationSpace := Space1DDescription{}

				RelativeZoomTranslationSpace.URI = interfaceToString(mapRelativeZoomTranslationSpace["URI"])
				// parse XRange
				if mapXRange, ok := mapRelativeZoomTranslationSpace["XRange"].(map[string]interface{}); ok {
					RelativeZoomTranslationSpace.XRange.Min = interfaceToFloat64(mapXRange["Min"])
					RelativeZoomTranslationSpace.XRange.Max = interfaceToFloat64(mapXRange["Max"])
				}

				Spaces.RelativeZoomTranslationSpace = RelativeZoomTranslationSpace
			}

			// parse ContinuousPanTiltVelocitySpace
			if mapContinuousPanTiltVelocitySpace, ok := mapSpaces["ContinuousPanTiltVelocitySpace"].(map[string]interface{}); ok {
				ContinuousPanTiltVelocitySpace := Space2DDescription{}

				ContinuousPanTiltVelocitySpace.URI = interfaceToString(mapContinuousPanTiltVelocitySpace["URI"])
				// parse XRange
				if mapXRange, ok := mapContinuousPanTiltVelocitySpace["XRange"].(map[string]interface{}); ok {
					ContinuousPanTiltVelocitySpace.XRange.Min = interfaceToFloat64(mapXRange["Min"])
					ContinuousPanTiltVelocitySpace.XRange.Max = interfaceToFloat64(mapXRange["Max"])
				}
				// parse YRange
				if mapYRange, ok := mapContinuousPanTiltVelocitySpace["YRange"].(map[string]interface{}); ok {
					ContinuousPanTiltVelocitySpace.YRange.Min = interfaceToFloat64(mapYRange["Min"])
					ContinuousPanTiltVelocitySpace.YRange.Max = interfaceToFloat64(mapYRange["Max"])
				}

				Spaces.ContinuousPanTiltVelocitySpace = ContinuousPanTiltVelocitySpace
			}

			// parse ContinuousZoomVelocitySpace
			if mapContinuousZoomVelocitySpace, ok := mapSpaces["ContinuousZoomVelocitySpace"].(map[string]interface{}); ok {
				ContinuousZoomVelocitySpace := Space1DDescription{}

				ContinuousZoomVelocitySpace.URI = interfaceToString(mapContinuousZoomVelocitySpace["URI"])
				// parse XRange
				if mapXRange, ok := mapContinuousZoomVelocitySpace["XRange"].(map[string]interface{}); ok {
					ContinuousZoomVelocitySpace.XRange.Min = interfaceToFloat64(mapXRange["Min"])
					ContinuousZoomVelocitySpace.XRange.Max = interfaceToFloat64(mapXRange["Max"])
				}

				Spaces.ContinuousZoomVelocitySpace = ContinuousZoomVelocitySpace
			}

			// parse PanTiltSpeedSpace
			if mapPanTiltSpeedSpace, ok := mapSpaces["PanTiltSpeedSpace"].(map[string]interface{}); ok {
				PanTiltSpeedSpace := Space1DDescription{}

				PanTiltSpeedSpace.URI = interfaceToString(mapPanTiltSpeedSpace["URI"])
				// parse XRange
				if mapXRange, ok := mapPanTiltSpeedSpace["XRange"].(map[string]interface{}); ok {
					PanTiltSpeedSpace.XRange.Min = interfaceToFloat64(mapXRange["Min"])
					PanTiltSpeedSpace.XRange.Max = interfaceToFloat64(mapXRange["Max"])
				}

				Spaces.PanTiltSpeedSpace = PanTiltSpeedSpace
			}

			// parse ZoomSpeedSpace
			if mapZoomSpeedSpace, ok := mapSpaces["ZoomSpeedSpace"].(map[string]interface{}); ok {
				ZoomSpeedSpace := Space1DDescription{}

				ZoomSpeedSpace.URI = interfaceToString(mapZoomSpeedSpace["URI"])
				// parse XRange
				if mapXRange, ok := mapZoomSpeedSpace["XRange"].(map[string]interface{}); ok {
					ZoomSpeedSpace.XRange.Min = interfaceToFloat64(mapXRange["Min"])
					ZoomSpeedSpace.XRange.Max = interfaceToFloat64(mapXRange["Max"])
				}

				Spaces.ZoomSpeedSpace = ZoomSpeedSpace
			}

			result.Spaces = Spaces
		}

		// parse PTZTimeout
		if mapPTZTimeout, ok := mapPTZConfigurationOptions["PTZTimeout"].(map[string]interface{}); ok {
			result.PTZTimeout.Min = interfaceToString(mapPTZTimeout["Min"])
			result.PTZTimeout.Max = interfaceToString(mapPTZTimeout["Max"])
		}

		// parse PTControlDirection
		if mapPTControlDirection, ok := mapPTZConfigurationOptions["PTControlDirection"].(map[string]interface{}); ok {
			// parse EFlip
			if mapEFlip, ok := mapPTControlDirection["EFlip"].(map[string]interface{}); ok {
				result.PTControlDirection.EFlip.Mode = interfaceToString(mapEFlip["Mode"])
			}

			// parse Reverse
			if mapReverse, ok := mapPTControlDirection["Reverse"].(map[string]interface{}); ok {
				result.PTControlDirection.Reverse.Mode = interfaceToString(mapReverse["Mode"])
			}
		}
	}

	return result, nil
}

func (device Device) GetStatus(profileToken string) (PTZStatus, error) {
	// create soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body: `<GetStatus xmlns="http://www.onvif.org/ver20/ptz/wsdl">
					<ProfileToken>` + profileToken + `</ProfileToken>
				</GetStatus>`,
	}

	result := PTZStatus{}

	//send soap
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return result, err
	}

	// parse response
	ifacePTZStatus, err := response.ValueForPath("Envelope.Body.GetStatusResponse.PTZStatus")
	if err != nil {
		return result, err
	}

	// parse interface
	if mapPTZStatus, ok := ifacePTZStatus.(map[string]interface{}); ok {
		// parse position
		if mapPosition, ok := mapPTZStatus["Position"].(map[string]interface{}); ok {
			Position := PTZVector{}

			//parse PanTilt
			if mapPanTilt, ok := mapPosition["PanTilt"].(map[string]interface{}); ok {
				Position.PanTilt.Space = interfaceToString(mapPanTilt["-space"])
				Position.PanTilt.X = interfaceToFloat64(mapPanTilt["-x"])
				Position.PanTilt.Y = interfaceToFloat64(mapPanTilt["-y"])
			}
			//parse Zoom
			if mapZoom, ok := mapPosition["Zoom"].(map[string]interface{}); ok {
				Position.Zoom.Space = interfaceToString(mapZoom["-space"])
				Position.Zoom.X = interfaceToFloat64(mapZoom["-x"])
			}
			result.Position = Position
		}

		// parse Move Status
		if mapMoveStatus, ok := mapPTZStatus["MoveStatus"].(map[string]interface{}); ok {
			result.MoveStatus.PanTilt = interfaceToString(mapMoveStatus["PanTilt"])
			result.MoveStatus.Zoom = interfaceToString(mapMoveStatus["Zoom"])
		}

		result.UtcTime = interfaceToString(mapPTZStatus["UtcTime"])
	}

	return result, nil
}

func (device Device) ContinuousMove(profileToken string, velocity PTZVector) error {
	// create soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body: `<ContinuousMove xmlns="http://www.onvif.org/ver20/ptz/wsdl">
					<ProfileToken>` + profileToken + `</ProfileToken>
					<Velocity>
						<PanTilt xmlns="http://www.onvif.org/ver10/schema" x="` + float64ToString(velocity.PanTilt.X) + `" y="` + float64ToString(velocity.PanTilt.Y) + `"/>
						<Zoom xmlns="http://www.onvif.org/ver10/schema" x="` + float64ToString(velocity.Zoom.X) + `"/>
					</Velocity>
				</ContinuousMove>`,
		XMLNs:  ptzXMLNs,
		Action: "http://www.onvif.org/ver20/ptz/wsdl/ContinuousMove",
	}

	//send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return err
	}

	_, err = response.ValueForPath("Envelope.Body.ContinuousMoveResponse")
	if err != nil {
		return err
	}

	return nil
}

func (device Device) AbsoluteMove(profileToken string, position PTZVector) error {
	// create soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body: `<AbsoluteMove xmlns="http://www.onvif.org/ver20/ptz/wsdl">
					<ProfileToken>` + profileToken + `</ProfileToken>
					<Position>
						<PanTilt xmlns="http://www.onvif.org/ver10/schema" x="` + float64ToString(position.PanTilt.X) + `" y="` + float64ToString(position.PanTilt.Y) + `"/>
						<Zoom x="` + float64ToString(position.Zoom.X) + `"/>
					</Position>
				</AbsoluteMove>`,
		XMLNs:  ptzXMLNs,
		Action: "http://www.onvif.org/ver20/ptz/wsdl/AbsoluteMove",
	}

	//send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return err
	}

	_, err = response.ValueForPath("Envelope.Body.AbsoluteMoveResponse")
	if err != nil {
		return err
	}

	return nil
}

/// PTZ Control RPC
// x: positive => go to right || negative => go to left
// y: positive => go to up || negative => go to down
// z: positive => zoom in || negative => zoom out
func (device Device) RelativeMove(profileToken string, translation PTZVector) error {
	// create soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body: `<RelativeMove xmlns="http://www.onvif.org/ver20/ptz/wsdl">
					<ProfileToken>` + profileToken + `</ProfileToken>
					<Translation>
						<PanTilt xmlns="http://www.onvif.org/ver10/schema" x="` + float64ToString(translation.PanTilt.X) + `" y="` + float64ToString(translation.PanTilt.Y) + `"/>
						<Zoom x="` + float64ToString(translation.Zoom.X) + `"/>
					</Translation>
				</RelativeMove>`,
		XMLNs:  ptzXMLNs,
		Action: "http://www.onvif.org/ver20/ptz/wsdl/RelativeMove",
	}

	//send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return err
	}

	_, err = response.ValueForPath("Envelope.Body.RelativeMoveResponse")
	if err != nil {
		return err
	}

	return nil
}

func (device Device) Stop(profileToken string) error {
	//create soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body: `<Stop xmlns="http://www.onvif.org/ver20/ptz/wsdl">
				<ProfileToken>` + profileToken + `</ProfileToken>
				<PanTilt>true</PanTilt><Zoom>true</Zoom>
			  </Stop>`,
		XMLNs:  ptzXMLNs,
		Action: "http://www.onvif.org/ver20/ptz/wsdl/Stop",
	}

	//send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return err
	}

	_, err = response.ValueForPath("Envelope.Body.StopResponse")
	if err != nil {
		return err
	}

	return nil
}

func (device Device) GotoHomePosition(profileToken string) error {
	// create soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body: `<GotoHomePosition xmlns="http://www.onvif.org/ver20/ptz/wsdl">
					<ProfileToken>` + profileToken + `</ProfileToken>
				</GotoHomePosition>`,
		XMLNs:  ptzXMLNs,
		Action: "http://www.onvif.org/ver20/ptz/wsdl/GotoHomePosition",
	}

	// send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return err
	}

	_, err = response.ValueForPath("Envelope.Body.GotoHomePositionResponse")
	if err != nil {
		return err
	}

	return nil
}

func (device Device) SetHomePosition(profileToken string) error {
	// create soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body: `<SetHomePosition xmlns="http://www.onvif.org/ver20/ptz/wsdl">
					<ProfileToken>` + profileToken + `</ProfileToken>
				</SetHomePosition>`,
	}

	//send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return err
	}

	_, err = response.ValueForPath("Envelope.Body.SetHomePositionResponse")
	if err != nil {
		return err
	}

	return nil
}

// return preset token of new preset
func (device Device) SetPreset(profileToken string, presetName string) (string, error) {
	// create soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body: `<SetPreset xmlns="http://www.onvif.org/ver20/ptz/wsdl">
					<ProfileToken>` + profileToken + `</ProfileToken>
					<PresetName>` + presetName + `</PresetName>
				</SetPreset>`,
	}
	var result string
	// send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return result, err
	}

	// parse response
	ifacePresetToken, err := response.ValueForPath("Envelope.Body.SetPresetResponse")
	if err != nil {
		return result, err
	}

	if mapProfileToken, ok := ifacePresetToken.(map[string]interface{}); ok {
		result = interfaceToString(mapProfileToken["PresetToken"])
	}

	return result, err
}

func (device Device) GetPresets(profileToken string) ([]PTZPreset, error) {
	// create soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body: `<GetPresets xmlns="http://www.onvif.org/ver20/ptz/wsdl">
					<ProfileToken>` + profileToken + `</ProfileToken>
				</GetPresets>`,
	}
	result := []PTZPreset{}

	// send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return result, err
	}

	// parse response
	ifacePresets, err := response.ValuesForPath("Envelope.Body.GetPresetsResponse.Preset")
	if err != nil {
		return result, err
	}

	// parse interface
	for _, ifacePreset := range ifacePresets {
		if mapPreset, ok := ifacePreset.(map[string]interface{}); ok {
			Preset := PTZPreset{}

			Preset.Token = interfaceToString(mapPreset["-token"])
			Preset.Name = interfaceToString(mapPreset["Name"])
			// parse PTZPosition
			if mapPosition, ok := mapPreset["PTZPosition"].(map[string]interface{}); ok {
				Position := PTZVector{}
				// parse PanTilt
				if mapPanTilt, ok := mapPosition["PanTilt"].(map[string]interface{}); ok {
					Position.PanTilt.Space = interfaceToString(mapPanTilt["-space"])
					Position.PanTilt.X = interfaceToFloat64(mapPanTilt["-x"])
					Position.PanTilt.Y = interfaceToFloat64(mapPanTilt["-y"])
				}
				// parse Zoom
				if mapZoom, ok := mapPosition["Zoom"].(map[string]interface{}); ok {
					Position.Zoom.Space = interfaceToString(mapZoom["-space"])
					Position.Zoom.X = interfaceToFloat64(mapZoom["-x"])
				}
				Preset.PTZPosition = Position
			}
			// push into result
			result = append(result, Preset)
		}
	}

	return result, nil
}

func (device Device) GotoPreset(profileToken string, presetToken string) error {
	// create soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body: `<GotoPreset xmlns="http://www.onvif.org/ver20/ptz/wsdl">
					<ProfileToken>` + profileToken + `</ProfileToken>
					<PresetToken>` + presetToken + `</PresetToken>
				</GotoPreset>`,
	}

	// send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return err
	}

	_, err = response.ValueForPath("Envelope.Body.GotoPresetResponse")
	if err != nil {
		return err
	}

	return nil
}

func (device Device) RemovePreset(profileToken string, presetToken string) error {
	// create soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body: `<RemovePreset xmlns="http://www.onvif.org/ver20/ptz/wsdl">
					<ProfileToken>` + profileToken + `</ProfileToken>
					<PresetToken>` + presetToken + `</PresetToken>
				</RemovePreset>`,
	}

	// send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return err
	}

	_, err = response.ValueForPath("Envelope.Body.RemovePresetResponse")
	if err != nil {
		return err
	}

	return nil
}
