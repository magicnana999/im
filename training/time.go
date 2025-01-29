package training

import (
	"fmt"
	"time"
)

func times() {
	now := time.Now()
	fmt.Println(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())

	fmt.Println(now.String())

	fmt.Println(now.Format("2006-01-02 15:04:05"))

	fmt.Println(now.Unix())
	fmt.Println(now.UnixMilli())
	fmt.Println(now.UnixMicro())
	fmt.Println(now.UnixNano())

	lastHeartbeat := now.UnixMilli()
	intervl := 10 * time.Second

	go func() {
		for {
			time.Sleep(time.Second)
			n := time.Now()
			fmt.Println(n.UnixMilli(), lastHeartbeat, intervl.Milliseconds())
			if n.UnixMilli()-lastHeartbeat > intervl.Milliseconds() {
				fmt.Println("OK")
				return
			}
		}
	}()

	time.Sleep(intervl + time.Second)
	fmt.Println("OK")
}
