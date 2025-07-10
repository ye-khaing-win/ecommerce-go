package router

import (
	"ecommerce-go/internal/api/handlers"
	"net/http"
)

func Router() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /categories", handlers.ListCategories)
	mux.HandleFunc("GET /categories/{id}", handlers.GetCategories)
	mux.HandleFunc("POST /categories", handlers.CreateCategories)
	mux.HandleFunc("PATCH /categories/{id}", handlers.UpdateCategories)
	mux.HandleFunc("DELETE /categories/{id}", handlers.DeleteCategories)

	return mux
}
