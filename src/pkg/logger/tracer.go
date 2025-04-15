package logger

import (
	"context"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
	traceing "go.opentelemetry.io/otel/trace"
)

var (
	tracer traceing.Tracer
)

func InitTracer(name string) {

	tracer = otel.Tracer(name)
	tp := trace.NewTracerProvider(trace.WithSampler(trace.NeverSample()))
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(b3.New())

}

func NewSpan(ctx context.Context, name string) (context.Context, traceing.Span) {
	c, s := tracer.Start(ctx, name)
	return c, s
}

func EndSpan(ctx context.Context) {
	span := traceing.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		span.End()
	}
}

func TraceID(ctx context.Context) string {
	span := traceing.SpanContextFromContext(ctx)
	if span.IsValid() {
		return span.TraceID().String()
	}
	return ""
}

func SpanID(ctx context.Context) string {
	span := traceing.SpanContextFromContext(ctx)

	if span.IsValid() {
		return span.SpanID().String()
	}
	return ""
}

func ShutdownTracer() error {
	if tp, ok := otel.GetTracerProvider().(*trace.TracerProvider); ok {
		return tp.Shutdown(context.Background())
	}
	return nil
}
