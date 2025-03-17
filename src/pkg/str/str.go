package str

import (
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

func ConvertSS2I64S(s []string) ([]int64, error) {
	result := make([]int64, 0, len(s)) // 预分配足够的容量
	for _, s := range s {
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, err
		}
		result = append(result, i)
	}
	return result, nil
}
