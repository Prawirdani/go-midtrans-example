package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/prawirdani/go-midtrans-example/db"
	"github.com/prawirdani/go-midtrans-example/internal/handler"
	"github.com/prawirdani/go-midtrans-example/internal/repository"
	"github.com/prawirdani/go-midtrans-example/internal/service"
	"github.com/prawirdani/go-midtrans-example/pkg/http/response"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
}

func main() {
	PORT := os.Getenv("PORT")

	conn := db.Connection()
	defer conn.Close()
	db.Init(conn)

	// Setup Repositories
	userRepo := repository.NewUserRepository(conn)
	productRepo := repository.NewProductRepository(conn)
	productHandler := handler.NewProductHandler(productRepo)
	transactionRepo := repository.NewTransactionRepository(conn, productRepo)
	// Setup Services
	paymentService := service.NewPaymentService(os.Getenv("MIDTRANS_SERVER_KEY"), transactionRepo)
	transactionService := service.NewTransactionService(
		productRepo,
		userRepo,
		transactionRepo,
	)
	// Setup Handlers
	transactionHandler := handler.NewTransactionHandler(transactionService, paymentService)
	paymentHandler := handler.NewPaymentHandler(paymentService, transactionService)

	// Setup routes
	mux := chi.NewRouter()

	mux.Use(middleware.Logger)
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_ = response.Send(w, response.WithMessage("Hello, World!"))
	})
	mux.Get("/products", productHandler.ProductList)
	mux.Get("/transactions", transactionHandler.GetTransactions)
	mux.Get("/transactions/{id}", transactionHandler.GetById)
	mux.Post("/transactions", transactionHandler.CreateTransaction)
	mux.Post("/payments", paymentHandler.ProcessPayment)
	mux.Post("/payments/callback", paymentHandler.PaymentCallback)

	users, err := userRepo.GetUsers(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Use this user data to create transaction:", users)

	log.Println("Server running on port", PORT)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", PORT), mux))
}
