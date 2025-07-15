package router

import (
	"ecommerce-go/internal/api/handlers"
	mw "ecommerce-go/internal/api/middlewares"
	"ecommerce-go/internal/app"
	"net/http"
)

func Router(app *app.Application) *http.ServeMux {
	mux := http.NewServeMux()

	h := handlers.NewCategoryHandler(app)

	mux.Handle("GET /categories", mw.Select(mw.Filter(mw.Sort(http.HandlerFunc(h.ListCategories)))))
	mux.HandleFunc("GET /categories/{id}", h.GetCategory)
	mux.HandleFunc("POST /categories", h.CreateCategory)
	mux.HandleFunc("PATCH /categories/{id}", h.UpdateCategory)
	mux.HandleFunc("DELETE /categories/{id}", h.DeleteCategory)

	return mux
}
