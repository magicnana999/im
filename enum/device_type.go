package enum

import "fmt"

type DeviceType int

const (
	Mobile DeviceType = iota + 1
	Desktop
)

var deviceNames = [...]string{
	"Mobile", "Desktop",
}

func (d DeviceType) String() string {
	if d.Valid() {
		return deviceNames[d-1]
	}
	return ""
}

func (d DeviceType) Code() int {
	return int(d)
}

func (d DeviceType) GetParser() (EnumParser, error) {
	return defaultDeviceTypeParser, nil
}

func (d DeviceType) Valid() bool {
	return d >= Mobile && d <= Desktop
}

var defaultDeviceTypeParser = DeviceTypeParser{}

type DeviceTypeParser struct{}

func (p DeviceTypeParser) Parse(code int) (any, error) {
	if code >= int(Mobile) && code <= int(Desktop) {
		return DeviceType(code), nil
	}
	return nil, fmt.Errorf("invalid DeviceType code: %d", code)
}
