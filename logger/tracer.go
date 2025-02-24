package logger

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
	traceing "go.opentelemetry.io/otel/trace"
)

var (
	Tracer traceing.Tracer
)

func InitTracer(name string) {

	Tracer = otel.Tracer(name)
	tp := trace.NewTracerProvider(trace.WithSampler(trace.NeverSample()))
	otel.SetTracerProvider(tp)
}

func NewSpan(ctx context.Context, name string) context.Context {
	c, _ := Tracer.Start(ctx, name)
	return c
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
