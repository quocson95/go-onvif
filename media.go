package onvif

import (
	"github.com/golang/glog"
	"strconv"
)

var mediaXMLNs = []string{
	`xmlns:trt="http://www.onvif.org/ver10/media/wsdl"`,
	`xmlns:tt="http://www.onvif.org/ver10/schema"`,
}

// GetProfiles fetch available media profiles of ONVIF camera
func (device Device) GetProfiles() ([]MediaProfile, error) {
	// Create SOAP
	soap := SOAP{
		Body:     "<trt:GetProfiles/>",
		XMLNs:    mediaXMLNs,
		User:     device.User,
		Password: device.Password,
	}

	// Send SOAP request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return []MediaProfile{}, err
	}
	// Get and parse list of profile to interface
	ifaceProfiles, err := response.ValuesForPath("Envelope.Body.GetProfilesResponse.Profiles")
	if err != nil {
		return []MediaProfile{}, err
	}

	// Create initial result
	result := []MediaProfile{}

	// Parse each available profile
	for _, ifaceProfile := range ifaceProfiles {
		if mapProfile, ok := ifaceProfile.(map[string]interface{}); ok {
			// Parse name and token
			profile := MediaProfile{}
			profile.Name = interfaceToString(mapProfile["Name"])
			profile.Token = interfaceToString(mapProfile["-token"])

			// Parse video source configuration
			videoSource := MediaSourceConfig{}
			if mapVideoSource, ok := mapProfile["VideoSourceConfiguration"].(map[string]interface{}); ok {
				videoSource.Name = interfaceToString(mapVideoSource["Name"])
				videoSource.Token = interfaceToString(mapVideoSource["-token"])
				videoSource.SourceToken = interfaceToString(mapVideoSource["SourceToken"])

				// Parse video bounds
				bounds := MediaBounds{}
				if mapVideoBounds, ok := mapVideoSource["Bounds"].(map[string]interface{}); ok {
					bounds.Height = interfaceToInt(mapVideoBounds["-height"])
					bounds.Width = interfaceToInt(mapVideoBounds["-width"])
				}
				videoSource.Bounds = bounds
			}
			profile.VideoSourceConfig = videoSource

			// Parse video encoder configuration
			videoEncoder := VideoEncoderConfig{}
			if mapVideoEncoder, ok := mapProfile["VideoEncoderConfiguration"].(map[string]interface{}); ok {
				videoEncoder.Name = interfaceToString(mapVideoEncoder["Name"])
				videoEncoder.Token = interfaceToString(mapVideoEncoder["-token"])
				videoEncoder.Encoding = interfaceToString(mapVideoEncoder["-encoding"])
				if videoEncoder.Encoding == "" {
					videoEncoder.Encoding = interfaceToString(mapVideoEncoder["Encoding"])
				}
				videoEncoder.Quality = interfaceToFloat64(mapVideoEncoder["Quality"])
				videoEncoder.SessionTimeout = interfaceToString(mapVideoEncoder["SessionTimeout"])

				// Parse video rate control
				rateControl := VideoRateControl{}
				if mapVideoRate, ok := mapVideoEncoder["RateControl"].(map[string]interface{}); ok {
					rateControl.BitrateLimit = interfaceToInt(mapVideoRate["BitrateLimit"])
					rateControl.EncodingInterval = interfaceToInt(mapVideoRate["EncodingInterval"])
					rateControl.FrameRateLimit = interfaceToInt(mapVideoRate["FrameRateLimit"])
				}
				videoEncoder.RateControl = rateControl

				// Parse video resolution
				resolution := MediaBounds{}
				if mapVideoRes, ok := mapVideoEncoder["Resolution"].(map[string]interface{}); ok {
					resolution.Height = interfaceToInt(mapVideoRes["Height"])
					resolution.Width = interfaceToInt(mapVideoRes["Width"])
				}
				videoEncoder.Resolution = resolution
			}
			profile.VideoEncoderConfig = videoEncoder

			// Parse audio source configuration
			audioSource := MediaSourceConfig{}
			if mapAudioSource, ok := mapProfile["AudioSourceConfiguration"].(map[string]interface{}); ok {
				audioSource.Name = interfaceToString(mapAudioSource["Name"])
				audioSource.Token = interfaceToString(mapAudioSource["-token"])
				audioSource.SourceToken = interfaceToString(mapAudioSource["SourceToken"])
			}
			profile.AudioSourceConfig = audioSource

			// Parse audio encoder configuration
			audioEncoder := AudioEncoderConfig{}
			if mapAudioEncoder, ok := mapProfile["AudioEncoderConfiguration"].(map[string]interface{}); ok {
				audioEncoder.Name = interfaceToString(mapAudioEncoder["Name"])
				audioEncoder.Token = interfaceToString(mapAudioEncoder["-token"])
				audioEncoder.Encoding = interfaceToString(mapAudioEncoder["Encoding"])
				audioEncoder.Bitrate = interfaceToInt(mapAudioEncoder["Bitrate"])
				audioEncoder.SampleRate = interfaceToInt(mapAudioEncoder["SampleRate"])
				audioEncoder.SessionTimeout = interfaceToString(mapAudioEncoder["SessionTimeout"])
			}
			profile.AudioEncoderConfig = audioEncoder

			// Parse PTZ configuration
			ptzConfig := PTZConfig{}
			if mapPTZ, ok := mapProfile["PTZConfiguration"].(map[string]interface{}); ok {
				ptzConfig.Name = interfaceToString(mapPTZ["Name"])
				ptzConfig.Token = interfaceToString(mapPTZ["-token"])
				ptzConfig.NodeToken = interfaceToString(mapPTZ["NodeToken"])
			}
			profile.PTZConfig = ptzConfig

			// Push profile to result
			result = append(result, profile)
		}
	}

	return result, nil
}

// GetStreamURI fetch stream URI of a media profile.
// Possible protocol is UDP, HTTP or RTSP
func (device Device) GetStreamURI(profileToken, protocol string) (MediaURI, error) {
	// Create SOAP
	soap := SOAP{
		XMLNs: mediaXMLNs,
		Body: `<trt:GetStreamUri>
			<trt:StreamSetup>
				<tt:Stream>RTP-Unicast</tt:Stream>
				<tt:Transport><tt:Protocol>` + protocol + `</tt:Protocol></tt:Transport>
			</trt:StreamSetup>
			<trt:ProfileToken>` + profileToken + `</trt:ProfileToken>
		</trt:GetStreamUri>`,
		User:     device.User,
		Password: device.Password,
	}

	// Send SOAP request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return MediaURI{}, err
	}

	// Parse response to interface
	ifaceURI, err := response.ValueForPath("Envelope.Body.GetStreamUriResponse.MediaUri")
	if err != nil {
		return MediaURI{}, err
	}

	// Parse interface to struct
	streamURI := MediaURI{}
	if mapURI, ok := ifaceURI.(map[string]interface{}); ok {
		streamURI.URI = interfaceToString(mapURI["Uri"])
		streamURI.Timeout = interfaceToString(mapURI["Timeout"])
		streamURI.InvalidAfterConnect = interfaceToBool(mapURI["InvalidAfterConnect"])
		streamURI.InvalidAfterReboot = interfaceToBool(mapURI["InvalidAfterReboot"])
	}

	return streamURI, nil
}

// GetSnapshot fetch snapshot URI of a media profile.
func (device Device) GetSnapshot(profileToken string) (string, error) {
	soap := SOAP{
		XMLNs:    mediaXMLNs,
		User:     device.User,
		Password: device.Password,
		Body: `<trt:GetSnapshotUri>
				<trt:ProfileToken>` + profileToken + `</trt:ProfileToken>
			 </trt:GetSnapshotUri>`,
	}
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return "", err
	}
	// Parse response to interface
	ifaceURI, err := response.ValueForPath("Envelope.Body.GetSnapshotUriResponse.MediaUri.Uri")
	if err != nil {
		return "", err
	}
	if mediaUri, ok := ifaceURI.(string); ok {
		return mediaUri, nil
	} else {
		return "", err
	}
}

