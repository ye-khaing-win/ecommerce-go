package handlers

import (
	"ecommerce-go/internal/api/dto"
	"ecommerce-go/internal/app"
	"ecommerce-go/internal/repos"
	"ecommerce-go/internal/validator"
	"ecommerce-go/pkg/utils"
	"encoding/json"
	"net/http"
	"os"
	"time"
)

type AuthHandler struct {
	App *app.Application
}

func NewAuthHandler(app *app.Application) *AuthHandler {
	return &AuthHandler{
		App: app,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest

	if err := utils.DecodeBody(r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validator.Validate(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	repo := repos.NewAdminRepository(h.App.Db)

	admin, err := repo.GetByEmail(req.Email)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, "incorrect credentials")
		return

	}
	if !admin.Active {
		utils.WriteError(w, http.StatusUnauthorized, "incorrect credentials")
		return
	}

	// VERIFY PASSWORD
	match, err := utils.VerifyPassword(req.Password, admin.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "authentication failed")
		return
	}

	if !match {
		utils.WriteError(w, http.StatusUnauthorized, "incorrect credentials")
		return
	}

	claims := utils.CustomClaims{
		UserID: admin.ID,
		Email:  admin.Email,
		Role:   admin.Role,
	}

	secret := os.Getenv("JWT_SECRET")
	expires := os.Getenv("JWT_EXPIRES")
	duration, err := time.ParseDuration(expires)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	token, err := utils.SignJWT(claims, secret, duration)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "error creating jwt token")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(duration),
		SameSite: http.SameSiteStrictMode,
	})

	res := struct {
		Status  string `json:"status,omitempty"`
		Message string `json:"message,omitempty"`
	}{
		Status:  "success",
		Message: "login success",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

}
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Unix(0, 0),
		SameSite: http.SameSiteStrictMode,
	})

	res := struct {
		Status  string `json:"status,omitempty"`
		Message string `json:"message,omitempty"`
	}{
		Status:  "success",
		Message: "logout success",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
}
