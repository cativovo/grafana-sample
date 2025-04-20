package middleware

import (
	"app/ctxvalue"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/otel/trace"
)

func Logger(l *slog.Logger, t trace.Tracer) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, span := t.Start(r.Context(), "Logger middleware")
			defer span.End()

			logger := l.With(
				"trace_id", span.SpanContext().TraceID().String(),
				"method", r.Method,
				"path", r.URL.Path,
				"scheme", r.URL.Scheme,
			)

			logger.Info("Processing request")

			defer func() {
				ww := w.(middleware.WrapResponseWriter)
				logger.Info("Finished processing request", "status", ww.Status())
			}()

			ctx = ctxvalue.WithLogger(ctx, logger)
			r = r.WithContext(ctx)

			ww, ok := w.(middleware.WrapResponseWriter)
			if !ok {
				ww = middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			}

			h.ServeHTTP(ww, r)
		})
	}
}
