package enum

type Enum interface {
	Name() string
	Code() uint8
	Valid() bool
}
