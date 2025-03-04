package redis

import (
	"context"
	"fmt"
	"github.com/timandy/routine"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestLock(t *testing.T) {

	initRedis()

	var wg sync.WaitGroup
	wg.Add(10)

	f := func() {
		fmt.Println(strconv.FormatInt(routine.Goid(), 10), "start")

		key := "im:test:lock"
		ctx := context.Background()
		val := "100"
		if AcquireLock(ctx, key, val, time.Second*10, 10) {
			defer ReleaseLock(ctx, key, val)
			time.Sleep(time.Second * 2)
			fmt.Println(strconv.FormatInt(routine.Goid(), 10), "done")
			wg.Done()
		}
	}

	for i := 0; i < 10; i++ {
		fmt.Println("start routine ", i)
		go f()
	}

	wg.Wait()
	fmt.Println("All routine is done")
}
