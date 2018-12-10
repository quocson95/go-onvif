package onvif

import (
	"fmt"
	"log"
	"testing"
)

func TestGetInformation(t *testing.T) {
	log.Println("Test GetInformation")

	res, err := testDevice.GetInformation()
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}

func TestGetCapabilities(t *testing.T) {
	log.Println("Test GetCapabilities")

	res, err := testDevice.GetCapabilities()
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}

func TestGetDiscoveryMode(t *testing.T) {
	log.Println("Test GetDiscoveryMode")

	res, err := testDevice.GetDiscoveryMode()
	if err != nil {
		t.Error(err)
	}

	fmt.Println(res)
}

func TestGetScopes(t *testing.T) {
	log.Println("Test GetScopes")

	res, err := testDevice.GetScopes()
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}

func TestGetHostname(t *testing.T) {
	log.Println("Test GetHostname")

	res, err := testDevice.GetHostname()
	if err != nil {
		t.Error(err)
	}

	js := prettyJSON(&res)
	fmt.Println(js)
}

/// Todo Test Get, Set SystemDateTime mode Manual
func GetSetSystemDateAndTime(t *testing.T){
	log.Println("Test GetSetSystemDateAndTime")

	utcDT := SystemDateAndTime{
		DateTimeType: "Manual",
		DaylightSavings:false,
		TimeZone: TimeZone{
			TZ:"CST-8",
		},
		UTCDateTime: UTCDateTime{
			Time: Time{
				Hour: 10,
				Minute: 20,
				Second: 10,
			},
			Date: Date{
				Year: 2018,
				Month: 11,
				Day: 15,
			},
		},
	}

	res, err := testDevice.GetSystemDateAndTime()
	if err != nil {
		t.Error(err)
	}
	js := prettyJSON(&res)
	fmt.Println(js)

	err = testDevice.SetSystemDateAndTime(utcDT)
	if err != nil {
		t.Error(err)
	}
}


/// Todo Test Set Get NTP, SetSytemDateAndTime mode NTP
// Todo Config NTP
func  TestGetSetNTP(t *testing.T)  {
	log.Println("Test GetSetNTP")

	ntpInformation := NTPInformation{
		FromDHCP: false,
		NTPNetworkHost: NetworkHost{
				Type: "IPv4",
				IPv4Address: "swisstime.ee.ethz.ch",
		},
	}
	res, err := testDevice.GetNTP()
	if err != nil {
		t.Error(err)
	}
	js := prettyJSON(&res)
	fmt.Println(js)

	err = testDevice.SetNTP(ntpInformation)
	if err != nil {
		t.Error(err)
	}
}


//// Todo Config SystemDateAndTime with mode NTP
func  GetSetSystemTimeNTP(t *testing.T)  {
	log.Println("Test GetSetSystemTimeNTP")

	utcDT := SystemDateAndTime{
		DateTimeType: "NTP",
		DaylightSavings:false,
	}
	res, err := testDevice.GetSystemDateAndTime()
	if err != nil {
		t.Error(err)
	}
	js := prettyJSON(&res)
	fmt.Println(js)

	err = testDevice.SetSystemDateAndTime(utcDT)
	if err != nil {
		t.Error(err)
	}
}

//// Todo Reboot Device
func TestSystemReboot(t *testing.T)  {
	log.Println("Test SystemReboot")

	res, err := testDevice.SystemReboot()
	if err != nil {
		t.Error(err)
	}
	js := prettyJSON(&res)
	fmt.Println(js)
}


/// Todo Get Set device
func GetSetDNS(t *testing.T) {
	log.Println("Test GetSetDNS")

	res, err := testDevice.GetDNS()
	if err != nil {
		t.Error(err)
	}
	js := prettyJSON(&res)
	fmt.Println(js)

	dnsInformation :=DNSInformation{
		FromDHCP: false,
		SearchDomain: "domain",
		DNSAddress:IPAddress{
			Type: "IPv4",
			IPv4Address:"172.16.0.1",
		},
	}

	err = testDevice.SetDNS(dnsInformation)
	if err != nil {
		t.Error(err)
	}
}


