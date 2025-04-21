package consurrent

import (
	"fmt"
	"sync"
	"time"
)

func cas() {
	var m sync.Map

	key := "hello_world"
	for i := 0; i < 10; i++ {
		go func() {
			v, ok := m.LoadOrStore(key, i)
			if ok {
				fmt.Println(i, "OK", v)
			} else {
				fmt.Println(i, "ERR", v)
			}
		}()
	}

	time.Sleep(time.Second)
}
