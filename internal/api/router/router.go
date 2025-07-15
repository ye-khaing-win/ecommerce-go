package router

import (
	"ecommerce-go/internal/app"
	"net/http"
)

func Router(app *app.Application) *http.ServeMux {
	mux := http.NewServeMux()

	RegisterCategoryRoutes(mux, app)
	RegisterItemRoutes(mux, app)

	return mux
}
