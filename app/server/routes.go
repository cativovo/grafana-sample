package server

import (
	"app/ctxvalue"
	"app/otel"
	"fmt"
	"net/http"
)

func handleProcess() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer.Start(r.Context(), "handleProcess")
		defer span.End()

		logger := ctxvalue.Logger(r.Context())

		logger.Info("Starting 2 process1")
		process1(ctx)
		process1(ctx)
		logger.Info("Process completed")

		fmt.Fprint(w, `{"message": "ok"}`)
	}
}
