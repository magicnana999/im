package enum

type DeviceMobileType uint8

const (
	Android DeviceMobileType = iota + 1
	IOS
)

func (d DeviceMobileType) Name() string {
	if d.Valid() {
		return [...]string{"Android", "IOS"}[d-1]
	}
	return ""
}

func (d DeviceMobileType) Code() uint8 {
	return uint8(d)
}

func (d DeviceMobileType) Valid() bool {
	return d.Code() >= uint8(Android) && d.Code() <= uint8(IOS)
}