// Todo get set protocols
func GetSetNetworkProtocols(t *testing.T){
	log.Println("Test GetSetNetworkProtocols")

	res, err := testDevice.GetNetworkProtocols()
	if err != nil {
		t.Error(err)
	}
	js := prettyJSON(&res)
	fmt.Println(js)

	// Todo set protocols
	//protocols := []models.NetworkProtocol{}
	//protocol :=NetworkProtocol{
	//	Name: "HTTP",
	//	Enabled: true,
	//	Port: 8080,
	//}
	//protocols = append(protocols, protocol)
	//protocol =NetworkProtocol{
	//	Name: "RTSP",
	//	Enabled: false,
	//	Port: 80,
	//}
	//protocols = append(protocols, protocol)
	//
	//res, err := testDevice.SetNetworkProtocols(protocols)
}

func GetSetScopes(t *testing.T)  {
	log.Println("Test GetSetScopes")

	res, err := testDevice.GetScopes()
	if err != nil {
		t.Error(err)
	}
	js := prettyJSON(&res)
	fmt.Println(js)

	scopes := []string{
		"onvif://www.onvif.org/name/HeroSpeed",
		"onvif://www.onvif.org/location/Guangzhou",
	}

	err = testDevice.SetScopes(scopes)
	if err != nil {
		t.Error(err)
	}
}

func AddRemoveScopes(t *testing.T)  {
	log.Println("Test AddRemoveScopes")

	res, err := testDevice.GetScopes()
	if err != nil {
		t.Error(err)
	}
	js := prettyJSON(&res)
	fmt.Println(js)

	scopes := []string{
		"onvif://www.onvif.org/location/Guangzhou",
	}

	res, err = testDevice.RemoveScopes(scopes)
	if err != nil {
		t.Error(err)
	}
	js = prettyJSON(&res)
	fmt.Println(js)

	scopes = []string{
		"onvif://www.onvif.org/location/Guangzhou",
	}

	err = testDevice.AddScopes(scopes)
	if err != nil {
		t.Error(err)
	}
}

func GetSetNetworkDefaultGateway(t *testing.T){
	log.Println("Test GetSetNetworkDefaultGateway")

	res, err := testDevice.GetNetworkDefaultGateway()
	if err != nil {
		t.Error(err)
	}
	js := prettyJSON(&res)
	fmt.Println(js)
	//gateway :=NetworkGateway{
	//	IPv4Address: "172.16.0.2",
	//}
	//res, err := testDevice.SetNetworkDefaultGateway(gateway)
}

func TestGetUsers(t *testing.T)  {
	log.Println("Test GetUsers")

	res, err := testDevice.GetUsers()
	if err != nil {
		t.Error(err)
	}
	js := prettyJSON(&res)
	fmt.Println(js)
}

func CreateSetDeleteUser(t *testing.T)  {
	log.Println("Test CreateSetDeleteUser")

	//users := []models.User{
	//	{
	//		"phuocnv2",
	//		"123456",
	//		"Administrator",
	//	},
	//}
	//res, err := testDevice.CreateUsers(users)
	//res, err := testDevice.SetUser(users[0])

	usernames := []string{
		"phuocnv2",
	}
	err := testDevice.DeleteUsers(usernames)
	if err != nil {
		t.Error(err)
	}
}

func TestGetRelayOutputs(t *testing.T){
	log.Println("Test GetRelayOutputs")

	res, err := testDevice.GetRelayOutputs()
	if err != nil {
		t.Error(err)
	}
	js := prettyJSON(&res)
	fmt.Println(js)
}

func TestGetServices(t *testing.T){
	log.Println("Test GetServices")

	res, err := testDevice.GetServices()
	if err != nil {
		t.Error(err)
	}
	js := prettyJSON(&res)
	fmt.Println(js)
}

func TestGetServiceCapabilities(t *testing.T){
	log.Println("Test GetServiceCapabilities")

	res, err := testDevice.GetServiceCapabilities()
	if err != nil {
		t.Error(err)
	}
	js := prettyJSON(&res)
	fmt.Println(js)
}

