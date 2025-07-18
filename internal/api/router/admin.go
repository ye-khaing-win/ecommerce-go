package router

import (
	"ecommerce-go/internal/api/handlers"
	mw "ecommerce-go/internal/api/middlewares"
	"ecommerce-go/internal/app"
	"net/http"
)

func RegisterAdminRoutes(mux *http.ServeMux, app *app.Application) {
	h := handlers.NewAdminHandler(app)

	mux.Handle("GET /admins", mw.Filter(h.AllowedFilters)(mw.Sort(h.AllowedSorts)(http.HandlerFunc(h.ListAdmin))))
	mux.HandleFunc("POST /admins", h.CreateAdmin)
}
