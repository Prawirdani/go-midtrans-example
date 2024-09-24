package handler

import (
	"net/http"

	"github.com/prawirdani/go-midtrans-example/internal/repository"
	res "github.com/prawirdani/go-midtrans-example/pkg/http/response"
)

type ProductHandler struct {
	repo repository.ProductRepository
}

func NewProductHandler(repo repository.ProductRepository) *ProductHandler {
	return &ProductHandler{repo: repo}
}

func (h *ProductHandler) ProductList(w http.ResponseWriter, r *http.Request) {
	products, err := h.repo.GetProducts(r.Context())
	if err != nil {
		res.HandleError(w, err)
		return
	}
	_ = res.Send(w, res.WithData(products))
}
