package controller

import (
	"encoding/json"
	"net/http"
)

// TODO: consider about an interface with Error and Respond methods
func Error(w http.ResponseWriter, r *http.Request, statusCode int, err error) {
	Respond(w, r, statusCode, map[string]string{"error": err.Error()})
}

func Respond(w http.ResponseWriter, r *http.Request, statusCode int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
