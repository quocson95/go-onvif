package onvif

import (
	"encoding/json"
	"strconv"
	"strings"
)

var testDevice = Device{
	XAddr: "http://192.168.1.75:5000/onvif/device_service",
}

func interfaceToString(src interface{}) string {
	str, _ := src.(string)
	return str
}

func interfaceToBool(src interface{}) bool {
	strBool := interfaceToString(src)
	return strings.ToLower(strBool) == "true"
}

func interfaceToInt(src interface{}) int {
	strNumber := interfaceToString(src)
	number, _ := strconv.Atoi(strNumber)
	return number
}

func prettyJSON(src interface{}) string {
	result, _ := json.MarshalIndent(&src, "", "    ")
	return string(result)
}

func interfaceToFloat64(src interface{}) float64 {
	strNumber := interfaceToString(src)
	number, _ := strconv.ParseFloat(strNumber, 64)
	return number
}

func intToString(src int) string  {
	strInt := strconv.Itoa(src)
	return strInt
}
func boolToString(src bool) string{
	if src {
		return "true"
	}

	return "false"
}