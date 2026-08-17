package response

import (
	"encoding/json"
	"net/http"
)

// APIResponse is the single envelope every endpoint returns (success or error).
type APIResponse[T any] struct {
	Success   bool   `json:"success"`
	Message   string `json:"message,omitempty"`   // available on any response; hidden when empty
	Data      T      `json:"data,omitempty"`      // omitted on error
	ErrorCode string `json:"errorCode,omitempty"` // present on every error
}

// IsSuccess / IsError — the checker for callers (e.g. service-to-service).
func (r APIResponse[T]) IsSuccess() bool { return r.Success && r.ErrorCode == "" }
func (r APIResponse[T]) IsError() bool   { return !r.IsSuccess() }

func write(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

// OK writes a 200 success envelope.
func OK[T any](w http.ResponseWriter, data T) {
	write(w, http.StatusOK, APIResponse[T]{Success: true, Data: data})
}

// OKMessage writes a 200 success envelope with a message.
func OKMessage[T any](w http.ResponseWriter, message string, data T) {
	write(w, http.StatusOK, APIResponse[T]{Success: true, Message: message, Data: data})
}

// Created writes a 201 success envelope.
func Created[T any](w http.ResponseWriter, data T) {
	write(w, http.StatusCreated, APIResponse[T]{Success: true, Data: data})
}

// Fail writes an error envelope; the error flows through the same APIResponse.
func Fail(w http.ResponseWriter, err error) {
	ae := From(err)
	write(w, ae.HTTPStatus, APIResponse[any]{
		Success:   false,
		Message:   ae.Message,
		ErrorCode: string(ae.Code),
	})
}
