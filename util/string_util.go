package util

func IsNotBlank(s *string) bool {
	return s != nil && *s != "" && *s != "nil"
}

func IsBlank(s *string) bool {
	return s == nil || *s == "" || *s == "nil"
}
