package onvif

import (
	"fmt"
	"log"
	"testing"
)

func TestGetNodes(t *testing.T)  {
	log.Println("Test GetNodes")

	res, err := testDevice.GetNodes()
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}

func TestGetNode(t *testing.T)  {
	log.Println("Test GetNode")

	res, err := testDevice.GetNode("onvif_ptz_0")
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}

func TestGetConfigurations(t *testing.T)  {
	log.Println("Test GetConfigurations")

	res, err := testDevice.GetNode("onvif_ptz_0")
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}

func TestGetConfiguration(t *testing.T)  {
	log.Println("Test GetConfiguration")

	res, err := testDevice.GetConfiguration("onvif_ptz_0")
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}


func TestGetConfigurationOptions(t *testing.T)  {
	log.Println("Test GetConfigurationOptions")

	res, err := testDevice.GetConfigurationOptions("onvif_ptz_0")
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}

func TestGetStatus(t *testing.T)  {
	log.Println("Test GetStatus")

	res, err := testDevice.GetStatus("mainStream_Profile_Token")
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}

func TestContinuousMove(t *testing.T)  {
	log.Println("Test ContinuousMove")
	velocity := PTZVector{
		PanTilt: Vector2D{
			X: 1,
			Y: 1,
		},
		Zoom: Vector1D{
			X: 0,
		},
	}
	err := testDevice.ContinuousMove("mainStream_Profile_Token", velocity)
	if err != nil {
		t.Error(err)
	}
}

func TestAbsoluteMove(t *testing.T)  {
	log.Println("Test AbsoluteMove")
	position := PTZVector{
		PanTilt: Vector2D{
			X: 1,
			Y: 1,
		},
		Zoom: Vector1D{
			X: 0,
		},
	}
	err := testDevice.AbsoluteMove("mainStream_Profile_Token", position)
	if err != nil {
		t.Error(err)
	}
}


func TestRelativeMove(t *testing.T)  {
	log.Println("Test RelativeMove")
	translation := PTZVector{
		PanTilt: Vector2D{
			X: 1,
			Y: 1,
		},
		Zoom: Vector1D{
			X: 0,
		},
	}
	err := testDevice.RelativeMove("mainStream_Profile_Token", translation)
	if err != nil {
		t.Error(err)
	}
}

func TestStop(t *testing.T)  {
	log.Println("Test Stop")

	err := testDevice.Stop("mainStream_Profile_Token")
	if err != nil {
		t.Error(err)
	}
}

func TestGotoHomePosition(t *testing.T)  {
	log.Println("Test GotoHomePosition")

	err := testDevice.GotoHomePosition("mainStream_Profile_Token")
	if err != nil {
		t.Error(err)
	}
}

func TestSetHomePosition(t *testing.T)  {
	log.Println("Test SetHomePosition")

	log.Println("Change Position")
	translation := PTZVector{
		PanTilt: Vector2D{
			X: 1,
			Y: 1,
		},
		Zoom: Vector1D{
			X: 0,
		},
	}

	testDevice.RelativeMove("MediaProfile000", translation)

	log.Println("Set Home Position")
	err := testDevice.SetHomePosition("MediaProfile000")
	if err != nil {
		t.Error(err)
	}
}

func TestSetPreset(t *testing.T)  {
	log.Println("Test SetPreset")

	res, err := testDevice.SetPreset("MediaProfile000", "preset1")
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}

func TestGetPresets(t *testing.T)  {
	log.Println("Test SetPreset")

	res, err := testDevice.GetPresets("MediaProfile000")
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}

func TestGotoPreset(t *testing.T)  {
	log.Println("Test GotoPreset")

	err := testDevice.GotoPreset("MediaProfile000", "1")
	if err != nil {
		t.Error(err)
	}
}

func TestRemovePreset(t *testing.T)  {
	log.Println("Test RemovePreset")

	err := testDevice.RemovePreset("MediaProfile000", "2")
	if err != nil {
		t.Error(err)
	}
}




