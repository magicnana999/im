package main

import (
	"fmt"
	"github.com/asynkron/goconsole"
	"time"

	"github.com/RussellLuo/timingwheel"
)

func main() {
	tw := timingwheel.NewTimingWheel(1*time.Second, 60)
	tw.Start()      // 启动时间轮
	defer tw.Stop() // 程序结束时停止

	for i := 1; i <= 60; i++ {
		tw.AfterFunc(time.Duration(i)*time.Second, func() {
			fmt.Printf("Hello, World! (%d seconds passed)\n", i)
		})
		fmt.Printf("add %d\n", i)

	}

	console.ReadLine()
	fmt.Println("Main program ends.")
}
