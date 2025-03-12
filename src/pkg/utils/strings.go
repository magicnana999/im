package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func IsNotBlank(s string) bool {
	str := strings.TrimSpace(s)
	return str != "" && str != "nil"
}

func IsBlank(s string) bool {
	str := strings.TrimSpace(s)
	return str == "" || str == "nil"
}

func Any2String(value any) (string, error) {
	t := reflect.TypeOf(value)

	switch t.Kind() {
	case reflect.Int:
		return strconv.FormatInt(int64(value.(int)), 10), nil
	case reflect.Int8:
		return strconv.FormatInt(int64(value.(int8)), 10), nil
	case reflect.Int16:
		return strconv.FormatInt(int64(value.(int16)), 10), nil
	case reflect.Int32:
		return strconv.FormatInt(int64(value.(int32)), 10), nil
	case reflect.Int64:
		return strconv.FormatInt(value.(int64), 10), nil
	case reflect.Uint:
		return strconv.FormatUint(uint64(value.(uint)), 10), nil
	case reflect.Uint8:
		return strconv.FormatUint(uint64(value.(uint8)), 10), nil
	case reflect.Uint16:
		return strconv.FormatUint(uint64(value.(uint16)), 10), nil
	case reflect.Uint32:
		return strconv.FormatUint(uint64(value.(uint32)), 10), nil
	case reflect.Uint64:
		return strconv.FormatUint(value.(uint64), 10), nil
	case reflect.Float32:
		return fmt.Sprintf("%f", value.(float32)), nil
	case reflect.Float64:
		return fmt.Sprintf("%f", value.(float64)), nil
	case reflect.String:
		return value.(string), nil
	case reflect.Bool:
		return strconv.FormatBool(value.(bool)), nil
	case reflect.Slice, reflect.Array, reflect.Map, reflect.Struct:
		js, err := json.Marshal(value)
		if err != nil {
			return "", fmt.Errorf("error json marshal slice: %v", err)
		}
		return string(js), nil
	case reflect.Ptr:
		if t.Implements(reflect.TypeOf((*error)(nil)).Elem()) {
			return value.(error).Error(), nil
		}

		js, err := json.Marshal(value)
		if err != nil {
			return "", fmt.Errorf("error json marshal slice: %v", err)
		}
		return string(js), nil
	default:
		return "", fmt.Errorf("error: unsupported type %v", t.Kind())
	}
}
