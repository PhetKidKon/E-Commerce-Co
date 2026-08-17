package response

import (
	"errors"
	"net/http"
)

// Code is a stable, machine-readable error code (frontend switches on this).
type Code string

const (
	CodeBadRequest   Code = "BAD_REQUEST"
	CodeUnauthorized Code = "UNAUTHORIZED"
	CodeForbidden    Code = "FORBIDDEN"
	CodeNotFound     Code = "NOT_FOUND"
	CodeConflict     Code = "CONFLICT"
	CodeValidation   Code = "VALIDATION_ERROR"
	CodeInternal     Code = "INTERNAL_ERROR"
)

// AppError carries an error code, a human message, the HTTP status, and an
// optional wrapped cause. It implements the error interface.
type AppError struct {
	Code       Code
	Message    string
	HTTPStatus int
	Err        error // wrapped cause: for logs, never exposed to the client
}

func (e *AppError) Error() string { return string(e.Code) + ": " + e.Message }
func (e *AppError) Unwrap() error { return e.Err }

// Wrap attaches a cause (for logging) without changing the public message.
func (e *AppError) Wrap(cause error) *AppError { c := *e; c.Err = cause; return &c }

// New builds an AppError.
func New(code Code, status int, msg string) *AppError {
	return &AppError{Code: code, Message: msg, HTTPStatus: status}
}

// Constructors — use these in the service layer.
func BadRequest(m string) *AppError   { return New(CodeBadRequest, http.StatusBadRequest, m) }
func Unauthorized(m string) *AppError { return New(CodeUnauthorized, http.StatusUnauthorized, m) }
func Forbidden(m string) *AppError    { return New(CodeForbidden, http.StatusForbidden, m) }
func NotFound(m string) *AppError     { return New(CodeNotFound, http.StatusNotFound, m) }
func Conflict(m string) *AppError     { return New(CodeConflict, http.StatusConflict, m) }
func Validation(m string) *AppError   { return New(CodeValidation, http.StatusUnprocessableEntity, m) }
func Internal(m string) *AppError     { return New(CodeInternal, http.StatusInternalServerError, m) }

// From converts any error into an *AppError (unknown errors default to Internal
// so internal details never leak to the client, but stay wrapped for logs).
func From(err error) *AppError {
	if err == nil {
		return nil
	}
	var ae *AppError
	if errors.As(err, &ae) {
		return ae
	}
	return Internal("internal server error").Wrap(err)
}
