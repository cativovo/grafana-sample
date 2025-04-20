package main

import (
	"app/otel"
	"app/server"
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	shutdown := otel.InitTracerProvider()
	defer shutdown()

	l := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	port := 6000
	s := server.NewServer(fmt.Sprintf(":%d", port), l)

	l.Info("Server is ready", slog.Int("port", port))

	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
