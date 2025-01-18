package tracer

import (
	"context"
	"fmt"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/trace"
	traceing "go.opentelemetry.io/otel/trace"
	"os"
)

const name = "go.opentelemetry.io/otel/tracer"

var (
	Tracer = otel.Tracer(name)
	PD     *trace.TracerProvider
)

func init() {

	otel.SetTextMapPropagator(b3.New())

	f, _ := os.Create("trace.txt")

	exp, _ := stdouttrace.New(
		stdouttrace.WithWriter(f),
		stdouttrace.WithPrettyPrint(),
		stdouttrace.WithoutTimestamps(),
	)

	PD = trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithSampler(trace.NeverSample()),
	)

	otel.SetTracerProvider(PD)

}

func demo() {
	ctx, span := Tracer.Start(context.Background(), "root-span")
	defer span.End()

	sc := traceing.SpanContextFromContext(ctx)
	if sc.IsValid() {
		fmt.Println(sc.TraceID().String(), sc.SpanID().String())
		fmt.Println(Tracer)
	}

	ctx, span = Tracer.Start(context.WithValue(ctx, "key", "value"), "sub-span")
	defer span.End()
	sc = traceing.SpanContextFromContext(ctx)
	if sc.IsValid() {
		fmt.Println(sc.TraceID().String(), sc.SpanID().String())
		fmt.Println(Tracer)

	}
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
