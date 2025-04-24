package server

import (
	"app/metrics"
	"app/server/middleware"
	"app/service"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/trace"
)

func NewServer(
	addr string,
	l *slog.Logger,
	m *metrics.Metrics,
	t trace.Tracer,
	s service.Service,
) *http.Server {
	r := chi.NewMux()

	r.Group(func(r chi.Router) {
		r.Use(middleware.Metrics(m))
		r.Use(middleware.Trace(t))
		r.Use(middleware.Logger(l, t))

		r.Route("/something", func(r chi.Router) {
			r.Get("/{id}", handleGetSomething(s))
			r.Post("/", handleCreateSomething(s))
		})
	})

	r.Get("/health", handleCheckHealth())
	r.Get("/metrics", promhttp.Handler().ServeHTTP)

	return &http.Server{
		Addr:    addr,
		Handler: r,
	}
}
