package tracer

import (
	"context"
	"fmt"
	"testing"
)

func Test_all(t *testing.T) {
	ctx := context.Background()
	ctx = NewSpan(ctx, "root")
	fmt.Println(TraceID(ctx), SpanID(ctx))
	EndSpan(ctx)

	ctx = NewSpan(ctx, "sub")
	fmt.Println(TraceID(ctx), SpanID(ctx))
	EndSpan(ctx)
}

func Test_demo(t *testing.T) {
	demo()
}
