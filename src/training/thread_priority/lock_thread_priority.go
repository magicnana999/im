package thread_priority

import (
	"fmt"
	"golang.org/x/sys/unix"
	"runtime"
	"sync"
	"time"
)

type worker struct {
	name     string
	priority int
	locked   bool
	wg       sync.WaitGroup
}

func (w *worker) run() {

	w.wg.Add(1)

	if w.locked {
		runtime.LockOSThread()
	}

	unix.Setpriority(unix.PRIO_PROCESS, 0, w.priority)

	for i := 0; i < 50; i++ {
		if w.name == "high" {
			fmt.Println(i, w.name, time.Now().Format("2006-01-02 15:04:05:000"))
		}
		time.Sleep(time.Second * 1)
	}

	w.wg.Done()
}

func foo() {

	var wg sync.WaitGroup
	w1 := worker{name: "high", priority: -10, locked: true, wg: wg}
	w2 := worker{name: "middle", priority: 20, locked: false, wg: wg}

	for i := 0; i < 8; i++ {
		go w1.run()
	}

	for i := 0; i < 9999999; i++ {
		go w2.run()
	}

	wg.Wait()
}
