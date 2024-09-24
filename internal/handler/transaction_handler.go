package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/prawirdani/go-midtrans-example/internal/model"
	"github.com/prawirdani/go-midtrans-example/internal/service"
	req "github.com/prawirdani/go-midtrans-example/pkg/http/request"
	res "github.com/prawirdani/go-midtrans-example/pkg/http/response"
)

type TransactionHandler struct {
	trxService service.TransactionService
}

func NewTransactionHandler(
	trxService service.TransactionService,
	paymentService service.PaymentService,
) *TransactionHandler {
	return &TransactionHandler{
		trxService: trxService,
	}
}

func (h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var reqBody model.TransactionRequest
	if err := req.BindValidate(r, &reqBody); err != nil {
		res.HandleError(w, err)
		return
	}

	trx, err := h.trxService.CreateTransaction(r.Context(), reqBody)
	if err != nil {
		res.HandleError(w, err)
		return
	}

	_ = res.Send(w, res.WithStatus(201), res.WithMessage("Transaction created"), res.WithData(trx))
}

func (h *TransactionHandler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	trxs, err := h.trxService.ListTrasaction(r.Context())
	if err != nil {
		res.HandleError(w, err)
		return
	}

	_ = res.Send(w, res.WithData(trxs))
}

func (h *TransactionHandler) GetById(w http.ResponseWriter, r *http.Request) {
	transactionId := chi.URLParam(r, "id")
	tx, err := h.trxService.FindTransaction(r.Context(), transactionId)
	if err != nil {
		res.HandleError(w, err)
		return
	}

	_ = res.Send(w, res.WithData(tx))
}

// func (h *TransactionHandler) ProcessPayment(w http.ResponseWriter, r *http.Request) {
// 	transactionId, err := strconv.Atoi(chi.URLParam(r, "id"))
// 	if err != nil {
// 		res.HandleError(w, err)
// 		return
// 	}
//
// 	d, err := h.paymentService.Pay(r.Context(), transactionId)
// 	if err != nil {
// 		res.HandleError(w, err)
// 		return
// 	}
//
// 	_ = res.Send(w, res.WithData(d), res.WithMessage("Checkout success"))
// }
