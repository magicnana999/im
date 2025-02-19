package constants

type OSType string

const (
	Windows OSType = "Windows"
	MacOS          = "MaxOS"
	Linux          = "Linux"
	Ios            = "iOS"
	Xiaomi         = "Xiaomi"
	Huawei         = "Huawei"
	Samsung        = "Samsung"
	Honor          = "Honer"
	Oppo           = "Oppo"
	Vivo           = "Vivo"
)

func (o OSType) String() string {
	return string(o)
}

func (o OSType) GetDeviceType() DeviceType {
	switch o {
	case Windows, MacOS, Linux:
		return Desktop
	case Ios, Xiaomi, Huawei, Samsung, Honor, Oppo, Vivo:
		return Mobile
	default:
		return Desktop
	}
}
