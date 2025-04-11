package main

import (
	"fmt"
	"github.com/RussellLuo/timingwheel"
	"sync/atomic"
	"time"
)

func main() {

	var count atomic.Int64

	wheelsize := int64(1)

	tw := timingwheel.NewTimingWheel(1*time.Second, wheelsize)
	tw.Start()
	defer tw.Stop()

	// 配置测试参数
	connCount := 900000

	for i := 0; i < connCount; i++ {
		delay := time.Duration(i%int(wheelsize)) * time.Second
		tw.AfterFunc(delay, func() {
			count.Add(1)
		})
	}
	time.Sleep(1 * time.Second) // 运行 2 分钟
	fmt.Println(count)
}
