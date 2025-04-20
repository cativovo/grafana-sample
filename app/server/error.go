package server

import (
	"app"
	"encoding/json"
	"net/http"
)

var code = map[app.ErrCode]int{
	app.ErrCodeInvalid:  http.StatusBadRequest,
	app.ErrCodeNotFound: http.StatusNotFound,
	app.ErrCodeConflict: http.StatusConflict,
	app.ErrCodeInternal: http.StatusInternalServerError,
}

func getStatusCode(err error) int {
	k := app.GetErrorCode(err)
	c, ok := code[k]
	if !ok {
		return http.StatusInternalServerError
	}

	return c
}

func newHTTPError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	p := map[string]any{
		"error": message,
	}
	json.NewEncoder(w).Encode(p)
}
