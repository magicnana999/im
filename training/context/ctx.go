package context

import (
	"context"
	"fmt"
	"time"
)

func CtxWithValue() {
	c := context.WithValue(context.Background(), "key", "parent")
	sub := context.WithValue(c, "key", "sub")

	fmt.Println(c.Value("key"))
	fmt.Println(sub.Value("key"))
}

func CtxWithCancel(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Done", ctx.Err())
			return
		case <-time.After(time.Millisecond * 200):
			fmt.Println("After 200 Millisecond")
			//default:
			//	time.Sleep(time.Second)
			//	fmt.Println("do")
		}
	}
}

func CtxWithDeadline(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Done", ctx.Err())
			return
		case <-time.After(time.Millisecond * 200):
			fmt.Println("After 200 Millisecond")
		}
	}
}

func CtxWithTimeout(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Done", ctx.Err())
			return
		case <-time.After(time.Millisecond * 200):
			fmt.Println("After 200 Millisecond")
		}
	}
}
