package onvif

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

var testDevice = Device{
	XAddr:    "http://192.168.0.11/onvif/device_service",
	User:     "admin",
	Password: "Admin123",
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

func intToString(src int) string {
	strInt := strconv.Itoa(src)
	return strInt
}

func float64ToString(src float64) string {
	strFloat64 := fmt.Sprint(src)
	return strFloat64
}

func boolToString(src bool) string {
	if src {
		return "true"
	}

	return "false"
}

// kiem tra co phai loi chung thuc hay khong
func CheckAuthorizedError(msg string) bool {
	msg = strings.ToLower(msg)
	return strings.Index(msg, "authorized") != -1
}
