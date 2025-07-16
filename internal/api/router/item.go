package router

import (
	"ecommerce-go/internal/api/handlers"
	"ecommerce-go/internal/app"
	"net/http"
)

func RegisterItemRoutes(mux *http.ServeMux, app *app.Application) {
	h := handlers.NewItemHandler(app)

	mux.HandleFunc("GET /items", h.ListItems)
	mux.HandleFunc("GET /items/{id}", h.Get)
	mux.HandleFunc("POST /items", h.CreateItem)
}
