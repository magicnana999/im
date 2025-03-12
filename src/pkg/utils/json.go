package utils

import "encoding/json"

func IgnoreErrorMarshal(v any) []byte {
	b, e := json.Marshal(v)

	if e != nil {
		return []byte("")
	}
	return b
}
