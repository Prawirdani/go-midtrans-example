package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/prawirdani/go-midtrans-example/pkg/http/mux"
	"github.com/prawirdani/go-midtrans-example/pkg/http/response"
)

func main() {
	mux := mux.New()

	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_ = response.Send(w, response.WithMessage("Hello, World!"))
	})

	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
