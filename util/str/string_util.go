package str

import "strings"

func IsNotBlank(s string) bool {
	str := strings.TrimSpace(s)
	return str != "" && str != "nil"
}

func IsBlank(s string) bool {
	str := strings.TrimSpace(s)
	return str == "" || str == "nil"
}
