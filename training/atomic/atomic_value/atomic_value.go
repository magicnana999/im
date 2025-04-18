package main

import (
	"fmt"
	"sync/atomic"
)

var token atomic.Value

func main() {
	token.Store("token")
	fmt.Println(token.Load().(string))
}
