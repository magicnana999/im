package errext

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	e := New(100, "login error")
	fmt.Println(e)
	fmt.Println(e.SetDetail("haha"))
	fmt.Println(e.SetDetail("hehe"))
	s, eee := e.JsonString()
	fmt.Println(string(s), eee)
}
