package onvif

import (
	"fmt"
	"log"
	"testing"
)

func TestGetAllStreamURI(t *testing.T) {
	log.Println("Test GetStreamURI")
	ovfDevice := Device{
		XAddr:    "http://192.168.0.102/onvif/device_service",
		User:     "admin",
		Password: "Trait17!",
	}
	caps, err := ovfDevice.GetCapabilities()

	if err != nil {
		fmt.Printf("Get cap error %s\n", err.Error())
		return
	}

	ovfMedia := Device{
		XAddr:    caps.Media.XAddr,
		User:     ovfDevice.User,
		Password: ovfDevice.Password,
	}

	profiles, err := ovfMedia.GetProfiles()
	for _, profile := range profiles {
		fmt.Printf("Profile %v\n", profile)
		uri, err := ovfMedia.GetStreamURI(profile.Token, "RTSP")
		fmt.Println("Stream: ", uri, err)
	}
}
