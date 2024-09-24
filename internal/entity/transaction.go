package entity

import (
	"github.com/google/uuid"
	"github.com/prawirdani/go-midtrans-example/pkg/errors"
)

var ErrTransactionNotFound = errors.NotFound("transaction not found")

type TransactionDetails struct {
	ID       int     `json:"id"`
	Product  Product `json:"product"`
	Quantity int     `json:"quantity"`
	Subtotal int     `json:"subtotal"`
}

func (td *TransactionDetails) CalculateSubtotal() {
	td.Subtotal = td.Product.Price * td.Quantity
}

type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusCompleted TransactionStatus = "completed"
	TransactionStatusFailed    TransactionStatus = "failed"
	TransactionStatusCancelled TransactionStatus = "cancelled"
)

type Transaction struct {
	ID      uuid.UUID            `json:"id"`
	User    User                 `json:"user"`
	Details []TransactionDetails `json:"details"`
	Total   int                  `json:"total"`
	Status  TransactionStatus    `json:"status"`
}

func (t *Transaction) CalculateTotal() {
	total := 0
	for _, detail := range t.Details {
		total += detail.Subtotal
	}
	t.Total = total
}

func NewTransaction(user User, details []TransactionDetails) Transaction {
	t := Transaction{
		ID:      uuid.New(),
		User:    user,
		Details: details,
	}
	t.CalculateTotal()
	return t
}
