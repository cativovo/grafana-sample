package server

import (
	"app/server/middleware"
	"app/service"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel/trace"
)

func NewServer(
	addr string,
	l *slog.Logger,
	t trace.Tracer,
	s service.Service,
) *http.Server {
	r := chi.NewMux()
	r.Use(middleware.Trace(t))
	r.Use(middleware.Logger(l, t))

	r.Route("/something", func(r chi.Router) {
		r.Get("/{id}", handleGetSomething(s))
		r.Post("/", handleCreateSomething(s))
	})

	r.Get("/health", handleCheckHealth())

	return &http.Server{
		Addr:    addr,
		Handler: r,
	}
}
