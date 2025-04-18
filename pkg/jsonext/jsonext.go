package jsonext

import (
	"encoding/json"
	"google.golang.org/protobuf/proto"
)

func MarshalNoErr(v interface{}) []byte {
	j, err := json.Marshal(v)
	if err != nil {
		return []byte("")
	}
	return j
}

func PbMarshalNoErr(m proto.Message) []byte {
	j, err := proto.Marshal(m)
	if err != nil {
		return []byte("")
	}
	return j
}
