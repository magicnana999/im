package utils

import (
	"errors"
	"fmt"
	"testing"
)

func TestAny2String(t *testing.T) {
	fmt.Println(Any2String(int(100)))
	fmt.Println(Any2String(int8(100)))
	fmt.Println(Any2String(int16(100)))
	fmt.Println(Any2String(int32(100)))

	fmt.Println(Any2String(uint(100)))
	fmt.Println(Any2String(uint8(100)))
	fmt.Println(Any2String(uint16(100)))
	fmt.Println(Any2String(uint32(100)))

	fmt.Println(Any2String(float32(100)))
	fmt.Println(Any2String(float64(100)))

	fmt.Println(Any2String(true))

	fmt.Println(Any2String([]string{"x", "y", "z"}))
	fmt.Println(Any2String([3]string{"x", "y", "z"}))

	fmt.Println(Any2String(map[string]int{"x": 1, "y": 2, "z": 3}))
	fmt.Println(Any2String(struct{ Name, Email string }{
		Name:  "Alice",
		Email: "alice@example.com",
	}))
	fmt.Println(Any2String(errors.New("hshs").Error()))

	slice := []string{"x", "y", "z"}
	fmt.Printf("Type of slice: %T\n", slice) // []string

	// This is an array
	array := [3]string{"x", "y", "z"}
	fmt.Printf("Type of array: %T\n", array) // [3]string

	fmt.Println(Any2String(errors.New("haha")))

	fmt.Println(Any2String(&struct{ Name, Email string }{
		Name:  "Alice",
		Email: "alice@example.com",
	}))
}
