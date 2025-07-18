package router

import (
	"ecommerce-go/internal/api/handlers"
	mw "ecommerce-go/internal/api/middlewares"
	"ecommerce-go/internal/app"
	"net/http"
)

func RegisterItemRoutes(mux *http.ServeMux, app *app.Application) {
	h := handlers.NewItemHandler(app)

	mux.Handle("GET /items", mw.Filter(h.AllowedFilters)(mw.Sort(h.AllowedSorts)(http.HandlerFunc(h.ListItems))))
	mux.HandleFunc("GET /items/{id}", h.GetItems)
	mux.HandleFunc("POST /items", h.CreateItem)
	mux.HandleFunc("PATCH /items/{id}", h.UpdateItem)
	mux.HandleFunc("DELETE /items/{id}", h.DeleteItem)
}
