package middleware

import (
	"app/metrics"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5/middleware"
)

func Metrics(m *metrics.Metrics) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			timer := m.Latency(r.URL.Path)
			defer timer.ObserveDuration()

			ww, ok := w.(middleware.WrapResponseWriter)
			if !ok {
				ww = middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			}
			defer func() {
				m.Count(strconv.Itoa(ww.Status()))
			}()

			h.ServeHTTP(ww, r)
		})
	}
}
