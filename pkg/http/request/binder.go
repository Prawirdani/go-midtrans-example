package request

import (
	"encoding/json"
	"net/http"

	"github.com/prawirdani/go-midtrans-example/pkg/errors"
)

// BindValidate is a helper function to bind and validate json request body, requires a struct that implements RequestBody interface
func BindValidate(r *http.Request, reqBody RequestBody) error {
	if r.Header.Get("Content-Type") != "application/json" {
		return errors.BadRequest("Content-Type must be application/json")
	}

	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(reqBody); err != nil {
		return parseBodyErr(err)
	}

	if err := reqBody.Sanitize(); err != nil {
		return err
	}

	return reqBody.Validate()
}
