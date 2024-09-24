package handler

import (
	"encoding/json"
	"net/http"

	"github.com/prawirdani/go-midtrans-example/internal/model"
	"github.com/prawirdani/go-midtrans-example/internal/service"
	"github.com/prawirdani/go-midtrans-example/pkg/errors"
	req "github.com/prawirdani/go-midtrans-example/pkg/http/request"
	res "github.com/prawirdani/go-midtrans-example/pkg/http/response"
)

type PaymentHandler struct {
	transactionService service.TransactionService
	paymentService     service.PaymentService
}

func NewPaymentHandler(ps service.PaymentService, ts service.TransactionService) *PaymentHandler {
	return &PaymentHandler{
		paymentService:     ps,
		transactionService: ts,
	}
}

func (h *PaymentHandler) ProcessPayment(w http.ResponseWriter, r *http.Request) {
	var reqBody model.PaymentRequest
	if err := req.BindValidate(r, &reqBody); err != nil {
		res.HandleError(w, err)
		return
	}
	paymentResult, err := h.paymentService.Process(r.Context(), reqBody.TransactionID)
	if err != nil {
		res.HandleError(w, err)
		return
	}
	_ = res.Send(w, res.WithData(paymentResult))
}

// PaymentCallback is a handler to handle payment callback from midtrans
// This handler will be called by midtrans after payment
func (h *PaymentHandler) PaymentCallback(w http.ResponseWriter, r *http.Request) {
	// 1. Initialize empty map
	var notificationPayload map[string]interface{}

	// 2. Parse JSON request body and use it to set json to payload
	err := json.NewDecoder(r.Body).Decode(&notificationPayload)
	if err != nil {
		res.HandleError(w, errors.BadRequest("invalid request body"))
		return
	}
	// 3. Get order-id from payload
	orderId, exists := notificationPayload["order_id"].(string)
	if !exists {
		res.HandleError(w, errors.NotFound("order_id not found"))
		return
	}

	// 4. Verify payment
	success, _ := h.paymentService.VerifyPayment(r.Context(), orderId)

	if success {
		// 5. Complete Transaction
		err = h.transactionService.CompleteTransaction(r.Context(), orderId)
		if err != nil {
			res.HandleError(w, err)
			return
		}
	}

	_ = res.Send(w, res.WithStatus(http.StatusOK))
}
