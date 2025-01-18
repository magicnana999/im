package enum

type DeviceType uint8

const (
	Mobile DeviceType = iota + 1
	PC
)

func (d DeviceType) Name() string {
	if d.Valid() {
		return [...]string{"Mobile", "PC"}[d-1]
	}
	return ""
}

func (d DeviceType) Code() uint8 {
	return uint8(d)
}

func (d DeviceType) Valid() bool {
	return d.Code() >= uint8(Mobile) && d.Code() <= uint8(PC)
}
