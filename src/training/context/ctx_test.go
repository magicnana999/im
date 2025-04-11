package context

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestCtxWithValue(t *testing.T) {
	CtxWithValue()
}

func TestCtxWithCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	go CtxWithCancel(ctx)
	time.Sleep(5 * time.Second)
	cancel()
	time.Sleep(1 * time.Second)
	fmt.Println("OK")
}

func TestCtxWithDeadline(t *testing.T) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second))
	go CtxWithDeadline(ctx)
	time.Sleep(5 * time.Second)
	cancel() //已经deadline，调用无效
	fmt.Println("OK")
}

func TestCtxWithTimeout(t *testing.T) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second))
	go CtxWithTimeout(ctx)
	time.Sleep(5 * time.Second)
	cancel() //已经deadline，调用无效
	fmt.Println("OK")
}
