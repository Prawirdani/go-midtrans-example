package service

import (
	"context"
	"fmt"
	"log"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/prawirdani/go-midtrans-example/internal/model"
	"github.com/prawirdani/go-midtrans-example/pkg/errors"
)

type PaymentService interface {
	Process(ctx context.Context, transactionID string) (*model.PaymentResult, error)
	VerifyPayment(
		ctx context.Context,
		paymentOrderID string,
	) (bool, error)
}

type paymentService struct {
	transactionService TransactionService
	snapClient         snap.Client
	coreClient         coreapi.Client
}

func NewPaymentService(
	serverKey string,
	ts TransactionService,
) PaymentService {
	snapClient := snap.Client{}
	snapClient.New(serverKey, midtrans.Sandbox)

	coreClient := coreapi.Client{}
	coreClient.New(serverKey, midtrans.Sandbox)

	return &paymentService{
		transactionService: ts,
		snapClient:         snapClient,
		coreClient:         coreClient,
	}
}

func (s *paymentService) Process(
	ctx context.Context,
	transactionID string,
) (*model.PaymentResult, error) {
	transaction, err := s.transactionService.FindTransaction(ctx, transactionID)
	if err != nil {
		return nil, err
	}

	details := make([]midtrans.ItemDetails, 0)
	for i := range transaction.Details {
		detail := transaction.Details[i]
		details = append(details, midtrans.ItemDetails{
			ID:    fmt.Sprintf("Product-%d", detail.Product.ID),
			Price: int64(detail.Product.Price),
			Qty:   int32(detail.Quantity),
			Name:  detail.Product.Name,
		})
	}

	snapReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  transaction.ID.String(),
			GrossAmt: int64(transaction.Total),
		},
		Items: &details,
		EnabledPayments: []snap.SnapPaymentType{
			snap.PaymentTypeGopay,
			snap.PaymentTypeBankTransfer,
			snap.PaymentTypeShopeepay,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: transaction.User.FirstName,
			LName: transaction.User.LastName,
			Email: transaction.User.Email,
			Phone: transaction.User.Phone,
		},
		Expiry: &snap.ExpiryDetails{
			Duration: int64(10),
			Unit:     "minute",
		},
	}

	res, e := s.snapClient.CreateTransaction(snapReq)
	if e != nil {
		log.Println("Failed to create transaction", err)
		return nil, errors.BadRequest(e.GetMessage())
	}

	coResult := model.PaymentResult{
		Token:       res.Token,
		RedirectURL: res.RedirectURL,
	}
	return &coResult, nil
}

func (s *paymentService) VerifyPayment(
	ctx context.Context,
	paymentOrderID string,
) (bool, error) {
	transactionStatusResp, e := s.coreClient.CheckTransaction(paymentOrderID)
	if e != nil {
		return false, errors.InternalServer(e.GetMessage())
	}

	if transactionStatusResp != nil {
		switch transactionStatusResp.TransactionStatus {
		case "capture":
			if transactionStatusResp.FraudStatus == "challenge" {
				// TODO: set transaction status on your database to 'challenge'
			} else if transactionStatusResp.FraudStatus == "accept" {
				return true, nil
				// TODO: set transaction status on your database to 'success'
			}
		case "settlement":
			return true, nil
			// TODO: set transaction status on your databaase to 'success'
		case "deny":
			// TODO: you can ignore 'deny', because most of the time it allows payment retries
		case "cancel", "expire":
			// TODO: set transaction status on your databaase to 'failure'
		case "pending":
			// TODO: set transaction status on your databaase to 'pending' / waiting payment
		}
	}

	return false, nil
}
