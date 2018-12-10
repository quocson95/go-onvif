package onvif

import (
	"fmt"
	"log"
	"testing"
)

func TestGetProfiles(t *testing.T) {
	log.Println("Test GetProfiles")

	res, err := testDevice.GetProfiles()
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}

func TestGetStreamURI(t *testing.T) {
	log.Println("Test GetStreamURI")

	res, err := testDevice.GetStreamURI("IPCProfilesToken0", "UDP")
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}

func TestGetVideoEncoderConfigurations(t *testing.T){
	log.Println("Test GetVideoEncoderConfigurations")

	res, err := testDevice.GetVideoEncoderConfigurations()
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}


func TestSetVideoEncoderConfiguration(t *testing.T){
	log.Println("Test SetVideoEncoderConfiguration")

	profile := MediaProfile{
		VideoEncoderConfig: VideoEncoderConfig{
			Token:"thirdVideoStream_Encoder_Token",
			Name: "thirdVideoStream",
			Encoding:"H264",
			Quality: 95,
			Resolution: MediaBounds{
				Height: 480,
				Width: 720,
			},
			SessionTimeout: "PT0S",
			RateControl: VideoRateControl{
				FrameRateLimit: 20,
				EncodingInterval: 1,
				BitrateLimit: 192,
			},
			H264: H264Configuration{
				GovLength: 45,
				H264Profile: "Main",
			},
			Multicast: Multicast{
				Address: IPAddress{
					Type: "IPv4",
					IPv4Address: "0.0.0.0",
				},
				Port: 8600,
				TTL: 1,
				AutoStart: false,
			},
		},
	}


	err := testDevice.SetVideoEncoderConfiguration(profile)
	if err != nil {
		t.Error(err)
	}
}

func TestGetCompatibleVideoEncoderConfigurations(t *testing.T){
	log.Println("Test GetCompatibleVideoEncoderConfigurations")

	res, err := testDevice.GetCompatibleVideoEncoderConfigurations( "")
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}

func TestGetVideoEncoderConfigurationOptions(t *testing.T){
	log.Println("Test GetVideoEncoderConfigurationOptions")

	res, err := testDevice.GetVideoEncoderConfigurationOptions( "mainVideoStream_Encoder_Token", "")
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}

func TestGetGuaranteedNumberOfVideoEncoderInstances(t *testing.T){
	log.Println("Test GetGuaranteedNumberOfVideoEncoderInstances")

	//res, err := testDevice.GetGuaranteedNumberOfVideoEncoderInstances( "")
	res, err := testDevice.GetGuaranteedNumberOfVideoEncoderInstances( "mainVideoStream_Encoder_Token")
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}

func TestGetProfileMedia(t *testing.T){
	log.Println("Test GetProfileMedia")

	res, err := testDevice.GetProfileMedia( "mainStream_Profile_Token")
	//res, err := testDevice.GetProfileMedia( "fourStream_Profile_Token")
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}

func CreateDeleteProfile(t *testing.T) {
	log.Println("Test CreateDeleteProfile")

	//res, err := testDevice.CreateProfile( "fourStream","fourStream_Profile_Token")
	err := testDevice.DeleteProfile( "fourStream_Profile_Token")
	if err != nil {
		t.Error(err)
	}
}


func TestGetVideoSources(t *testing.T) {
	log.Println("Test GetVideoSources")

	res, err := testDevice.GetVideoSources()
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}

func TestGetVideoSourceConfiguration(t *testing.T) {
	log.Println("Test GetVideoSourceConfiguration")

	res, err := testDevice.GetVideoSourceConfiguration( "VideoStream_Config_Token")
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}

func TestGetVideoSourceConfigurations(t *testing.T) {
	log.Println("Test GetVideoSourceConfigurations")

	res, err := testDevice.GetVideoSourceConfigurations()
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}

func TestGetCompatibleVideoSourceConfigurations(t *testing.T) {
	log.Println("Test GetCompatibleVideoSourceConfigurations")

	res, err := testDevice.GetCompatibleVideoSourceConfigurations( "mainStream_Profile_Token")
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}

func TestGetVideoSourceConfigurationOptions(t *testing.T){
	log.Println("Test GetVideoSourceConfigurationOptions")

	res, err := testDevice.GetVideoSourceConfigurationOptions( "VideoStream_Config_Token", "")
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}

func TestGetMetadataConfiguration(t *testing.T){
	log.Println("Test GetMetadataConfiguration")

	res, err := testDevice.GetMetadataConfiguration( "metadata0")
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}

func TestGetMetadataConfigurations(t *testing.T){
	log.Println("Test GetMetadataConfigurations")

	res, err := testDevice.GetMetadataConfigurations()
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}

func TestGetCompatibleMetadataConfigurations(t *testing.T){
	log.Println("Test GetCompatibleMetadataConfigurations")

	res, err := testDevice.GetCompatibleMetadataConfigurations( "mainStream_Profile_Token")
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}

func TestGetMetadataConfigurationOptions(t *testing.T){
	log.Println("Test GetMetadataConfigurationOptions")

	res, err := testDevice.GetMetadataConfigurationOptions( "metadata0", "")
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}

func TestGetAudioSources(t *testing.T)  {
	log.Println("Test GetAudioSources")

	res, err := testDevice.GetAudioSources()
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}

func TestGetAudioSourceConfiguration(t *testing.T)  {
	log.Println("Test GetAudioSourceConfiguration")

	res, err := testDevice.GetAudioSourceConfiguration("AudioStream_Config_Token")
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}

func TestGetAudioSourceConfigurations(t *testing.T)  {
	log.Println("Test GetAudioSourceConfigurations")

	res, err := testDevice.GetAudioSourceConfigurations()
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}

func TestGetCompatibleAudioSourceConfigurations(t *testing.T)  {
	log.Println("Test GetCompatibleAudioSourceConfigurations")

	res, err := testDevice.GetCompatibleAudioSourceConfigurations("mainStream_Profile_Token")
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}

func TestGetAudioSourceConfigurationOptions(t *testing.T)  {
	log.Println("Test GetAudioSourceConfigurationOptions")

	res, err := testDevice.GetAudioSourceConfigurationOptions("AudioStream_Config_Token", "")
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}


func TestGetAudioEncoderConfiguration(t *testing.T)  {
	log.Println("Test GetAudioEncoderConfiguration")

	res, err := testDevice.GetAudioEncoderConfiguration("AudioStream_Encoder_Token")
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}


func TestGetAudioEncoderConfigurations(t *testing.T)  {
	log.Println("Test GetAudioEncoderConfigurations")

	res, err := testDevice.GetAudioEncoderConfigurations()
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}

func TestGetCompatibleAudioEncoderConfigurations(t *testing.T)  {
	log.Println("Test GetCompatibleAudioEncoderConfigurations")

	res, err := testDevice.GetCompatibleAudioEncoderConfigurations("")
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}

func TestGetAudioEncoderConfigurationOptions(t *testing.T)  {
	log.Println("Test GetAudioEncoderConfigurationOptions")

	res, err := testDevice.GetAudioEncoderConfigurationOptions("AudioStream_Encoder_Token", "")
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}

