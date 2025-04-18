package define

type OSType string

const (
	Windows OSType = "Windows"
	MacOS   OSType = "MacOS"
	Linux   OSType = "Linux"
	Ios     OSType = "iOS"
	Xiaomi  OSType = "Xiaomi"
	Huawei  OSType = "Huawei"
	Samsung OSType = "Samsung"
	Honor   OSType = "Honor"
	Oppo    OSType = "Oppo"
	Vivo    OSType = "Vivo"
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
