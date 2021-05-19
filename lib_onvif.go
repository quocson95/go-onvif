package onvif

import (
	"encoding/json"
	"github.com/golang/glog"
	"net"
	"time"
)

var (
	mapProfile    = make(map[string]string) // key: device XAddr, value: profile of device
	mapPtzXAddr   = make(map[string]string)
	mapMediaXAddr = make(map[string]string)
)

type OnvifData struct {
	Error string
	Data  interface{}
}

func DiscoveryDevice(interfaceName string, duration int) string {
	result := OnvifData{}
	itf, err := net.InterfaceByName(interfaceName) //here your interface

	if err != nil {
		result.Error = err.Error()
		str, _ := json.Marshal(result)
		return string(str)
	}

	item, _ := itf.Addrs()
	var ip net.IP
	for _, addr := range item {
		switch v := addr.(type) {
		case *net.IPNet:
			if !v.IP.IsLoopback() {
				if v.IP.To4() != nil { //Verify if IP is IPV4
					ip = v.IP
				}
			}
		}
	}

	if ip == nil {
		result.Error = "cannot get device ip"
		str, _ := json.Marshal(result)
		return string(str)
	}

	// Discover device on interface's network
	result.Error = ""
	devices, err := discoverDevices(ip.String(), time.Duration(duration)*time.Millisecond)
	if err != nil {
		result.Error = err.Error()
	}
	result.Data = devices
	str, _ := json.Marshal(result)
	return string(str)
}

// DiscoveryDevice send a WS-Discovery message and wait for all matching device to respond
func GetMediaInformation(host, username, password string) string {
	result := OnvifData{}
	profile := CameraProfile{}

	// Get media address
	od := Device{
		XAddr:    host,
		User:     username,
		Password: password,
	}

	sys, err := od.GetInformation()
	if err == nil {
		profile.Manufacturer = sys.Manufacturer
		profile.Model = sys.Model
		profile.Serial = sys.SerialNumber
		profile.HardwareID = sys.HardwareID
	}

	caps, err := od.GetCapabilities()
	if err != nil {
		profile.LastError = "profile.onvif.getcapabilities.error"
		glog.Warning("Get capabilities error")

		result.Error = "res.error.getcapabilities"
		result.Data = profile
		str, _ := json.Marshal(result)
		return string(str)
	}
	// Get media profile
	odm := Device{
		XAddr:    caps.Media.XAddr,
		User:     username,
		Password: password,
	}

	profiles, err := odm.GetProfiles()
	// Find highest resolution snapshot
	highestHeight := 0
	for _, ovfprofile := range profiles {
		if ovfprofile.VideoEncoderConfig.Resolution.Height > highestHeight {
			highestHeight = ovfprofile.VideoEncoderConfig.Resolution.Height
		}
	}

	if err != nil {
		profile.LastError = "profile.onvif.getprofiles.error"
		glog.Warning("Get profiles error")
		result.Error = "res.error.getprofiles"
		result.Data = profile
		str, _ := json.Marshal(result)
		return string(str)
	}

	for _, ovfprofile := range profiles {
		glog.Infof("Get profile %s", ovfprofile.Token)
		// Get streaming uri
		uri, err := odm.GetStreamURI(ovfprofile.Token, "RTSP")

		if err != nil {
			glog.Warning("Get streaming error")
			continue
		}

		// Only get snapshot on highest resolution
		var snapshotUri string
		if ovfprofile.VideoEncoderConfig.Resolution.Height >= highestHeight {
			// Get streaming uri
			snapshotUri, err = odm.GetSnapshot(ovfprofile.Token)
		}

		if err != nil {
			profile.LastError = "profile.onvif.getsnapshot.error"
			glog.Warningf("Get snapshot for %s error %v", ovfprofile.Token, err)
		}

		profile.Streams = append(profile.Streams, Stream{
			ProfileToken: ovfprofile.Token,
			StreamURI:    uri.URI,
			Resolution: Resolution{
				Width:  ovfprofile.VideoEncoderConfig.Resolution.Width,
				Height: ovfprofile.VideoEncoderConfig.Resolution.Height,
			},
			SnapshotURI:      snapshotUri,
			VideoEncToken:    ovfprofile.VideoEncoderConfig.Token,
			VideoCodec:       ovfprofile.VideoEncoderConfig.Encoding,
			VideoSourceToken: ovfprofile.VideoSourceConfig.Token,
		})

		glog.Infof("Get profile %s done", ovfprofile.Token)
	}

	profile.Authorize = true
	if len(profile.Streams) == 0 {
		profile.LastError = "profile.onvif.getcapabilities.error"
		result.Error = "res.error.getstream"
		data, _ := json.Marshal(profile)
		result.Data = string(data)
		str, _ := json.Marshal(result)
		return string(str)
	}
	result.Error = ""
	result.Data = profile
	str, _ := json.Marshal(result)
	return string(str)
}

