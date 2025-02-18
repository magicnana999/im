package constants

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