func (device Device) GetVideoEncoderConfigurations() ([]VideoEncoderConfig, error) {
	soap := SOAP{
		Body:     `<GetVideoEncoderConfigurations xmlns="http://www.onvif.org/ver10/media/wsdl"/>`,
		User:     device.User,
		Password: device.Password,
	}
	result := []VideoEncoderConfig{}
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return result, err
	}
	ifaceVideoEncoders, err := response.ValuesForPath("Envelope.Body.GetVideoEncoderConfigurationsResponse.Configurations")
	if err != nil {
		return result, err
	}

	for _, ifaceVideoEncoder := range ifaceVideoEncoders {
		if mapVideoEncoder, ok := ifaceVideoEncoder.(map[string]interface{}); ok {
			videoEncoder := VideoEncoderConfig{}

			videoEncoder.Name = interfaceToString(mapVideoEncoder["Name"])
			videoEncoder.Token = interfaceToString(mapVideoEncoder["-token"])
			videoEncoder.Encoding = interfaceToString(mapVideoEncoder["-encoding"])
			if videoEncoder.Encoding == "" {
				videoEncoder.Encoding = interfaceToString(mapVideoEncoder["Encoding"])
			}
			videoEncoder.Quality = interfaceToFloat64(mapVideoEncoder["Quality"])
			videoEncoder.SessionTimeout = interfaceToString(mapVideoEncoder["SessionTimeout"])
			videoEncoder.GuaranteedFrameRate = interfaceToBool(mapVideoEncoder["GuaranteedFrameRate"])

			// parse Resolution
			if mapResolution, ok := mapVideoEncoder["Resolution"].(map[string]interface{}); ok {
				resolution := MediaBounds{}

				resolution.Width = interfaceToInt(mapResolution["Width"])
				resolution.Height = interfaceToInt(mapResolution["Height"])

				videoEncoder.Resolution = resolution
			}

			// parse Rate Control
			if mapRateControl, ok := mapVideoEncoder["RateControl"].(map[string]interface{}); ok {
				rateControl := VideoRateControl{}

				rateControl.FrameRateLimit = interfaceToInt(mapRateControl["FrameRateLimit"])
				rateControl.EncodingInterval = interfaceToInt(mapRateControl["EncodingInterval"])
				rateControl.BitrateLimit = interfaceToInt(mapRateControl["BitrateLimit"])

				videoEncoder.RateControl = rateControl
			}

			//parse H264
			if mapH264, ok := mapVideoEncoder["H264"].(map[string]interface{}); ok {
				videoEncoder.H264.GovLength = interfaceToInt(mapH264["GovLength"])
				videoEncoder.H264.H264Profile = interfaceToString(mapH264["H264Profile"])
			}

			// parse Multicast
			if mapMulticast, ok := mapVideoEncoder["Multicast"].(map[string]interface{}); ok {
				videoEncoder.Multicast.TTL = interfaceToInt(mapMulticast["TTL"])
				videoEncoder.Multicast.Port = interfaceToInt(mapMulticast["Port"])
				videoEncoder.Multicast.AutoStart = interfaceToBool(mapMulticast["AutoStart"])

				// parse Address
				if mapAddress, ok := mapMulticast["Address"].(map[string]interface{}); ok {
					videoEncoder.Multicast.Address.Type = interfaceToString(mapAddress["Type"])
					videoEncoder.Multicast.Address.IPv4Address = interfaceToString(mapAddress["IPv4Address"])
				}
			}

			// add to result
			result = append(result, videoEncoder)
		}
	}

	glog.Info(result)
	return result, nil
}

func (device Device) SetVideoEncoderConfiguration(videoEncoderConfig VideoEncoderConfig) error {
	soap := SOAP{
		XMLNs:    mediaXMLNs,
		User:     device.User,
		Password: device.Password,
		Body: `<trt:SetVideoEncoderConfiguration xmlns="http://www.onvif.org/ver10/media/wsdl">
					<trt:Configuration token="` + videoEncoderConfig.Token + `">
						<tt:Name>` + videoEncoderConfig.Name + `</tt:Name>
						<tt:Encoding>` + videoEncoderConfig.Encoding + `</tt:Encoding>
						<tt:Quality>` + float64ToString(videoEncoderConfig.Quality) + `</tt:Quality>
						<tt:SessionTimeout>` + videoEncoderConfig.SessionTimeout + `</tt:SessionTimeout>
						<tt:Resolution>
							<tt:Width>` + strconv.Itoa(videoEncoderConfig.Resolution.Width) + `</tt:Width>
							<tt:Height>` + strconv.Itoa(videoEncoderConfig.Resolution.Height) + `</tt:Height>
						</tt:Resolution>
						<tt:RateControl>
							<tt:FrameRateLimit>` + strconv.Itoa(videoEncoderConfig.RateControl.FrameRateLimit) + `</tt:FrameRateLimit>
							<tt:EncodingInterval>` + strconv.Itoa(videoEncoderConfig.RateControl.EncodingInterval) + `</tt:EncodingInterval>
							<tt:BitrateLimit>` + strconv.Itoa(videoEncoderConfig.RateControl.BitrateLimit) + `</tt:BitrateLimit>
						</tt:RateControl>
						<tt:H264>
							<tt:GovLength>` + strconv.Itoa(videoEncoderConfig.H264.GovLength) + `</tt:GovLength>
							<tt:H264Profile>` + videoEncoderConfig.H264.H264Profile + `</tt:H264Profile>
						</tt:H264>
						<tt:Multicast>
							<tt:Address>
								<tt:Type>` + videoEncoderConfig.Multicast.Address.Type + `</tt:Type>
								<tt:IPv4Address>` + videoEncoderConfig.Multicast.Address.IPv4Address + `</tt:IPv4Address>
							</tt:Address>
							<tt:Port>` + strconv.Itoa(videoEncoderConfig.Multicast.Port) + `</tt:Port>
							<tt:TTL>` + strconv.Itoa(videoEncoderConfig.Multicast.TTL) + `</tt:TTL>
							<tt:AutoStart>` + strconv.FormatBool(videoEncoderConfig.Multicast.AutoStart) + `</tt:AutoStart>
						</tt:Multicast>
					</trt:Configuration>
					<trt:ForcePersistence>true</trt:ForcePersistence>
				</trt:SetVideoEncoderConfiguration>`,
	}

	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return err
	}
	_, err = response.ValueForPath("Envelope.Body.SetVideoEncoderConfigurationResponse")
	if err != nil {
		return err
	}
	return nil
}

func (device Device) GetCompatibleVideoEncoderConfigurations(profileToken string) ([]VideoEncoderConfig, error) {
	soap := SOAP{
		Body: `<GetCompatibleVideoEncoderConfigurations xmlns="http://www.onvif.org/ver10/media/wsdl">
					<ProfileToken xmlns="http://www.onvif.org/ver10/schema">` + profileToken + `</ProfileToken></GetCompatibleVideoEncoderConfigurations>`,
		User:     device.User,
		Password: device.Password,
	}
	result := []VideoEncoderConfig{}
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return result, err
	}
	ifaceVideoEncoders, err := response.ValuesForPath("Envelope.Body.GetCompatibleVideoEncoderConfigurationsResponse.Configurations")
	if err != nil {
		return result, err
	}

	for _, ifaceVideoEncoder := range ifaceVideoEncoders {
		if mapVideoEncoder, ok := ifaceVideoEncoder.(map[string]interface{}); ok {
			videoEncoder := VideoEncoderConfig{}

			videoEncoder.Name = interfaceToString(mapVideoEncoder["Name"])
			videoEncoder.Token = interfaceToString(mapVideoEncoder["-token"])
			videoEncoder.Encoding = interfaceToString(mapVideoEncoder["Encoding"])
			videoEncoder.Quality = interfaceToFloat64(mapVideoEncoder["Quality"])
			videoEncoder.SessionTimeout = interfaceToString(mapVideoEncoder["SessionTimeout"])

			// parse Resolution
			if mapResolution, ok := mapVideoEncoder["Resolution"].(map[string]interface{}); ok {
				resolution := MediaBounds{}

				resolution.Width = interfaceToInt(mapResolution["Width"])
				resolution.Height = interfaceToInt(mapResolution["Height"])

				videoEncoder.Resolution = resolution
			}

			// parse Rate Control
			if mapRateControl, ok := mapVideoEncoder["RateControl"].(map[string]interface{}); ok {
				rateControl := VideoRateControl{}

				rateControl.FrameRateLimit = interfaceToInt(mapRateControl["FrameRateLimit"])
				rateControl.EncodingInterval = interfaceToInt(mapRateControl["EncodingInterval"])
				rateControl.BitrateLimit = interfaceToInt(mapRateControl["BitrateLimit"])

				videoEncoder.RateControl = rateControl
			}

			// add to result
			result = append(result, videoEncoder)
		}
	}

	return result, nil
}

