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



func GetVideoEncoderConfigurations(t *testing.T){
	log.Println("Test GetVideoEncoderConfigurations")

	res, err := testDevice.GetVideoEncoderConfigurations()
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}

func GetCompatibleVideoEncoderConfigurations(t *testing.T){
	log.Println("Test GetCompatibleVideoEncoderConfigurations")

	res, err := testDevice.GetCompatibleVideoEncoderConfigurations( "")
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}

func GetVideoEncoderConfigurationOptions(t *testing.T){
	log.Println("Test GetVideoEncoderConfigurationOptions")

	res, err := testDevice.GetVideoEncoderConfigurationOptions( "mainVideoStream_Encoder_Token", "")
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}

func GetGuaranteedNumberOfVideoEncoderInstances(t *testing.T){
	log.Println("Test GetGuaranteedNumberOfVideoEncoderInstances")

	//res, err := testDevice.GetGuaranteedNumberOfVideoEncoderInstances( "")
	res, err := testDevice.GetGuaranteedNumberOfVideoEncoderInstances( "mainVideoStream_Encoder_Token")
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}

func GetProfileMedia(t *testing.T){
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


func GetVideoSources(t *testing.T) {
	log.Println("Test GetVideoSources")

	res, err := testDevice.GetVideoSources()
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}

func GetVideoSourceConfiguration(t *testing.T) {
	log.Println("Test GetVideoSourceConfiguration")

	res, err := testDevice.GetVideoSourceConfiguration( "VideoStream_Config_Token")
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}

func GetVideoSourceConfigurations(t *testing.T) {
	log.Println("Test GetVideoSourceConfigurations")

	res, err := testDevice.GetVideoSourceConfigurations()
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}

func GetCompatibleVideoSourceConfigurations(t *testing.T) {
	log.Println("Test GetCompatibleVideoSourceConfigurations")

	res, err := testDevice.GetCompatibleVideoSourceConfigurations( "mainStream_Profile_Token")
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}

func GetVideoSourceConfigurationOptions(t *testing.T){
	log.Println("Test GetVideoSourceConfigurationOptions")

	res, err := testDevice.GetVideoSourceConfigurationOptions( "VideoStream_Config_Token", "")
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}

func GetMetadataConfiguration(t *testing.T){
	log.Println("Test GetMetadataConfiguration")

	res, err := testDevice.GetMetadataConfiguration( "metadata0")
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}

func GetMetadataConfigurations(t *testing.T){
	log.Println("Test GetMetadataConfigurations")

	res, err := testDevice.GetMetadataConfigurations()
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}

func GetCompatibleMetadataConfigurations(t *testing.T){
	log.Println("Test GetCompatibleMetadataConfigurations")

	res, err := testDevice.GetCompatibleMetadataConfigurations( "mainStream_Profile_Token")
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}

func GetMetadataConfigurationOptions(t *testing.T){
	log.Println("Test GetMetadataConfigurationOptions")

	res, err := testDevice.GetMetadataConfigurationOptions( "metadata0", "")
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}

func GetAudioSources(t *testing.T)  {
	log.Println("Test GetAudioSources")

	res, err := testDevice.GetAudioSources()
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}

func GetAudioSourceConfiguration(t *testing.T)  {
	log.Println("Test GetAudioSourceConfiguration")

	res, err := testDevice.GetAudioSourceConfiguration("AudioStream_Config_Token")
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}

func GetAudioSourceConfigurations(t *testing.T)  {
	log.Println("Test GetAudioSourceConfigurations")

	res, err := testDevice.GetAudioSourceConfigurations()
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}

func GetCompatibleAudioSourceConfigurations(t *testing.T)  {
	log.Println("Test GetCompatibleAudioSourceConfigurations")

	res, err := testDevice.GetCompatibleAudioSourceConfigurations("mainStream_Profile_Token")
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}

func GetAudioSourceConfigurationOptions(t *testing.T)  {
	log.Println("Test GetAudioSourceConfigurationOptions")

	res, err := testDevice.GetAudioSourceConfigurationOptions("AudioStream_Config_Token", "")
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}


func GetAudioEncoderConfiguration(t *testing.T)  {
	log.Println("Test GetAudioEncoderConfiguration")

	res, err := testDevice.GetAudioEncoderConfiguration("AudioStream_Encoder_Token")
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}


func GetAudioEncoderConfigurations(t *testing.T)  {
	log.Println("Test GetAudioEncoderConfigurations")

	res, err := testDevice.GetAudioEncoderConfigurations()
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}

func GetCompatibleAudioEncoderConfigurations(t *testing.T)  {
	log.Println("Test GetCompatibleAudioEncoderConfigurations")

	res, err := testDevice.GetCompatibleAudioEncoderConfigurations("")
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}

func GetAudioEncoderConfigurationOptions(t *testing.T)  {
	log.Println("Test GetAudioEncoderConfigurationOptions")

	res, err := testDevice.GetAudioEncoderConfigurationOptions("AudioStream_Encoder_Token", "")
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}

