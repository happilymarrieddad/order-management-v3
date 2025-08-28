package middleware

import (
	"encoding/json"
	"net/http"
)

// WriteJSON is a helper for sending structured JSON responses.
func WriteJSON(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}

// WriteError is a convenience wrapper around WriteJSON for sending errors.
func WriteError(w http.ResponseWriter, code int, message string) {
	WriteJSON(w, code, map[string]string{"error": message})
}
