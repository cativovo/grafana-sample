package main

import (
	"context"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

// OTLP Exporter
func newOTLPExporter(ctx context.Context) (sdktrace.SpanExporter, error) {
	otlpEndpoint := os.Getenv("OTLP_ENDPOINT")
	if otlpEndpoint == "" {
		panic("You MUST set OTLP_ENDPOINT env variable!")
	}

	// Change default HTTPS -> HTTP
	insecureOpt := otlptracehttp.WithInsecure()

	// Update default OTLP reciver endpoint
	endpointOpt := otlptracehttp.WithEndpoint(otlpEndpoint)

	return otlptracehttp.New(ctx, insecureOpt, endpointOpt)
}

// TracerProvider is an OpenTelemetry TracerProvider.
// It provides Tracers to instrumentation so it can trace operational flow through a system.
func newTraceProvider(ctx context.Context, exp sdktrace.SpanExporter) *sdktrace.TracerProvider {
	// Ensure default SDK resources and the required service name are set.
	r, err := resource.New(ctx, resource.WithAttributes(semconv.ServiceNameKey.String("app")))
	if err != nil {
		panic(err)
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(r),
	)
}

func initTracerProvider() (shutdown func()) {
	ctx := context.Background()
	exp, err := newOTLPExporter(ctx)
	if err != nil {
		panic(err)
	}

	tp := newTraceProvider(ctx, exp)

	otel.SetTracerProvider(tp)

	tracer = tp.Tracer("app")

	return func() {
		if err := tp.Shutdown(ctx); err != nil {
			panic(err)
		}
	}
}
