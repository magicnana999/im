package enum

type DevicePCType uint8

const (
	Windows DevicePCType = iota + 1
	Mac
)

func (d DevicePCType) Name() string {
	if d.Valid() {
		return [...]string{"Windows", "Mac"}[d-1]
	}
	return ""
}

func (d DevicePCType) Code() uint8 {
	return uint8(d)
}

func (d DevicePCType) Valid() bool {
	return d.Code() >= uint8(Windows) && d.Code() <= uint8(Mac)
}
