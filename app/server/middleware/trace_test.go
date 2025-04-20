package middleware_test

import (
	"app/server/middleware"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

func TestTracer(t *testing.T) {
	spanRecorder := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider()
	tp.RegisterSpanProcessor(spanRecorder)
	otel.SetTracerProvider(tp)

	tracer := tp.Tracer("test-tracer")

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	handler := middleware.Trace(tracer)(testHandler)
	handler.ServeHTTP(w, r)

	spans := spanRecorder.Ended()
	gotLen := len(spans)
	if gotLen != 1 {
		t.Errorf("want: 1, got: %d", gotLen)
	}

	gotAttr := spans[0].Attributes()
	wantAttr := []attribute.KeyValue{
		attribute.String(string(semconv.HTTPRequestMethodKey), "GET"),
		semconv.URLPath(r.URL.Path),
		semconv.URLScheme(r.URL.Scheme),
		semconv.HTTPResponseStatusCode(http.StatusOK),
	}
	if !reflect.DeepEqual(gotAttr, wantAttr) {
		t.Errorf("want: %v, got: %v", wantAttr, gotAttr)
	}
}
