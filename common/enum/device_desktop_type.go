package enum

import "fmt"

type DeviceDesktopType int

const (
	Windows DeviceDesktopType = iota + 1
	Mac
	Linux
)

var devicePCNames = [...]string{
	"Windows", "Mac", "Linux",
}

func (d DeviceDesktopType) String() string {
	if d.Valid() {
		return devicePCNames[d-1]
	}
	return "Unknown"
}

func (d DeviceDesktopType) Code() int {
	return int(d)
}

func (d DeviceDesktopType) GetParser() (EnumParser, error) {
	return defaultDevicePCTypeParser, nil
}

func (d DeviceDesktopType) Valid() bool {
	return d >= Windows && d <= Mac
}

var defaultDevicePCTypeParser = DevicePCTypeParser{}

type DevicePCTypeParser struct{}

func (p DevicePCTypeParser) Parse(code int) (any, error) {
	if code >= int(Windows) && code <= int(Mac) {
		return DeviceDesktopType(code), nil
	}
	return nil, fmt.Errorf("invalid DeviceDesktopType code: %d", code)
}
