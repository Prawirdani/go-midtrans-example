package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/prawirdani/go-midtrans-example/pkg/http/mux"
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

	mux := mux.New()

	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_ = response.Send(w, response.WithMessage("Hello, World!"))
	})

	fmt.Println("Server running on port", PORT)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", PORT), mux))
}
