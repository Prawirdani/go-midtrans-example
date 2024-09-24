package model

import (
	"fmt"
	"net/http"

	"github.com/prawirdani/go-midtrans-example/pkg/errors"
)

type TransactionDetailRequest struct {
	ProductID int `json:"productId"`
	Quantity  int `json:"quantity"`
}

type TransactionRequest struct {
	UserID  string                     `json:"userId"`
	Details []TransactionDetailRequest `json:"details"`
}

func (tr *TransactionRequest) Validate() error {
	errMap := make(map[string]string)

	if tr.UserID == "" {
		errMap["userId"] = "userId is required"
	}

	if len(tr.Details) == 0 {
		errMap["details"] = "details is required"
	}

	for i, detail := range tr.Details {
		if detail.ProductID == 0 {
			errMap[fmt.Sprintf("details[%d]", i)] = "productId is required"
		}
		if detail.Quantity == 0 {
			errMap[fmt.Sprintf("details[%d]", i)] = "quantity is required"
		}
	}

	if len(errMap) > 0 {
		err := errors.ApiError{
			Status:  http.StatusUnprocessableEntity,
			Message: "Invalid request",
			Cause:   errMap,
		}

		return &err
	}

	return nil
}
