package service

import (
	"context"

	"github.com/prawirdani/go-midtrans-example/internal/entity"
	"github.com/prawirdani/go-midtrans-example/internal/model"
	"github.com/prawirdani/go-midtrans-example/internal/repository"
)

type TransactionService interface {
	CreateTransaction(
		ctx context.Context,
		req model.TransactionRequest,
	) (*entity.Transaction, error)
	ListTrasaction(ctx context.Context) ([]entity.Transaction, error)
	FindTransaction(ctx context.Context, id string) (*entity.Transaction, error)
	CompleteTransaction(ctx context.Context, id string) error
}

type transactionService struct {
	productRepo     repository.ProductRepository
	userRepo        repository.UserRepository
	transactionRepo repository.TransactionRepository
}

func NewTransactionService(
	productRepo repository.ProductRepository,
	userRepo repository.UserRepository,
	transactionRepo repository.TransactionRepository,
) TransactionService {
	return &transactionService{
		productRepo:     productRepo,
		userRepo:        userRepo,
		transactionRepo: transactionRepo,
	}
}

func (s *transactionService) CreateTransaction(
	ctx context.Context,
	req model.TransactionRequest,
) (*entity.Transaction, error) {
	// Get user
	user, err := s.userRepo.GetUserByID(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	// TODO: Should batch select
	products := make([]entity.Product, len(req.Details))
	for _, detail := range req.Details {
		product, err := s.productRepo.GetProductByID(ctx, detail.ProductID)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	// Create transaction details, by matching the request product id with products
	details := make([]entity.TransactionDetails, len(req.Details))
	for i, detail := range req.Details {
		for _, product := range products {
			if product.ID == detail.ProductID {
				details[i] = entity.TransactionDetails{
					Product:  product,
					Quantity: detail.Quantity,
				}
				details[i].CalculateSubtotal()
				break
			}
		}
	}

	// Create transaction
	transaction := entity.NewTransaction(user, details)

	return s.transactionRepo.Insert(ctx, transaction)
}

func (s *transactionService) ListTrasaction(ctx context.Context) ([]entity.Transaction, error) {
	return s.transactionRepo.Select(ctx)
}

func (s *transactionService) FindTransaction(
	ctx context.Context,
	id string,
) (*entity.Transaction, error) {
	return s.transactionRepo.SelectByID(ctx, id)
}

func (s *transactionService) CompleteTransaction(ctx context.Context, id string) error {
	t, err := s.transactionRepo.SelectByID(ctx, id)
	if err != nil {
		return err
	}

	t.Status = entity.TransactionStatusCompleted
	return s.transactionRepo.SaveChanges(ctx, *t)
}
