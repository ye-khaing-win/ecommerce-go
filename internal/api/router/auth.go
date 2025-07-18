package router

import (
	"ecommerce-go/internal/api/handlers"
	"ecommerce-go/internal/app"
	"net/http"
)

func RegisterAuthRoutes(mux *http.ServeMux, app *app.Application) {
	h := handlers.AuthHandler{App: app}

	mux.HandleFunc("POST /auth/login", h.Login)
	mux.HandleFunc("POST /auth/logout", h.Logout)
}
