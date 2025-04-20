package server

import (
	"app/server/middleware"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewServer(addr string, l *slog.Logger) *http.Server {
	r := chi.NewMux()
	r.Use(middleware.Trace)
	r.Use(middleware.Logger(l))

	r.Post("/process", handleProcess())

	return &http.Server{
		Addr:    addr,
		Handler: r,
	}
}
