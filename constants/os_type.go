package constants

const (
	Unknown = "Unknown"
	Windows = "Windows"
	MacOS   = "MaxOS"
	Linux   = "Linux"
	Ios     = "iOS"
	Xiaomi  = "Xiaomi"
	Huawei  = "Huawei"
	Samsung = "Samsung"
	Honor   = "Honer"
	Oppo    = "Oppo"
	Vivo    = "Vivo"
)

func GetDeviceType(o string) DeviceType {
	switch o {
	case Windows, MacOS, Linux:
		return Desktop
	case Ios, Xiaomi, Huawei, Samsung, Honor, Oppo, Vivo:
		return Mobile
	default:
		return Desktop
	}
}