// truyen vao mot trong 2 tham so
func (device Device) GetVideoEncoderConfigurationOptions(configurationToken string, profileToken string) (VideoEncoderConfigurationOptions, error) {
	// create token body
	tokenBody := ``
	if configurationToken != "" {
		tokenBody = `<ConfigurationToken>` + configurationToken + `</ConfigurationToken>`
	} else {
		tokenBody = `<ProfileToken>` + profileToken + `</ProfileToken>`
	}

	// create soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body:     `<GetVideoEncoderConfigurationOptions xmlns="http://www.onvif.org/ver10/media/wsdl">` + tokenBody + `</GetVideoEncoderConfigurationOptions>`,
	}

	result := VideoEncoderConfigurationOptions{}

	// send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return result, err
	}

	// parse response into interface
	ifaceVideoEncoderConfOption, err := response.ValueForPath("Envelope.Body.GetVideoEncoderConfigurationOptionsResponse.Options")
	if err != nil {
		return result, err
	}

	// parse interface
	if mapOptions, ok := ifaceVideoEncoderConfOption.(map[string]interface{}); ok {
		// parse Quality Range
		if mapQualityRange, ok := mapOptions["QualityRange"].(map[string]interface{}); ok {
			qualityRange := IntRange{}

			qualityRange.Min = interfaceToInt(mapQualityRange["Min"])
			qualityRange.Max = interfaceToInt(mapQualityRange["Max"])

			result.QualityRange = qualityRange
		}

		// parse H264
		if mapH264, ok := mapOptions["H264"].(map[string]interface{}); ok {
			h264Options := H264Options{}

			// parse Resolution Available
			if mapResolution, ok := mapH264["ResolutionsAvailable"].(map[string]interface{}); ok {
				h264Options.ResolutionsAvailable.Height = interfaceToInt(mapResolution["Height"])
				h264Options.ResolutionsAvailable.Width = interfaceToInt(mapResolution["Width"])
			}
			// parse GovLengthRange
			if mapGovLengthRange, ok := mapH264["GovLengthRange"].(map[string]interface{}); ok {
				h264Options.GovLengthRange.Min = interfaceToInt(mapGovLengthRange["Min"])
				h264Options.GovLengthRange.Max = interfaceToInt(mapGovLengthRange["Max"])
			}
			// parse Frame Rate Range
			if mapFrameRateRange, ok := mapH264["FrameRateRange"].(map[string]interface{}); ok {
				h264Options.FrameRateRange.Min = interfaceToInt(mapFrameRateRange["Min"])
				h264Options.FrameRateRange.Max = interfaceToInt(mapFrameRateRange["Max"])
			}
			// parse Encoding Interval Range
			if mapEncodingIntervalRange, ok := mapH264["EncodingIntervalRange"].(map[string]interface{}); ok {
				h264Options.EncodingIntervalRange.Min = interfaceToInt(mapEncodingIntervalRange["Min"])
				h264Options.EncodingIntervalRange.Max = interfaceToInt(mapEncodingIntervalRange["Max"])
			}
			// parse H264 Profiles Supported
			if H264ProfilesSupported, ok := mapH264["H264ProfilesSupported"].([]interface{}); ok {
				h264Supported := []string{}
				for _, H264ProfileSupported := range H264ProfilesSupported {
					h264Supported = append(h264Supported, interfaceToString(H264ProfileSupported))
				}
				h264Options.H264ProfilesSupported = h264Supported
			}

			result.H264 = h264Options
		}
	}

	return result, nil
}

