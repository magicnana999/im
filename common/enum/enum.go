package enum

type Enum interface {
	String() string
	Code() int
	Valid() bool
	GetParser() (EnumParser, error)
}

type EnumParser interface {
	Parse(int) (any, error)
}
