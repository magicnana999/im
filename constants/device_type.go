package constants

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

func (d DeviceType) Valid() bool {
	return d >= Mobile && d <= Desktop
}
