package enum

import "fmt"

type DeviceMobileType int

const (
	Android DeviceMobileType = iota + 1
	IOS
)

var deviceMobileNames = [...]string{
	"Android", "iOS",
}

func (d DeviceMobileType) Name() string {
	if d.Valid() {
		return deviceMobileNames[d-1]
	}
	return "Unknown"
}

func (d DeviceMobileType) Code() int {
	return int(d)
}

func (d DeviceMobileType) GetParser() (EnumParser, error) {
	return defaultDeviceMobileTypeParser, nil
}

func (d DeviceMobileType) Valid() bool {
	return d >= Android && d <= IOS
}

var defaultDeviceMobileTypeParser = DeviceMobileTypeParser{}

type DeviceMobileTypeParser struct{}

func (p DeviceMobileTypeParser) Parse(code int) (any, error) {
	if code >= int(Android) && code <= int(IOS) {
		return DeviceMobileType(code), nil
	}
	return nil, fmt.Errorf("invalid DeviceMobileType code: %d", code)
}
