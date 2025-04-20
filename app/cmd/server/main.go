package main

import (
	"app/otel"
	"app/repository"
	"app/server"
	"app/service"
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	tracer, shutdown := otel.NewTracer()
	defer shutdown()

	repo := repository.NewRepository(tracer)
	service := service.NewService(tracer, repo)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	port := 6000
	s := server.NewServer(
		fmt.Sprintf(":%d", port),
		logger,
		tracer,
		service,
	)

	logger.Info("Server is ready", "port", port)

	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
