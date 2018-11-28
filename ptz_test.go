package onvif

import (
	"fmt"
	"log"
	"testing"
)

func GetNodes(t *testing.T)  {
	log.Println("Test GetNodes")

	res, err := testDevice.GetNodes()
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}

func GetNode(t *testing.T)  {
	log.Println("Test GetNode")

	res, err := testDevice.GetNode("onvif_ptz_0")
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}

func GetConfigurations(t *testing.T)  {
	log.Println("Test GetConfigurations")

	res, err := testDevice.GetNode("onvif_ptz_0")
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}

func GetConfiguration(t *testing.T)  {
	log.Println("Test GetConfiguration")

	res, err := testDevice.GetConfiguration("onvif_ptz_0")
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}


func GetConfigurationOptions(t *testing.T)  {
	log.Println("Test GetConfigurationOptions")

	res, err := testDevice.GetConfigurationOptions("onvif_ptz_0")
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}

func GetStatus(t *testing.T)  {
	log.Println("Test GetStatus")

	res, err := testDevice.GetStatus("mainStream_Profile_Token")
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}

func ContinuousMove(t *testing.T)  {
	log.Println("Test ContinuousMove")
	velocity := PTZSpeed{
		PanTilt: Vector2D{
			X: 1,
			Y: 1,
		},
		Zoom: Vector1D{
			X: -1,
		},
	}
	err := testDevice.ContinuousMove("mainStream_Profile_Token", velocity)
	if err != nil {
		t.Error(err)
	}
}

func AbsoluteMove(t *testing.T)  {
	log.Println("Test AbsoluteMove")
	position := PTZSpeed{
		PanTilt: Vector2D{
			X: 1,
			Y: 1,
		},
		Zoom: Vector1D{
			X: -1,
		},
	}
	err := testDevice.AbsoluteMove("mainStream_Profile_Token", position)
	if err != nil {
		t.Error(err)
	}
}


func RelativeMove(t *testing.T)  {
	log.Println("Test RelativeMove")
	translation := PTZSpeed{
		PanTilt: Vector2D{
			X: 1,
			Y: 1,
		},
		Zoom: Vector1D{
			X: -1,
		},
	}
	err := testDevice.RelativeMove("mainStream_Profile_Token", translation)
	if err != nil {
		t.Error(err)
	}
}

func Stop(t *testing.T)  {
	log.Println("Test Stop")

	err := testDevice.Stop("mainStream_Profile_Token")
	if err != nil {
		t.Error(err)
	}
}
