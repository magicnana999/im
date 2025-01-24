package ip

import (
	"fmt"
	"testing"
)

func TestGetLocalIP(t *testing.T) {
	s, _ := GetLocalIP()
	fmt.Println(s)
}
