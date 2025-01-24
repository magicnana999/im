package json

import "encoding/json"

func IgnoreErrorMarshal(v interface{}) []byte {
	b, e := json.Marshal(v)

	if e != nil {
		return []byte("")
	}
	return b
}