func GetXAddress(od Device) (OnvifXAddress, error) {
	result := OnvifXAddress{}
	caps, err := od.GetCapabilities()
	if err != nil || caps.Media.XAddr == "" {
		return result, err
	}
	mapPtzXAddr[od.XAddr] = caps.Ptz.XAddr
	mapMediaXAddr[od.XAddr] = caps.Media.XAddr

	result.PtzXAddress = caps.Ptz.XAddr
	result.MediaXAddress = caps.Media.XAddr
	result.EventXAddress = caps.EventsCap.XAddr
	return result, nil
}

func PtzStart(host, username, password string, x, y, z float64) string {
	result := OnvifData{}
	od := Device{
		XAddr:    host,
		User:     username,
		Password: password,
	}
	// get ptz XAddr and media XAddr
	ptzXAddr := mapPtzXAddr[od.XAddr]
	mediaXAddr := mapMediaXAddr[od.XAddr]
	if ptzXAddr == "" || mediaXAddr == "" {
		glog.Info("Find PTZ And Media Address")
		caps, err := GetXAddress(od)
		if err != nil {
			if CheckAuthorizedError(err.Error()) {
				result.Error = "res.error.unauthorized"
			} else {
				result.Error = "res.error.getptzxaddr"
			}
			str, _ := json.Marshal(result)
			return string(str)
		}
		ptzXAddr = caps.PtzXAddress
		mediaXAddr = caps.MediaXAddress
	}
	// get profile
	profileToken := mapProfile[od.XAddr]
	if profileToken == "" {
		glog.Info("Find Profile")
		// Media device control
		odMedia := Device{
			XAddr:    mediaXAddr,
			User:     username,
			Password: password,
		}
		profiles, err := odMedia.GetProfiles()
		if err != nil {
			glog.Info(err)
			if CheckAuthorizedError(err.Error()) {
				result.Error = "res.error.unauthorized"
			} else {
				result.Error = "res.error.getprofile"
			}
			str, _ := json.Marshal(result)
			return string(str)
		}
		mapProfile[od.XAddr] = profiles[0].Token
		profileToken = mapProfile[od.XAddr]
	}
	glog.Info("PTZ XAddr: ", ptzXAddr)
	glog.Info("Profile Token: ", profileToken)
	// PTZ device control
	odPtz := Device{
		XAddr:    ptzXAddr,
		User:     username,
		Password: password,
	}
	err := odPtz.ContinuousMove(profileToken, PTZVector{
		PanTilt: Vector2D{
			X: x,
			Y: y,
		},
		Zoom: Vector1D{
			X: z,
		},
	})
	if err != nil {
		glog.Warning("PTZ Start Error: ", err.Error())
		if CheckAuthorizedError(err.Error()) {
			result.Error = "res.error.unauthorized"
		} else {
			result.Error = "res.error.ptzstart"
		}
		str, _ := json.Marshal(result)
		return string(str)
	}
	result.Error = ""
	str, _ := json.Marshal(result)
	return string(str)
}

func PtzStop(host, username, password string) string {
	result := OnvifData{}
	od := Device{
		XAddr:    host,
		User:     username,
		Password: password,
	}

	// get ptz XAddr and media XAddr
	ptzXAddr := mapPtzXAddr[od.XAddr]
	mediaXAddr := mapMediaXAddr[od.XAddr]
	if ptzXAddr == "" || mediaXAddr == "" {
		glog.Info("Find PTZ And Media Address")

		caps, err := GetXAddress(od)
		if err != nil {
			if CheckAuthorizedError(err.Error()) {
				result.Error = "res.error.unauthorized"
			} else {
				result.Error = "res.error.getptzxaddr"
			}
			str, _ := json.Marshal(result)
			return string(str)
		}
		ptzXAddr = caps.PtzXAddress
		mediaXAddr = caps.MediaXAddress
	}

	// get profile
	profileToken := mapProfile[od.XAddr]
	if profileToken == "" {
		glog.Info("Find Profile")
		// Media device control
		odMedia := Device{
			XAddr:    mediaXAddr,
			User:     username,
			Password: password,
		}
		profiles, err := odMedia.GetProfiles()
		if err != nil {
			glog.Info(err)
			if CheckAuthorizedError(err.Error()) {
				result.Error = "res.error.unauthorized"
			} else {
				result.Error = "res.error.getprofile"
			}
			str, _ := json.Marshal(result)
			return string(str)
		}
		mapProfile[od.XAddr] = profiles[0].Token
		profileToken = mapProfile[od.XAddr]
	}
	glog.Info("PTZ XAddr: ", ptzXAddr)
	glog.Info("Profile Token: ", profileToken)
	// PTZ device control
	odPtz := Device{
		XAddr:    ptzXAddr,
		User:     username,
		Password: password,
	}
	err := odPtz.Stop(profileToken)
	if err != nil {
		glog.Warning("PTZ Stop Error: ", err.Error())
		if CheckAuthorizedError(err.Error()) {
			result.Error = "res.error.unauthorized"
		} else {
			result.Error = "res.error.ptzstop"
		}
		str, _ := json.Marshal(result)
		return string(str)
	}
	result.Error = ""
	str, _ := json.Marshal(result)
	return string(str)
}
