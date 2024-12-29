package logger

import (
	"context"
	"fmt"
	"testing"
)

func TestStartSpan(t *testing.T) {
	ctx := CurrentSpan(context.Background(), "hello world")
	tid := TraceID(ctx)
	sid := SpanID(ctx)
	fmt.Println(tid)
	fmt.Println(sid)
	Log.Infof("%s %s %s", tid, sid, "hadf啥地方啥地方是")

}
