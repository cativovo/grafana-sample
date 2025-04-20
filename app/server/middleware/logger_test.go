package middleware_test

import (
	"app/ctxvalue"
	"app/server/middleware"
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/otel/trace/noop"
)

var tracer = noop.NewTracerProvider().Tracer("noop tracer")

func TestLoggerMiddleware(t *testing.T) {
	t.Run("middleware.WrapResponseWriter", func(t *testing.T) {
		var buf bytes.Buffer
		logger := slog.New(slog.NewJSONHandler(&buf, nil))

		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, ok := w.(chimiddleware.WrapResponseWriter)
			if !ok {
				t.Errorf("middleware.WrapResponseWriter is unavailable on the writer. add the interface methods.")
			}
		})

		r := httptest.NewRequest("GET", "/", nil)
		w := chimiddleware.NewWrapResponseWriter(httptest.NewRecorder(), 1)

		handler := middleware.Logger(logger, tracer)(testHandler)
		handler.ServeHTTP(w, r)
	})

	t.Run("logger contents", func(t *testing.T) {
		var buf bytes.Buffer
		logger := slog.New(slog.NewJSONHandler(&buf, nil))

		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l := ctxvalue.Logger(r.Context())
			l.Info("Test", "test", "test")
			w.WriteHeader(http.StatusOK)
		})

		r := httptest.NewRequest("GET", "/", nil)
		w := chimiddleware.NewWrapResponseWriter(httptest.NewRecorder(), 1)

		handler := middleware.Logger(logger, tracer)(testHandler)
		handler.ServeHTTP(w, r)

		testCases := []struct {
			wantPattern string
			wantMatches int
		}{
			{
				wantPattern: `"method":"GET"`,
				wantMatches: 3,
			},
			{
				wantPattern: `"path":"/"`,
				wantMatches: 3,
			},
			{
				wantPattern: `"scheme":""`,
				wantMatches: 3,
			},
			{
				wantPattern: `"msg":"Processing request"`,
				wantMatches: 1,
			},
			{
				wantPattern: `"msg":"Finished processing request"`,
				wantMatches: 1,
			},
			{
				wantPattern: `"status":200`,
				wantMatches: 1,
			},
			{
				wantPattern: `"msg":"Test"`,
				wantMatches: 1,
			},
			{
				wantPattern: `"test":"test"`,
				wantMatches: 1,
			},
			{
				wantPattern: `"level":"INFO"`,
				wantMatches: 3,
			},
			{
				wantPattern: `"time":".*"`,
				wantMatches: 3,
			},
		}

		str := buf.String()
		for _, tc := range testCases {
			re := regexp.MustCompile(tc.wantPattern)
			matches := re.FindAllString(str, -1)
			got := len(matches)
			if got != tc.wantMatches {
				t.Errorf("%s: want %d, got %d\n %s", tc.wantPattern, tc.wantMatches, got, str)
			}
		}
	})
}
