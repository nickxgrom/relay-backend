package controller

import (
	"encoding/json"
	"net/http"
	"relay-backend/internal/utils/exception"
)

// TODO: consider about an interface with Error and Respond methods
func Error(w http.ResponseWriter, r *http.Request, statusCode int, err error) {
	Respond(w, r, statusCode, map[string]string{"error": err.Error()})
}

// TODO: replace all errors to this method
func HTTPError(w http.ResponseWriter, r *http.Request, err exception.Exception) {
	Respond(w, r, err.StatusCode, map[string]interface{}{"error": err.Error(), "details": err.Details})
}

func Respond(w http.ResponseWriter, r *http.Request, statusCode int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
