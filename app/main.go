package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	shutdown := initTracerProvider()
	defer shutdown()

	mux := http.NewServeMux()
	mux.Handle("POST /process", handleProcess(logger))

	port := "6000"
	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	logger.Info("Server started", "port", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}

func handleProcess(logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracer.Start(r.Context(), "HTTP POST /process")
		defer span.End()

		logger.Info("Starting process")
		process(ctx)
		logger.Info("Process completed")

		fmt.Fprint(w, `{"message": "ok"}`)
	})
}

func process2(ctx context.Context) {
	_, span := tracer.Start(ctx, "process2")
	defer span.End()

	time.Sleep(time.Second * 1)
}

func process(ctx context.Context) {
	ctx, span := tracer.Start(ctx, "process")
	defer span.End()

	time.Sleep(time.Second * 2)
	process2(ctx)
}
