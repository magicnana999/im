package atomic

import (
	"fmt"
	"sync"
	"time"
)

type token struct {
	token string
}

func foo() {
	t := &token{"version1"}

	var wg sync.WaitGroup
	wg.Add(2)

	var w sync.WaitGroup
	w.Add(1)
	go func(t *token) {
		fmt.Println(t.token)
		w.Wait()
		time.Sleep(time.Second * 10)
		fmt.Println(t.token)
		wg.Done()
	}(t)

	time.Sleep(time.Second * 2)

	go func(t *token) {
		t.token = "version2"
		wg.Done()
		w.Done()
	}(t)

	wg.Wait()
	fmt.Println("OK")
}
