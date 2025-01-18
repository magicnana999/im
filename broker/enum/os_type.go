package enum

type OSType uint8

const (
	OSWindows OSType = iota + 1
	MacOS
	Linux
	OSIos
	Xiaomi
	Huawei
	Samsung
	Honor
	Oppo
	Vivo
)

func (o OSType) Name() string {
	if o.Valid() {
		return [...]string{"Windows", "MacOS", "Linux", "iOS", "Xiaomi", "Huawei", "Samsung", "Honor", "Oppo", "Vivo"}[o-1]
	}
	return ""
}

func (o OSType) Code() uint8 {
	return uint8(o)
}

func (o OSType) Valid() bool {
	return o.Code() >= uint8(OSWindows) && o.Code() <= uint8(Vivo)
}

func (o OSType) PushAvailable() bool {
	return o.Code() >= uint8(OSIos)
}
