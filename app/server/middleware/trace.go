package middleware

import (
	"app/otel"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

func Trace(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer.Start(r.Context(), fmt.Sprintf("%s %s", r.Method, r.URL.Path))
		defer span.End()

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		defer func() {
			span.SetAttributes(
				attribute.String(string(semconv.HTTPRequestMethodKey), r.Method),
				semconv.URLPath(r.URL.Path),
				semconv.URLScheme(r.URL.Scheme),
				semconv.HTTPResponseStatusCode(ww.Status()),
			)
		}()

		r = r.WithContext(ctx)
		h.ServeHTTP(ww, r)
	})
}
