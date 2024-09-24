package errors

import (
	"log"
	"net/http"
)

var (
	BadRequest       = build(http.StatusBadRequest)
	Conflict         = build(http.StatusConflict)
	NotFound         = build(http.StatusNotFound)
	Unauthorized     = build(http.StatusUnauthorized)
	Unprocessable    = build(http.StatusUnprocessableEntity)
	MethodNotAllowed = build(http.StatusMethodNotAllowed)
	Forbidden        = build(http.StatusForbidden)
	InternalServer   = build(http.StatusInternalServerError)
)

type ApiError struct {
	Status  int
	Message string
	Cause   interface{}
}

func (e *ApiError) Error() string {
	return e.Message
}

func Parse(err error) *ApiError {
	switch e := err.(type) {
	case *ApiError:
		return e
	default:
		log.Println("An unexpected error occurred:", err)
		return &ApiError{
			Status:  500,
			Message: "An unexpected error occurred, try again latter",
		}
	}
}

func build(status int) func(msg string) *ApiError {
	return func(m string) *ApiError {
		return &ApiError{
			Status:  status,
			Message: m,
		}
	}
}
