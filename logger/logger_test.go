package logger

import (
	"context"
	"fmt"
	"github.com/timandy/routine"
	"sync"
	"testing"
)

func Test_demo(t *testing.T) {
	{
		Info("无Trace,哈哈")
		Info("无Trace,呵呵")
	}

	{
		c := NewSpan(context.Background(), "root")
		Info("有Trace,哈哈1")
		Info("有Trace,哈哈2")
		EndSpan(c)
		Info("无Trace,呵呵")
	}
}

var (
	threadLocal = routine.NewThreadLocal[string]()
)

func Test_timandy(t *testing.T) {
	{
		fmt.Println(routine.Goid())

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			fmt.Println(routine.Goid())
			wg.Done()
		}()
		wg.Wait()
	}

	{
		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			threadLocal.Set("111")
			fmt.Println(threadLocal.Get())
			wg.Done()
		}()

		go func() {
			threadLocal.Set("222")
			fmt.Println(threadLocal.Get())
			wg.Done()
		}()
		wg.Wait()
	}

}
