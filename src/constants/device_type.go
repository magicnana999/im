package constants

type DeviceType string

const (
	Mobile  DeviceType = "Mobile"
	Desktop            = "Desktop"
)

func (d DeviceType) String() string {
	return string(d)
}
