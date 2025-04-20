package middleware

import (
	"app/ctxvalue"
	"app/otel"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

func Logger(logger *slog.Logger) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, span := otel.Tracer.Start(r.Context(), "Logger middleware")
			defer span.End()

			l := logger.With(
				slog.String("trace_id", span.SpanContext().TraceID().String()),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("scheme", r.URL.Scheme),
			)

			l.Info("Processing request")

			defer func() {
				ww := w.(middleware.WrapResponseWriter)
				l.Info("Finished processing request", slog.Int("status", ww.Status()))
			}()

			ctx = ctxvalue.WithLogger(ctx, l)
			r = r.WithContext(ctx)

			ww, ok := w.(middleware.WrapResponseWriter)
			if !ok {
				ww = middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			}

			h.ServeHTTP(ww, r)
		})
	}
}
