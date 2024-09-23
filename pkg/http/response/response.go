package response

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/prawirdani/go-midtrans-example/pkg/errors"
)

type ResponseBody struct {
	Data    any        `json:"data,omitempty"`
	Message string     `json:"message,omitempty"`
	Status  int        `json:"-"`
	Error   *ErrorBody `json:"error"`
}

type ErrorBody struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details"`
}

// Send is a function to send json response to client
// It uses option pattern to accepts multiple options to customize the response
func Send(w http.ResponseWriter, opts ...Option) error {
	res := ResponseBody{
		Status: http.StatusOK, // Default
	}
	for _, opt := range opts {
		opt(&res)
	}
	return writeJSON(w, res.Status, res)
}

// Response writer for handling error
func HandleError(w http.ResponseWriter, err error) {
	e := errors.Parse(err)
	response := ResponseBody{
		Error: &ErrorBody{
			Code:    e.Status,
			Message: e.Message,
			Details: e.Cause,
		},
	}

	writeErr := writeJSON(w, e.Status, response)
	if writeErr != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Utility function to help writing json to response body.
func writeJSON(w http.ResponseWriter, status int, response interface{}) error {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)

	if err := enc.Encode(response); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err := w.Write(buf.Bytes())
	return err
}
