package main

import (
	"fmt"
	console "github.com/asynkron/goconsole"
	"time"
)

func main() {
	ch := make(chan int)

	go func() {
		for {
			time.Sleep(1 * time.Second)
			ch <- time.Now().Second()
		}
	}()

	go func() {
		for {
			select {
			case t := <-ch:
				fmt.Println(t)
			}
		}
	}()

	console.ReadLine()
}
