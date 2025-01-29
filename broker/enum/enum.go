package enum

type Enum interface {
	Name() string
	Code() int
	Valid() bool
	GetParser() (EnumParser, error)
}

type EnumParser interface {
	Parse(int) (any, error)
}
