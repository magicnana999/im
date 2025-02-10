package main

import (
	"fmt"
	"testing"
)

func Test_parseFlag(t *testing.T) {

	//os.Args = []string{"cmd", "--name", "broker123", "--interval", "40"}
	option := parseFlag()
	fmt.Println(option)

}