func (device Device) GetGuaranteedNumberOfVideoEncoderInstances(configurationToken string) (GuaranteedNumberOfVideoEncoderInstances, error) {
	// create soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body: `<GetGuaranteedNumberOfVideoEncoderInstances xmlns="http://www.onvif.org/ver10/media/wsdl">
					<ConfigurationToken>` + configurationToken + `</ConfigurationToken>
				</GetGuaranteedNumberOfVideoEncoderInstances>`,
	}

	result := GuaranteedNumberOfVideoEncoderInstances{}

	//send reuest
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return result, err
	}

	// parse response into interface
	ifaceVideoEncoderInstances, err := response.ValueForPath("Envelope.Body.GetGuaranteedNumberOfVideoEncoderInstancesResponse")
	if err != nil {
		return result, err
	}

	// parse interface
	if mapVideoEncoderInstances, ok := ifaceVideoEncoderInstances.(map[string]interface{}); ok {
		result.TotalNumber = interfaceToInt(mapVideoEncoderInstances["TotalNumber"])
		result.H264 = interfaceToInt(mapVideoEncoderInstances["H264"])
	}

	return result, err
}

func (device Device) GetProfileMedia(profileToken string) (MediaProfile, error) {
	// Create SOAP
	soap := SOAP{
		Body: `<GetProfile xmlns="http://www.onvif.org/ver10/media/wsdl">
						<ProfileToken>` + profileToken + `</ProfileToken>
					</GetProfile>`,
		User:     device.User,
		Password: device.Password,
	}

	result := MediaProfile{}
	// Send SOAP request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return result, err
	}

	// Get and parse list of profile to interface
	ifaceProfile, err := response.ValueForPath("Envelope.Body.GetProfileResponse.Profile")
	if err != nil {
		return result, err
	}

	// Parse available profile
	if mapProfile, ok := ifaceProfile.(map[string]interface{}); ok {
		// Parse name and token
		result.Name = interfaceToString(mapProfile["Name"])
		result.Token = interfaceToString(mapProfile["-token"])

		// Parse video source configuration
		videoSource := MediaSourceConfig{}
		if mapVideoSource, ok := mapProfile["VideoSourceConfiguration"].(map[string]interface{}); ok {
			videoSource.Name = interfaceToString(mapVideoSource["Name"])
			videoSource.Token = interfaceToString(mapVideoSource["-token"])
			videoSource.SourceToken = interfaceToString(mapVideoSource["SourceToken"])

			// Parse video bounds
			bounds := MediaBounds{}
			if mapVideoBounds, ok := mapVideoSource["Bounds"].(map[string]interface{}); ok {
				bounds.Height = interfaceToInt(mapVideoBounds["-height"])
				bounds.Width = interfaceToInt(mapVideoBounds["-width"])
			}
			videoSource.Bounds = bounds
		}
		result.VideoSourceConfig = videoSource

		// Parse video encoder configuration
		videoEncoder := VideoEncoderConfig{}
		if mapVideoEncoder, ok := mapProfile["VideoEncoderConfiguration"].(map[string]interface{}); ok {
			videoEncoder.Name = interfaceToString(mapVideoEncoder["Name"])
			videoEncoder.Token = interfaceToString(mapVideoEncoder["-token"])
			videoEncoder.Encoding = interfaceToString(mapVideoEncoder["-encoding"])
			if videoEncoder.Encoding == "" {
				videoEncoder.Encoding = interfaceToString(mapVideoEncoder["Encoding"])
			}
			videoEncoder.Quality = interfaceToFloat64(mapVideoEncoder["Quality"])
			videoEncoder.SessionTimeout = interfaceToString(mapVideoEncoder["SessionTimeout"])

			// Parse video rate control
			rateControl := VideoRateControl{}
			if mapVideoRate, ok := mapVideoEncoder["RateControl"].(map[string]interface{}); ok {
				rateControl.BitrateLimit = interfaceToInt(mapVideoRate["BitrateLimit"])
				rateControl.EncodingInterval = interfaceToInt(mapVideoRate["EncodingInterval"])
				rateControl.FrameRateLimit = interfaceToInt(mapVideoRate["FrameRateLimit"])
			}
			videoEncoder.RateControl = rateControl

			// Parse video resolution
			resolution := MediaBounds{}
			if mapVideoRes, ok := mapVideoEncoder["Resolution"].(map[string]interface{}); ok {
				resolution.Height = interfaceToInt(mapVideoRes["Height"])
				resolution.Width = interfaceToInt(mapVideoRes["Width"])
			}
			videoEncoder.Resolution = resolution
		}
		result.VideoEncoderConfig = videoEncoder

		// Parse audio source configuration
		audioSource := MediaSourceConfig{}
		if mapAudioSource, ok := mapProfile["AudioSourceConfiguration"].(map[string]interface{}); ok {
			audioSource.Name = interfaceToString(mapAudioSource["Name"])
			audioSource.Token = interfaceToString(mapAudioSource["-token"])
			audioSource.SourceToken = interfaceToString(mapAudioSource["SourceToken"])
		}
		result.AudioSourceConfig = audioSource

		// Parse audio encoder configuration
		audioEncoder := AudioEncoderConfig{}
		if mapAudioEncoder, ok := mapProfile["AudioEncoderConfiguration"].(map[string]interface{}); ok {
			audioEncoder.Name = interfaceToString(mapAudioEncoder["Name"])
			audioEncoder.Token = interfaceToString(mapAudioEncoder["-token"])
			audioEncoder.Encoding = interfaceToString(mapAudioEncoder["Encoding"])
			audioEncoder.Bitrate = interfaceToInt(mapAudioEncoder["Bitrate"])
			audioEncoder.SampleRate = interfaceToInt(mapAudioEncoder["SampleRate"])
			audioEncoder.SessionTimeout = interfaceToString(mapAudioEncoder["SessionTimeout"])
		}
		result.AudioEncoderConfig = audioEncoder

		// Parse PTZ configuration
		ptzConfig := PTZConfig{}
		if mapPTZ, ok := mapProfile["PTZConfiguration"].(map[string]interface{}); ok {
			ptzConfig.Name = interfaceToString(mapPTZ["Name"])
			ptzConfig.Token = interfaceToString(mapPTZ["-token"])
			ptzConfig.NodeToken = interfaceToString(mapPTZ["NodeToken"])
		}
		result.PTZConfig = ptzConfig
	}

	return result, nil
}

func (device Device) CreateProfile(profileName string, profileToken string) (MediaProfile, error) {
	// create soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body: `<CreateProfile xmlns="http://www.onvif.org/ver10/media/wsdl">
					<Name>` + profileName + `</Name>
					<Token>` + profileToken + `</Token>
				</CreateProfile>`,
	}

	result := MediaProfile{}
	// Send SOAP request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return result, err
	}

	// Get and parse list of profile to interface
	ifaceProfile, err := response.ValueForPath("Envelope.Body.CreateProfileResponse.Profile")
	if err != nil {
		return result, err
	}

	// Parse profile
	if mapProfile, ok := ifaceProfile.(map[string]interface{}); ok {
		// Parse name and token
		result.Name = interfaceToString(mapProfile["Name"])
		result.Token = interfaceToString(mapProfile["-token"])

		// Parse video source configuration
		videoSource := MediaSourceConfig{}
		if mapVideoSource, ok := mapProfile["VideoSourceConfiguration"].(map[string]interface{}); ok {
			videoSource.Name = interfaceToString(mapVideoSource["Name"])
			videoSource.Token = interfaceToString(mapVideoSource["-token"])
			videoSource.SourceToken = interfaceToString(mapVideoSource["SourceToken"])

			// Parse video bounds
			bounds := MediaBounds{}
			if mapVideoBounds, ok := mapVideoSource["Bounds"].(map[string]interface{}); ok {
				bounds.Height = interfaceToInt(mapVideoBounds["-height"])
				bounds.Width = interfaceToInt(mapVideoBounds["-width"])
			}
			videoSource.Bounds = bounds
		}
		result.VideoSourceConfig = videoSource

		// Parse video encoder configuration
		videoEncoder := VideoEncoderConfig{}
		if mapVideoEncoder, ok := mapProfile["VideoEncoderConfiguration"].(map[string]interface{}); ok {
			videoEncoder.Name = interfaceToString(mapVideoEncoder["Name"])
			videoEncoder.Token = interfaceToString(mapVideoEncoder["-token"])
			videoEncoder.Encoding = interfaceToString(mapVideoEncoder["Encoding"])
			videoEncoder.Quality = interfaceToFloat64(mapVideoEncoder["Quality"])
			videoEncoder.SessionTimeout = interfaceToString(mapVideoEncoder["SessionTimeout"])

			// Parse video rate control
			rateControl := VideoRateControl{}
			if mapVideoRate, ok := mapVideoEncoder["RateControl"].(map[string]interface{}); ok {
				rateControl.BitrateLimit = interfaceToInt(mapVideoRate["BitrateLimit"])
				rateControl.EncodingInterval = interfaceToInt(mapVideoRate["EncodingInterval"])
				rateControl.FrameRateLimit = interfaceToInt(mapVideoRate["FrameRateLimit"])
			}
			videoEncoder.RateControl = rateControl

			// Parse video resolution
			resolution := MediaBounds{}
			if mapVideoRes, ok := mapVideoEncoder["Resolution"].(map[string]interface{}); ok {
				resolution.Height = interfaceToInt(mapVideoRes["Height"])
				resolution.Width = interfaceToInt(mapVideoRes["Width"])
			}
			videoEncoder.Resolution = resolution
		}
		result.VideoEncoderConfig = videoEncoder

		// Parse audio source configuration
		audioSource := MediaSourceConfig{}
		if mapAudioSource, ok := mapProfile["AudioSourceConfiguration"].(map[string]interface{}); ok {
			audioSource.Name = interfaceToString(mapAudioSource["Name"])
			audioSource.Token = interfaceToString(mapAudioSource["-token"])
			audioSource.SourceToken = interfaceToString(mapAudioSource["SourceToken"])
		}
		result.AudioSourceConfig = audioSource

		// Parse audio encoder configuration
		audioEncoder := AudioEncoderConfig{}
		if mapAudioEncoder, ok := mapProfile["AudioEncoderConfiguration"].(map[string]interface{}); ok {
			audioEncoder.Name = interfaceToString(mapAudioEncoder["Name"])
			audioEncoder.Token = interfaceToString(mapAudioEncoder["-token"])
			audioEncoder.Encoding = interfaceToString(mapAudioEncoder["Encoding"])
			audioEncoder.Bitrate = interfaceToInt(mapAudioEncoder["Bitrate"])
			audioEncoder.SampleRate = interfaceToInt(mapAudioEncoder["SampleRate"])
			audioEncoder.SessionTimeout = interfaceToString(mapAudioEncoder["SessionTimeout"])
		}
		result.AudioEncoderConfig = audioEncoder

		// Parse PTZ configuration
		ptzConfig := PTZConfig{}
		if mapPTZ, ok := mapProfile["PTZConfiguration"].(map[string]interface{}); ok {
			ptzConfig.Name = interfaceToString(mapPTZ["Name"])
			ptzConfig.Token = interfaceToString(mapPTZ["-token"])
			ptzConfig.NodeToken = interfaceToString(mapPTZ["NodeToken"])
		}
		result.PTZConfig = ptzConfig
	}

	return result, nil

}

func (device Device) DeleteProfile(profileToken string) error {
	// create soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body: `<DeleteProfile xmlns="http://www.onvif.org/ver10/media/wsdl">
					<ProfileToken>` + profileToken + `</ProfileToken>
				</DeleteProfile>`,
	}

	// send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return err
	}

	_, err = response.ValueForPath("Envelope.Body.DeleteProfileResponse")
	if err != nil {
		return err
	}

	return nil
}

func (device Device) GetVideoSources() ([]VideoSource, error) {
	//create soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body:     `<GetVideoSources xmlns="http://www.onvif.org/ver10/media/wsdl"/>`,
	}

	result := []VideoSource{}

	// send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return result, err
	}

	// parse response into interface
	ifaceVideoSources, err := response.ValuesForPath("Envelope.Body.GetVideoSourcesResponse.VideoSources")
	if err != nil {
		return result, err
	}

	// parse interface
	for _, ifaceVideoSource := range ifaceVideoSources {
		if mapVideoSource, ok := ifaceVideoSource.(map[string]interface{}); ok {
			videoSource := VideoSource{}

			videoSource.Token = interfaceToString(mapVideoSource["-token"])
			videoSource.Framerate = interfaceToFloat64(mapVideoSource["Framerate"])

			// parse resolution
			if mapResolution, ok := mapVideoSource["Resolution"].(map[string]interface{}); ok {
				videoSource.Resolution.Height = interfaceToInt(mapResolution["Height"])
				videoSource.Resolution.Width = interfaceToInt(mapResolution["Width"])
			}
			// parse imaging
			if mapImaging, ok := mapVideoSource["Imaging"].(map[string]interface{}); ok {
				imaging := ImagingSettings{}

				// parse Backlight Compensation
				if mapBacklightCompensation, ok := mapImaging["BacklightCompensation"].(map[string]interface{}); ok {
					imaging.BacklightCompensation.Mode = interfaceToString(mapBacklightCompensation["Mode"])
					imaging.BacklightCompensation.Level = interfaceToFloat64(mapBacklightCompensation["Level"])
				}

				imaging.Brightness = interfaceToFloat64(mapImaging["Brightness"])
				imaging.ColorSaturation = interfaceToFloat64(mapImaging["ColorSaturation"])
				imaging.Contrast = interfaceToFloat64(mapImaging["Contrast"])

				// parse Exposure
				if mapExposure, ok := mapImaging["Exposure"].(map[string]interface{}); ok {
					exposure := Exposure{}

					exposure.Mode = interfaceToString(mapExposure["Mode"])
					exposure.Priority = interfaceToString(mapExposure["Priority"])

					exposure.MinExposureTime = interfaceToFloat64(mapExposure["MinExposureTime"])
					exposure.MaxExposureTime = interfaceToFloat64(mapExposure["MaxExposureTime"])
					exposure.MinGain = interfaceToFloat64(mapExposure["MinGain"])
					exposure.MaxGain = interfaceToFloat64(mapExposure["MaxGain"])
					exposure.MinIris = interfaceToFloat64(mapExposure["MinIris"])
					exposure.MaxIris = interfaceToFloat64(mapExposure["MaxIris"])
					exposure.ExposureTime = interfaceToFloat64(mapExposure["ExposureTime"])
					exposure.Gain = interfaceToFloat64(mapExposure["Gain"])
					exposure.Iris = interfaceToFloat64(mapExposure["Iris"])

					// parse window
					if mapWindow, ok := mapExposure["Window"].(map[string]interface{}); ok {
						exposure.Window.Top = interfaceToInt(mapWindow["-top"])
						exposure.Window.Bottom = interfaceToInt(mapWindow["-bottom"])
						exposure.Window.Left = interfaceToInt(mapWindow["-left"])
						exposure.Window.Right = interfaceToInt(mapWindow["-right"])
					}

					imaging.Exposure = exposure
				}

				// parse focus
				if mapFocus, ok := mapImaging["Focus"].(map[string]interface{}); ok {
					focus := FocusConfiguration{}

					focus.AutoFocusMode = interfaceToString(mapFocus["AutoFocusMode"])
					focus.DefaultSpeed = interfaceToFloat64(mapFocus["DefaultSpeed"])
					focus.FarLimit = interfaceToFloat64(mapFocus["FarLimit"])
					focus.NearLimit = interfaceToFloat64(mapFocus["NearLimit"])

					imaging.Focus = focus
				}

				imaging.IrCutFilter = interfaceToString(mapImaging["IrCutFilter"])
				imaging.Sharpness = interfaceToFloat64(mapImaging["Sharpness"])

				// parse WideDynamicRange
				if mapWideDynamicRange, ok := mapImaging["WideDynamicRange"].(map[string]interface{}); ok {
					imaging.WideDynamicRange.Mode = interfaceToString(mapWideDynamicRange["Mode"])
					imaging.WideDynamicRange.Level = interfaceToFloat64(mapWideDynamicRange["Level"])
				}

				// parse WhiteBalance
				if mapWhiteBalance, ok := mapImaging["WhiteBalance"].(map[string]interface{}); ok {
					whiteBalance := WhiteBalance{}

					whiteBalance.Mode = interfaceToString(mapWhiteBalance["Mode"])
					whiteBalance.CbGain = interfaceToFloat64(mapWhiteBalance["CbGain"])
					whiteBalance.CrGain = interfaceToFloat64(mapWhiteBalance["CrGain"])

					imaging.WhiteBalance = whiteBalance
				}

				videoSource.Imaging = imaging
			}

			// push to result
			result = append(result, videoSource)
		}
	}

	return result, nil
}

func (device Device) GetVideoSourceConfiguration(configurationToken string) (VideoSourceConfiguration, error) {
	// create soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body: `<GetVideoSourceConfiguration xmlns="http://www.onvif.org/ver10/media/wsdl">
					<ConfigurationToken>` + configurationToken + `</ConfigurationToken>
				</GetVideoSourceConfiguration>`,
	}

	result := VideoSourceConfiguration{}

	// send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return result, err
	}

	// parse rsponse into interface
	ifaceVideoSouceConfiguration, err := response.ValueForPath("Envelope.Body.GetVideoSourceConfigurationResponse.Configuration")
	if err != nil {
		return result, err
	}

	// parse interface into struct
	if mapVideoSouceConfiguration, ok := ifaceVideoSouceConfiguration.(map[string]interface{}); ok {
		result.Token = interfaceToString(mapVideoSouceConfiguration["-token"])
		result.Name = interfaceToString(mapVideoSouceConfiguration["Name"])
		result.SourceToken = interfaceToString(mapVideoSouceConfiguration["SourceToken"])
		// parse bounds
		if mapBounds, ok := mapVideoSouceConfiguration["Bounds"].(map[string]interface{}); ok {
			bounds := IntRectangle{}

			bounds.X = interfaceToInt(mapBounds["-x"])
			bounds.Y = interfaceToInt(mapBounds["-y"])
			bounds.Height = interfaceToInt(mapBounds["-height"])
			bounds.Width = interfaceToInt(mapBounds["-width"])

			result.Bounds = bounds
		}
	}

	return result, nil
}

func (device Device) GetVideoSourceConfigurations() ([]VideoSourceConfiguration, error) {
	// create soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body:     `<GetVideoSourceConfigurations xmlns="http://www.onvif.org/ver10/media/wsdl"/>`,
	}

	result := []VideoSourceConfiguration{}

	// send soap request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return result, err
	}

	// parse response
	ifaceConfigurations, err := response.ValuesForPath("Envelope.Body.GetVideoSourceConfigurationsResponse.Configurations")
	if err != nil {
		return result, err
	}

	for _, ifaceConfiguration := range ifaceConfigurations {
		// parse interface into struct
		if mapVideoSouceConfiguration, ok := ifaceConfiguration.(map[string]interface{}); ok {
			videoSourceConfiguration := VideoSourceConfiguration{}

			videoSourceConfiguration.Token = interfaceToString(mapVideoSouceConfiguration["-token"])
			videoSourceConfiguration.Name = interfaceToString(mapVideoSouceConfiguration["Name"])
			videoSourceConfiguration.SourceToken = interfaceToString(mapVideoSouceConfiguration["SourceToken"])
			// parse bounds
			if mapBounds, ok := mapVideoSouceConfiguration["Bounds"].(map[string]interface{}); ok {
				bounds := IntRectangle{}

				bounds.X = interfaceToInt(mapBounds["-x"])
				bounds.Y = interfaceToInt(mapBounds["-y"])
				bounds.Height = interfaceToInt(mapBounds["-height"])
				bounds.Width = interfaceToInt(mapBounds["-width"])

				videoSourceConfiguration.Bounds = bounds
			}

			// push to result
			result = append(result, videoSourceConfiguration)
		}
	}

	return result, nil
}

func (device Device) GetCompatibleVideoSourceConfigurations(profileToken string) ([]VideoSourceConfiguration, error) {
	//create soap request
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body: `<GetCompatibleVideoSourceConfigurations xmlns="http://www.onvif.org/ver10/media/wsdl">
					<ProfileToken>` + profileToken + `</ProfileToken>
				</GetCompatibleVideoSourceConfigurations>`,
	}
	result := []VideoSourceConfiguration{}

	// send soap request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return result, err
	}

	// parse response
	ifaceConfigurations, err := response.ValuesForPath("Envelope.Body.GetCompatibleVideoSourceConfigurationsResponse.Configurations")
	if err != nil {
		return result, err
	}

	for _, ifaceConfiguration := range ifaceConfigurations {
		// parse interface into struct
		if mapVideoSouceConfiguration, ok := ifaceConfiguration.(map[string]interface{}); ok {
			videoSourceConfiguration := VideoSourceConfiguration{}

			videoSourceConfiguration.Token = interfaceToString(mapVideoSouceConfiguration["-token"])
			videoSourceConfiguration.Name = interfaceToString(mapVideoSouceConfiguration["Name"])
			videoSourceConfiguration.SourceToken = interfaceToString(mapVideoSouceConfiguration["SourceToken"])
			// parse bounds
			if mapBounds, ok := mapVideoSouceConfiguration["Bounds"].(map[string]interface{}); ok {
				bounds := IntRectangle{}

				bounds.X = interfaceToInt(mapBounds["-x"])
				bounds.Y = interfaceToInt(mapBounds["-y"])
				bounds.Height = interfaceToInt(mapBounds["-height"])
				bounds.Width = interfaceToInt(mapBounds["-width"])

				videoSourceConfiguration.Bounds = bounds
			}

			// push to result
			result = append(result, videoSourceConfiguration)
		}
	}

	return result, nil
}

// truyen vao mot trong 2 tham so
func (device Device) GetVideoSourceConfigurationOptions(configurationToken string, profileToken string) (VideoSourceConfigurationOption, error) {
	// create token body
	tokenBody := ``
	if configurationToken != "" {
		tokenBody = `<ConfigurationToken>` + configurationToken + `</ConfigurationToken>`
	} else {
		tokenBody = `<ProfileToken>` + profileToken + `</ProfileToken>`
	}

	//create soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body:     `<GetVideoSourceConfigurationOptions xmlns="http://www.onvif.org/ver10/media/wsdl">` + tokenBody + `</GetVideoSourceConfigurationOptions>`,
	}

	result := VideoSourceConfigurationOption{}
	// send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return result, err
	}

	// parse response
	ifaceConfiguration, err := response.ValueForPath("Envelope.Body.GetVideoSourceConfigurationOptionsResponse.Options")
	if err != nil {
		return result, err
	}

	// parse interface
	if mapConfigurationOption, ok := ifaceConfiguration.(map[string]interface{}); ok {
		result.MaximumNumberOfProfiles = interfaceToInt(mapConfigurationOption["MaximumNumberOfProfiles"])
		// parse BoundsRange
		if mapBoundsRange, ok := mapConfigurationOption["BoundsRange"].(map[string]interface{}); ok {
			boundsRange := IntRectangleRange{}
			// parse XRange
			if mapXRange, ok := mapBoundsRange["XRange"].(map[string]interface{}); ok {
				boundsRange.XRange.Max = interfaceToInt(mapXRange["Max"])
				boundsRange.XRange.Min = interfaceToInt(mapXRange["Min"])
			}
			// parse YRange
			if mapYRange, ok := mapBoundsRange["YRange"].(map[string]interface{}); ok {
				boundsRange.YRange.Max = interfaceToInt(mapYRange["Max"])
				boundsRange.YRange.Min = interfaceToInt(mapYRange["Min"])
			}
			// parse WidthRange
			if mapWidthRange, ok := mapBoundsRange["WidthRange"].(map[string]interface{}); ok {
				boundsRange.WidthRange.Max = interfaceToInt(mapWidthRange["Max"])
				boundsRange.WidthRange.Min = interfaceToInt(mapWidthRange["Min"])
			}
			// parse HeightRange
			if mapHeightRange, ok := mapBoundsRange["HeightRange"].(map[string]interface{}); ok {
				boundsRange.HeightRange.Max = interfaceToInt(mapHeightRange["Max"])
				boundsRange.HeightRange.Min = interfaceToInt(mapHeightRange["Min"])
			}
			result.BoundsRange = boundsRange
		}
		result.VideoSourceTokensAvailable = interfaceToString(mapConfigurationOption["VideoSourceTokensAvailable"])
	}

	return result, nil
}

func (device Device) GetMetadataConfiguration(configurationToken string) (MetadataConfiguration, error) {
	//send soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body: `<GetMetadataConfiguration xmlns="http://www.onvif.org/ver10/media/wsdl">
					<ConfigurationToken>` + configurationToken + `</ConfigurationToken>
				</GetMetadataConfiguration>`,
	}

	result := MetadataConfiguration{}
	// send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return result, err
	}

	// parse response into interface
	ifaceMetadata, err := response.ValueForPath("Envelope.Body.GetMetadataConfigurationResponse.Configuration")
	if err != nil {
		return result, err
	}

	// parse interface
	if mapMetadata, ok := ifaceMetadata.(map[string]interface{}); ok {
		result.Name = interfaceToString(mapMetadata["Name"])
		result.Token = interfaceToString(mapMetadata["-token"])
		result.SessionTimeout = interfaceToString(mapMetadata["SessionTimeout"])
		// parse Multicast
		if mapMulticast, ok := mapMetadata["Multicast"].(map[string]interface{}); ok {
			multicast := Multicast{}
			// parse address
			if mapAddress, ok := mapMulticast["Address"].(map[string]interface{}); ok {
				multicast.Address.Type = interfaceToString(mapAddress["Type"])
				multicast.Address.IPv4Address = interfaceToString(mapAddress["IPv4Address"])
			}

			multicast.AutoStart = interfaceToBool(mapMulticast["AutoStart"])
			multicast.Port = interfaceToInt(mapMulticast["Port"])
			multicast.TTL = interfaceToInt(mapMulticast["TTL"])

			result.Multicast = multicast
		}
	}

	return result, nil
}

func (device Device) GetMetadataConfigurations() ([]MetadataConfiguration, error) {
	// create soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body:     `<GetMetadataConfigurations xmlns="http://www.onvif.org/ver10/media/wsdl"/>`,
	}

	result := []MetadataConfiguration{}

	// send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return result, err
	}

	// parse response
	ifaceMetaConfigurations, err := response.ValuesForPath("Envelope.Body.GetMetadataConfigurationsResponse.Configurations")
	if err != nil {
		return result, err
	}

	// parse interface
	for _, ifaceMetaConfiguration := range ifaceMetaConfigurations {
		if mapMetadataConfiguration, ok := ifaceMetaConfiguration.(map[string]interface{}); ok {
			metadataConfiguration := MetadataConfiguration{}

			metadataConfiguration.Name = interfaceToString(mapMetadataConfiguration["Name"])
			metadataConfiguration.Token = interfaceToString(mapMetadataConfiguration["-token"])
			metadataConfiguration.SessionTimeout = interfaceToString(mapMetadataConfiguration["SessionTimeout"])
			// parse Multicast
			if mapMulticast, ok := mapMetadataConfiguration["Multicast"].(map[string]interface{}); ok {
				multicast := Multicast{}
				// parse address
				if mapAddress, ok := mapMulticast["Address"].(map[string]interface{}); ok {
					multicast.Address.Type = interfaceToString(mapAddress["Type"])
					multicast.Address.IPv4Address = interfaceToString(mapAddress["IPv4Address"])
				}

				multicast.AutoStart = interfaceToBool(mapMulticast["AutoStart"])
				multicast.Port = interfaceToInt(mapMulticast["Port"])
				multicast.TTL = interfaceToInt(mapMulticast["TTL"])

				metadataConfiguration.Multicast = multicast
			}

			// push into result
			result = append(result, metadataConfiguration)
		}
	}

	return result, nil
}

func (device Device) GetCompatibleMetadataConfigurations(profileToken string) ([]MetadataConfiguration, error) {
	// create soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body: `<GetCompatibleMetadataConfigurations xmlns="http://www.onvif.org/ver10/media/wsdl">
					<ProfileToken>` + profileToken + `</ProfileToken>
				</GetCompatibleMetadataConfigurations>`,
	}

	result := []MetadataConfiguration{}

	// send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return result, err
	}

	// parse response
	ifaceMetaConfigurations, err := response.ValuesForPath("Envelope.Body.GetCompatibleMetadataConfigurationsResponse.Configurations")
	if err != nil {
		return result, err
	}

	// parse interface
	for _, ifaceMetaConfiguration := range ifaceMetaConfigurations {
		if mapMetadataConfiguration, ok := ifaceMetaConfiguration.(map[string]interface{}); ok {
			metadataConfiguration := MetadataConfiguration{}

			metadataConfiguration.Name = interfaceToString(mapMetadataConfiguration["Name"])
			metadataConfiguration.Token = interfaceToString(mapMetadataConfiguration["-token"])
			metadataConfiguration.SessionTimeout = interfaceToString(mapMetadataConfiguration["SessionTimeout"])
			// parse Multicast
			if mapMulticast, ok := mapMetadataConfiguration["Multicast"].(map[string]interface{}); ok {
				multicast := Multicast{}
				// parse address
				if mapAddress, ok := mapMulticast["Address"].(map[string]interface{}); ok {
					multicast.Address.Type = interfaceToString(mapAddress["Type"])
					multicast.Address.IPv4Address = interfaceToString(mapAddress["IPv4Address"])
				}

				multicast.AutoStart = interfaceToBool(mapMulticast["AutoStart"])
				multicast.Port = interfaceToInt(mapMulticast["Port"])
				multicast.TTL = interfaceToInt(mapMulticast["TTL"])

				metadataConfiguration.Multicast = multicast
			}

			// push into result
			result = append(result, metadataConfiguration)
		}
	}

	return result, nil
}

// truyen vao mot trong 2 tham so
func (device Device) GetMetadataConfigurationOptions(configurationToken string, profileToken string) (MetadataConfigurationOptions, error) {
	// create token body
	tokenBody := ``
	if configurationToken != "" {
		tokenBody = `<ConfigurationToken>` + configurationToken + `</ConfigurationToken>`
	} else {
		tokenBody = `<ProfileToken>` + profileToken + `</ProfileToken>`
	}

	// create soap request
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body:     `<GetMetadataConfigurationOptions xmlns="http://www.onvif.org/ver10/media/wsdl">` + tokenBody + `</GetMetadataConfigurationOptions>`,
	}
	result := MetadataConfigurationOptions{}
	// send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return result, err
	}

	// parse response
	ifaceMetadataConfigurationOptions, err := response.ValueForPath("Envelope.Body.GetMetadataConfigurationOptionsResponse.Options")
	if err != nil {
		return result, err
	}
	// parse interface
	if mapMetadataConfigurationOptions, ok := ifaceMetadataConfigurationOptions.(map[string]interface{}); ok {
		result.GeoLocation = interfaceToBool(mapMetadataConfigurationOptions["GeoLocation"])
		// parse PTZStatusFilterOptions
		if mapPTZStatusFilterOptions, ok := mapMetadataConfigurationOptions["PTZStatusFilterOptions"].(map[string]interface{}); ok {
			PTZStatusFilterOptions := PTZStatusFilterOptions{}

			PTZStatusFilterOptions.PanTiltPositionSupported = interfaceToBool(mapPTZStatusFilterOptions["PanTiltPositionSupported"])
			PTZStatusFilterOptions.PanTiltStatusSupported = interfaceToBool(mapPTZStatusFilterOptions["PanTiltStatusSupported"])
			PTZStatusFilterOptions.ZoomPositionSupported = interfaceToBool(mapPTZStatusFilterOptions["ZoomPositionSupported"])
			PTZStatusFilterOptions.ZoomStatusSupported = interfaceToBool(mapPTZStatusFilterOptions["ZoomStatusSupported"])

			result.PTZStatusFilterOptions = PTZStatusFilterOptions
		}
	}
	return result, nil
}

func (device Device) GetAudioSources() ([]AudioSource, error) {
	// create soap request
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body:     `<GetAudioSources xmlns="http://www.onvif.org/ver10/media/wsdl"/>`,
	}

	result := []AudioSource{}

	//send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return result, err
	}

	// parse response
	ifaceAudioSources, err := response.ValuesForPath("Envelope.Body.GetAudioSourcesResponse.AudioSources")
	if err != nil {
		return result, err
	}

	// parse interface
	for _, ifaceAudioSource := range ifaceAudioSources {
		if mapAudioSource, ok := ifaceAudioSource.(map[string]interface{}); ok {
			audioSource := AudioSource{}

			audioSource.Token = interfaceToString(mapAudioSource["-token"])
			audioSource.Channels = interfaceToInt(mapAudioSource["Channels"])

			// push into result
			result = append(result, audioSource)
		}
	}

	return result, nil
}

func (device Device) GetAudioSourceConfiguration(configurationToken string) (AudioSourceConfiguration, error) {
	// create soap request
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body: `<GetAudioSourceConfiguration xmlns="http://www.onvif.org/ver10/media/wsdl">
					<ConfigurationToken>` + configurationToken + `</ConfigurationToken>
				</GetAudioSourceConfiguration>`,
	}

	result := AudioSourceConfiguration{}

	//send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return result, err
	}

	// parse response
	ifaceAudioSourceConfiguration, err := response.ValueForPath("Envelope.Body.GetAudioSourceConfigurationResponse.Configuration")
	if err != nil {
		return result, err
	}

	// parse interface
	if mapAudioSourceConfiguration, ok := ifaceAudioSourceConfiguration.(map[string]interface{}); ok {
		result.Token = interfaceToString(mapAudioSourceConfiguration["-token"])
		result.Name = interfaceToString(mapAudioSourceConfiguration["Name"])
		result.SourceToken = interfaceToString(mapAudioSourceConfiguration["SourceToken"])
	}

	return result, nil
}

func (device Device) GetAudioSourceConfigurations() ([]AudioSourceConfiguration, error) {
	// create soap request
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body:     `<GetAudioSourceConfigurations xmlns="http://www.onvif.org/ver10/media/wsdl"/>`,
	}

	result := []AudioSourceConfiguration{}

	//send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return result, err
	}

	// parse response
	ifaceAudioSourceConfigurations, err := response.ValuesForPath("Envelope.Body.GetAudioSourceConfigurationsResponse.Configurations")
	if err != nil {
		return result, err
	}

	// parse interface
	for _, ifaceAudioSourceConfiguration := range ifaceAudioSourceConfigurations {
		if mapAudioSourceConfiguration, ok := ifaceAudioSourceConfiguration.(map[string]interface{}); ok {
			audioSourceConfiguration := AudioSourceConfiguration{}

			audioSourceConfiguration.Token = interfaceToString(mapAudioSourceConfiguration["-token"])
			audioSourceConfiguration.Name = interfaceToString(mapAudioSourceConfiguration["Name"])
			audioSourceConfiguration.SourceToken = interfaceToString(mapAudioSourceConfiguration["SourceToken"])

			// push into result
			result = append(result, audioSourceConfiguration)
		}
	}

	glog.Info(result)
	return result, nil
}

func (device Device) GetCompatibleAudioSourceConfigurations(profileToken string) ([]AudioSourceConfiguration, error) {
	// create soap request
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body: `<GetCompatibleAudioSourceConfigurations xmlns="http://www.onvif.org/ver10/media/wsdl">
					<ProfileToken>` + profileToken + `</ProfileToken>
				</GetCompatibleAudioSourceConfigurations>`,
	}

	result := []AudioSourceConfiguration{}

	//send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return result, err
	}

	// parse response
	ifaceAudioSourceConfigurations, err := response.ValuesForPath("Envelope.Body.GetCompatibleAudioSourceConfigurationsResponse.Configurations")
	if err != nil {
		return result, err
	}

	// parse interface
	for _, ifaceAudioSourceConfiguration := range ifaceAudioSourceConfigurations {
		if mapAudioSourceConfiguration, ok := ifaceAudioSourceConfiguration.(map[string]interface{}); ok {
			audioSourceConfiguration := AudioSourceConfiguration{}

			audioSourceConfiguration.Token = interfaceToString(mapAudioSourceConfiguration["-token"])
			audioSourceConfiguration.Name = interfaceToString(mapAudioSourceConfiguration["Name"])
			audioSourceConfiguration.SourceToken = interfaceToString(mapAudioSourceConfiguration["SourceToken"])

			// push into result
			result = append(result, audioSourceConfiguration)
		}
	}

	return result, nil
}

// fetch input tokens available
// truyen vao mot trong 2 tham so
func (device Device) GetAudioSourceConfigurationOptions(configurationToken string, profileToken string) (string, error) {
	// create token body
	tokenBody := ``
	if configurationToken != "" {
		tokenBody = `<ConfigurationToken>` + configurationToken + `</ConfigurationToken>`
	} else {
		tokenBody = `<ProfileToken>` + profileToken + `</ProfileToken>`
	}

	// create soap request
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body:     `<GetAudioSourceConfigurationOptions xmlns="http://www.onvif.org/ver10/media/wsdl">` + tokenBody + `</GetAudioSourceConfigurationOptions>`,
	}

	var result string
	//send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return result, err
	}

	// parse response
	iface, err := response.ValueForPath("Envelope.Body.GetAudioSourceConfigurationOptionsResponse.Options")
	if err != nil {
		return result, err
	}

	// parse interface
	if mapAudioSourceConfigurationOption, ok := iface.(map[string]interface{}); ok {
		result = interfaceToString(mapAudioSourceConfigurationOption["InputTokensAvailable"])
	}

	glog.Info(result)
	return result, nil
}

func (device Device) GetAudioEncoderConfiguration(configurationToken string) (AudioEncoderConfig, error) {
	// create soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body: `<GetAudioEncoderConfiguration xmlns="http://www.onvif.org/ver10/media/wsdl">
					<ConfigurationToken>` + configurationToken + `</ConfigurationToken>
				</GetAudioEncoderConfiguration>`,
	}

	result := AudioEncoderConfig{}
	// send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return result, err
	}

	// parse response
	ifaceAudioEncoder, err := response.ValueForPath("Envelope.Body.GetAudioEncoderConfigurationResponse.Configuration")
	if err != nil {
		return result, err
	}

	// parse interface
	if mapAudioEncoder, ok := ifaceAudioEncoder.(map[string]interface{}); ok {
		result.Token = interfaceToString(mapAudioEncoder["-token"])
		result.Name = interfaceToString(mapAudioEncoder["Name"])
		result.Bitrate = interfaceToInt(mapAudioEncoder["Bitrate"])
		result.Encoding = interfaceToString(mapAudioEncoder["Encoding"])
		result.SampleRate = interfaceToInt(mapAudioEncoder["SampleRate"])
		result.SessionTimeout = interfaceToString(mapAudioEncoder["SessionTimeout"])
	}

	return result, nil
}

func (device Device) GetAudioEncoderConfigurations() ([]AudioEncoderConfig, error) {
	// create soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body:     `<GetAudioEncoderConfigurations xmlns="http://www.onvif.org/ver10/media/wsdl"/>`,
	}

	result := []AudioEncoderConfig{}
	// send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return result, err
	}

	// parse response
	ifaceAudioEncoderConfigurations, err := response.ValuesForPath("Envelope.Body.GetAudioEncoderConfigurationsResponse.Configurations")
	if err != nil {
		return result, err
	}

	// parse interface
	for _, ifaceAudioEncoderConfiguration := range ifaceAudioEncoderConfigurations {
		if mapAudioEncoderConf, ok := ifaceAudioEncoderConfiguration.(map[string]interface{}); ok {
			audioEncoderConfig := AudioEncoderConfig{}

			audioEncoderConfig.Token = interfaceToString(mapAudioEncoderConf["-token"])
			audioEncoderConfig.Name = interfaceToString(mapAudioEncoderConf["Name"])
			audioEncoderConfig.Bitrate = interfaceToInt(mapAudioEncoderConf["Bitrate"])
			audioEncoderConfig.Encoding = interfaceToString(mapAudioEncoderConf["Encoding"])
			audioEncoderConfig.SampleRate = interfaceToInt(mapAudioEncoderConf["SampleRate"])
			audioEncoderConfig.SessionTimeout = interfaceToString(mapAudioEncoderConf["SessionTimeout"])

			// push into result
			result = append(result, audioEncoderConfig)
		}
	}

	return result, nil
}

func (device Device) GetCompatibleAudioEncoderConfigurations(profileToken string) ([]AudioEncoderConfig, error) {
	// create soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body: `<GetCompatibleAudioEncoderConfigurations xmlns="http://www.onvif.org/ver10/media/wsdl">
					<ProfileToken>` + profileToken + `</ProfileToken>
				</GetCompatibleAudioEncoderConfigurations>`,
	}

	result := []AudioEncoderConfig{}
	// send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return result, err
	}

	// parse response
	ifaceAudioEncoderConfigurations, err := response.ValuesForPath("Envelope.Body.GetCompatibleAudioEncoderConfigurationsResponse.Configurations")
	if err != nil {
		return result, err
	}

	// parse interface
	for _, ifaceAudioEncoderConfiguration := range ifaceAudioEncoderConfigurations {
		if mapAudioEncoderConf, ok := ifaceAudioEncoderConfiguration.(map[string]interface{}); ok {
			audioEncoderConfig := AudioEncoderConfig{}

			audioEncoderConfig.Token = interfaceToString(mapAudioEncoderConf["-token"])
			audioEncoderConfig.Name = interfaceToString(mapAudioEncoderConf["Name"])
			audioEncoderConfig.Bitrate = interfaceToInt(mapAudioEncoderConf["Bitrate"])
			audioEncoderConfig.Encoding = interfaceToString(mapAudioEncoderConf["Encoding"])
			audioEncoderConfig.SampleRate = interfaceToInt(mapAudioEncoderConf["SampleRate"])
			audioEncoderConfig.SessionTimeout = interfaceToString(mapAudioEncoderConf["SessionTimeout"])

			// push into result
			result = append(result, audioEncoderConfig)
		}
	}
	return result, nil
}

// truyen vao mot trong 2 tham so
func (device Device) GetAudioEncoderConfigurationOptions(configurationToken string, profileToken string) ([]AudioEncoderConfigurationOption, error) {
	// create token body
	tokenBody := ``
	if configurationToken != "" {
		tokenBody = `<ConfigurationToken>` + configurationToken + `</ConfigurationToken>`
	} else {
		tokenBody = `<ProfileToken>` + profileToken + `</ProfileToken>`
	}

	// create soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		Body:     `<GetAudioEncoderConfigurationOptions xmlns="http://www.onvif.org/ver10/media/wsdl">` + tokenBody + `</GetAudioEncoderConfigurationOptions>`,
	}

	result := []AudioEncoderConfigurationOption{}

	// send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return result, err
	}

	// parse response
	ifaceAudioEncoderConfigurationOptions, err := response.ValuesForPath("Envelope.Body.GetAudioEncoderConfigurationOptionsResponse.Options.Options")
	if err != nil {
		return result, err
	}

	// parse interface
	for _, ifaceAudioEncoderConfigurationOption := range ifaceAudioEncoderConfigurationOptions {
		if mapAudioEncoderConfigurationOption, ok := ifaceAudioEncoderConfigurationOption.(map[string]interface{}); ok {
			audioEncoderConfigOption := AudioEncoderConfigurationOption{}

			audioEncoderConfigOption.Encoding = interfaceToString(mapAudioEncoderConfigurationOption["Encoding"])
			audioEncoderConfigOption.BitrateList = interfaceToInt(mapAudioEncoderConfigurationOption["BitrateList"])
			audioEncoderConfigOption.SampleRateList = interfaceToInt(mapAudioEncoderConfigurationOption["SampleRateList"])

			result = append(result, audioEncoderConfigOption)
		}
	}

	glog.Info(result)
	return result, nil
}

func (device Device) GetMasks(configurationToken string) ([]Mask, error) {
	// create soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		XMLNs: []string{
			`xmlns:tr2="http://www.onvif.org/ver20/media/wsdl"`,
		},
		Body: `<GetMasks xmlns="http://www.onvif.org/ver20/media/wsdl">
						<ConfigurationToken>` + configurationToken + `</ConfigurationToken>
					</GetMasks>`,
	}

	result := make([]Mask, 0)

	// send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return result, err
	}

	// parse response
	ifMasksResponse, err := response.ValuesForPath("Envelope.Body.GetMasksResponse.Masks")
	if err != nil {
		return result, err
	}

	for _, mask := range ifMasksResponse {
		if mapMask, ok := mask.(map[string]interface{}); ok {

			r := Mask{}

			r.Token = interfaceToString(mapMask["-token"])
			r.ConfigurationToken = interfaceToString(mapMask["ConfigurationToken"])
			r.Type = interfaceToString(mapMask["Type"])
			r.Enabled = interfaceToBool(mapMask["Enabled"])

			if mapPolygon, ok := mapMask["Polygon"].(map[string]interface{}); ok {
				if mapPoints, ok := mapPolygon["Point"].([]interface{}); ok {
					for _, point := range mapPoints {
						if mapPoint, ok := point.(map[string]interface{}); ok {
							p := Point{}
							p.X = interfaceToInt(mapPoint["-x"])
							p.Y = interfaceToInt(mapPoint["-y"])
							r.Polygon = append(r.Polygon, p)
						}
					}
				}
			}

			result = append(result, r)
		}
	}

	return result, nil
}

func (device Device) CreateMask(configurationToken string, pointStart, pointEnd Point) (string, error) {
	// create soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		XMLNs: []string{
			`xmlns:tr2="http://www.onvif.org/ver20/media/wsdl"`,
			`xmlns:tt="http://www.onvif.org/ver10/schema"`,
		},
		Action: "http://www.onvif.org/ver20/media/wsdl/CreateMask",
		Body: `<tr2:CreateMask xmlns="http://www.onvif.org/ver20/media/wsdl">
						<tr2:Mask>
							<tr2:ConfigurationToken>` + configurationToken + `</tr2:ConfigurationToken>
							<tr2:Polygon>
								<Point xmlns="http://www.onvif.org/ver10/schema" y="` + intToString(pointStart.Y) + `" x="` + intToString(pointStart.X) + `"></Point>
								<Point xmlns="http://www.onvif.org/ver10/schema" y="` + intToString(pointEnd.Y) + `" x="` + intToString(pointEnd.X) + `"></Point>
								<Point xmlns="http://www.onvif.org/ver10/schema" y="` + intToString(pointStart.Y) + `" x="` + intToString(pointStart.X) + `"></Point>
								<Point xmlns="http://www.onvif.org/ver10/schema" y="` + intToString(pointEnd.Y) + `" x="` + intToString(pointEnd.X) + `"></Point>
							</tr2:Polygon>
							<tr2:Type>Blurred</tr2:Type>
							<tr2:Enabled>true</tr2:Enabled>
						</tr2:Mask>
					</tr2:CreateMask>`,
	}

	// send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return "", err
	}

	// parse response
	ifOSDsResponse, err := response.ValueForPath("Envelope.Body.CreateMaskResponse.Token")
	if err != nil {
		return "", err
	}

	glog.Info("Data %v", interfaceToString(ifOSDsResponse))
	return interfaceToString(ifOSDsResponse), nil
}

func (device Device) UpdateMask(maskToken, configurationToken string, pointStart, pointEnd Point, enable bool) error {
	// create soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		XMLNs: []string{
			`xmlns:tr2="http://www.onvif.org/ver20/media/wsdl"`,
			`xmlns:tt="http://www.onvif.org/ver10/schema"`,
		},
		Action: "http://www.onvif.org/ver20/media/wsdl/SetMask",
		Body: `<tr2:SetMask xmlns="http://www.onvif.org/ver20/media/wsdl">
						<tr2:Mask token="` + maskToken + `">
							<tr2:ConfigurationToken>` + configurationToken + `</tr2:ConfigurationToken>
							<tr2:Polygon>
								<Point xmlns="http://www.onvif.org/ver10/schema" y="` + intToString(pointStart.Y) + `" x="` + intToString(pointStart.X) + `"></Point>
								<Point xmlns="http://www.onvif.org/ver10/schema" y="` + intToString(pointEnd.Y) + `" x="` + intToString(pointEnd.X) + `"></Point>
								<Point xmlns="http://www.onvif.org/ver10/schema" y="` + intToString(pointStart.Y) + `" x="` + intToString(pointStart.X) + `"></Point>
								<Point xmlns="http://www.onvif.org/ver10/schema" y="` + intToString(pointEnd.Y) + `" x="` + intToString(pointEnd.X) + `"></Point>
							</tr2:Polygon>
							<tr2:Type>Blurred</tr2:Type>
							<tr2:Enabled>` + boolToString(enable) + `</tr2:Enabled>
						</tr2:Mask>
					</tr2:SetMask>`,
	}
	// send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return err
	}

	// parse response
	_, err = response.ValueForPath("Envelope.Body.SetMaskResponse")
	if err != nil {
		return err
	}

	return nil
}

func (device Device) DeleteMask(maskToken string) error {
	// create soap
	soap := SOAP{
		User:     device.User,
		Password: device.Password,
		XMLNs: []string{
			`xmlns:tr2="http://www.onvif.org/ver20/media/wsdl"`,
			`xmlns:tt="http://www.onvif.org/ver10/schema"`,
		},
		Action: "http://www.onvif.org/ver20/media/wsdl/DeleteMask",
		Body: `<tr2:DeleteMask xmlns="http://www.onvif.org/ver20/media/wsdl">
					<Token>` + maskToken + `</Token>
			   </tr2:DeleteMask>`,
	}
	// send request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return err
	}

	// parse response
	_, err = response.ValueForPath("Envelope.Body.DeleteMaskResponse")
	if err != nil {
		return err
	}

	return nil
}
