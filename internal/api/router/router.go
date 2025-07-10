package router

import (
	"ecommerce-go/internal/api/handlers"
	"net/http"
)

func Router() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /categories", handlers.ListCategories)
	mux.HandleFunc("GET /categories/{id}", handlers.GetCategory)
	mux.HandleFunc("POST /categories", handlers.CreateCategory)
	mux.HandleFunc("PATCH /categories/{id}", handlers.UpdateCategory)
	mux.HandleFunc("DELETE /categories/{id}", handlers.DeleteCategory)

	return mux
}
