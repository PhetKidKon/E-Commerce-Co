// Package response is the shared API-response helper used by every service.
// This package has no main() so it can never run on its own; services import it.
package response

import (
	"encoding/json"
	"net/http"
)

// ApiResponse is the standard envelope returned by all services.
type ApiResponse[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    T      `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

// JSON writes any body with the given status code.
func JSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

// Success writes a 200 response wrapping data.
func Success[T any](w http.ResponseWriter, message string, data T) {
	JSON(w, http.StatusOK, ApiResponse[T]{Success: true, Message: message, Data: data})
}

// Error writes an error response with the given status code.
func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, ApiResponse[any]{Success: false, Error: message})
}
