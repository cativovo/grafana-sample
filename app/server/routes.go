package server

import (
	"app"
	"app/ctxvalue"
	"app/otel"
	"app/service"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func handleGetSomething(s service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer.Start(r.Context(), "handleGetSomething")
		defer span.End()

		logger := ctxvalue.Logger(ctx)

		something, err := s.GetSomething(ctx, chi.URLParam(r, "id"))
		if err != nil {
			msg := app.GetErrorMessage(err)
			logger.Error(msg, "error", err)
			newHTTPError(w, getStatusCode(err), msg)
			return
		}

		json.NewEncoder(w).Encode(something)
	}
}

func handleCreateSomething(s service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer.Start(r.Context(), "handleGetSomething")
		defer span.End()

		logger := ctxvalue.Logger(ctx)

		var c service.CreateSomethingReq
		if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
			logger.Error("Failed to unmarshal the body", "error", err)
			newHTTPError(w, http.StatusBadRequest, "Unable to parse the request")
			return
		}

		something, err := s.CreateSomething(ctx, c)
		if err != nil {
			msg := app.GetErrorMessage(err)
			logger.Error(msg, "error", err)
			newHTTPError(w, getStatusCode(err), msg)
			return
		}

		json.NewEncoder(w).Encode(something)
	}
}

func handleCheckHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{
			"status": "healthy",
		})
	}
}
