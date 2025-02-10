package enum

import "fmt"

type OSType int

const (
	Windows OSType = iota + 1
	MacOS
	Linux
	Ios
	Xiaomi
	Huawei
	Samsung
	Honor
	Oppo
	Vivo
)

var osNames = [...]string{
	"Windows", "MacOS", "Linux", "iOS", "Xiaomi", "Huawei", "Samsung", "Honor", "Oppo", "Vivo",
}

func (o OSType) String() string {
	if o.Valid() {
		return osNames[o-1]
	}
	return "Unknown"
}

func (o OSType) Code() int {
	return int(o)
}

func (o OSType) GetParser() (EnumParser, error) {
	return DefaultOSTypeParser, nil
}

func (o OSType) Valid() bool {
	return o >= Windows && o <= Vivo
}

func (o OSType) GetDeviceType() DeviceType {
	switch o.Code() {
	case int(Windows), int(MacOS), int(Linux):
		return Desktop
	case int(Ios), int(Xiaomi), int(Huawei), int(Samsung), int(Honor), int(Oppo), int(Vivo):
		return Mobile
	default:
		return DeviceType(0)
	}
}

var (
	DefaultOSTypeParser = OSTypeParser{}
)

type OSTypeParser struct{}

func (p OSTypeParser) Parse(code int) (any, error) {
	if code >= int(Windows) && code <= int(Vivo) {
		return OSType(code), nil
	}
	return nil, fmt.Errorf("invalid OSType code: %d", code)
}
