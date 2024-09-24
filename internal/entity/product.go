package entity

import "github.com/prawirdani/go-midtrans-example/pkg/errors"

var ErrProductNotFound = errors.NotFound("product not found")

type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}
