package training

import (
	"context"
	"fmt"
	"github.com/timandy/routine"
	"time"
)

func handle(name string, ctx context.Context) {

	time.Sleep(time.Second)
	for {
		select {
		case <-ctx.Done():
			fmt.Println(name, "done")
			return
		default:
			fmt.Println(name, routine.Goid(), ctx.Value("root"))
			fmt.Println(name, routine.Goid(), ctx.Value("sub1"))
			fmt.Println(name, routine.Goid(), ctx.Value("sub2"))
			return
		}
	}
}

func case3WithDeadline() {
	root, _ := context.WithTimeout(context.Background(), time.Second*3)

	deadline := time.Now().Add(2 * time.Second)
	sub, _ := context.WithDeadline(context.Background(), deadline)

	fun := func(name string, ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				fmt.Println(name, "done")
				return
			case <-time.After(time.Second):
				fmt.Println(name, time.Second.Seconds())
			case <-time.After(time.Second * 2):
				fmt.Println(name, time.Second.Seconds()*2)
			case <-time.After(time.Second * 3):
				fmt.Println(name, time.Second.Seconds()*3)
			case <-time.After(time.Second * 4):
				fmt.Println(name, time.Second.Seconds()*4)
			}
		}
	}

	fun("root", root)
	fun("sub1", sub)

	time.Sleep(time.Second * 5)
	fmt.Println("ok")

}

func case1RootCancel() {
	ctx := context.WithValue(context.Background(), "root", "root")

	root, cancel := context.WithCancel(ctx)
	sub1 := context.WithValue(root, "sub1", "sub1")
	sub2 := context.WithValue(root, "sub2", "sub2")

	go handle("root", root)
	go handle("sub1", sub1)
	go handle("sub2", sub2)

	cancel()

	time.Sleep(5 * time.Second)
	fmt.Println("ok")
}

func case2RootValue() {
	ctx := context.WithValue(context.Background(), "root", "root")

	root, _ := context.WithCancel(ctx)
	sub1 := context.WithValue(root, "sub1", "sub1")
	sub2 := context.WithValue(root, "sub2", "sub2")

	go handle("root", root)
	go handle("sub1", sub1)
	go handle("sub2", sub2)

	time.Sleep(5 * time.Second)
	fmt.Println("ok")
}
