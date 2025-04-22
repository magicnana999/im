package main

import (
	"fmt"
	"sync/atomic"
	"unsafe"
)

func unsafePointer() {

	type node struct {
		index int
	}

	var val unsafe.Pointer

	fmt.Println(val == nil)
	fmt.Println(atomic.LoadPointer(&val) == nil)

	for i := 0; i < 10; i++ {

		go func() {
			n := &node{i}

			ok := atomic.CompareAndSwapPointer(&val, nil, unsafe.Pointer(n))
			if ok {
				nn := (*node)(atomic.LoadPointer(&val))
				fmt.Println(i, nn.index)
			}
			fmt.Println(i, ok)
		}()

	}
}
