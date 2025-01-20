package core

type Codec interface {
	Encode(v any) ([]byte, error)
	Decode(data []byte, v any) error
}
