package model

import "github.com/prawirdani/go-midtrans-example/pkg/errors"

type PaymentRequest struct {
	TransactionID string `json:"transactionId"`
}

func (p *PaymentRequest) Validate() error {
	if p.TransactionID == "" {
		return errors.BadRequest("transactionId is required")
	}
	return nil
}

type PaymentResult struct {
	Token       string `json:"token"`
	RedirectURL string `json:"redirectUrl"`
}
