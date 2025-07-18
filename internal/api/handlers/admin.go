package handlers

import (
	"ecommerce-go/internal/app"
	"ecommerce-go/internal/models"
	"ecommerce-go/internal/repos"
	"ecommerce-go/internal/validator"
	"ecommerce-go/pkg/utils"
	"encoding/json"
	"net/http"
)

type AdminHandler struct {
	App            *app.Application
	AllowedFilters map[string]struct{}
	AllowedSorts   map[string]struct{}
}

func NewAdminHandler(app *app.Application) *AdminHandler {
	return &AdminHandler{
		App: app,
		AllowedFilters: map[string]struct{}{
			"role": {},
		},
		AllowedSorts: map[string]struct{}{
			"first_name": {},
			"last_name":  {},
			"email":      {},
			"created_at": {},
		},
	}
}

func (h *AdminHandler) ListAdmin(w http.ResponseWriter, r *http.Request) {
	repo := repos.NewAdminRepository(h.App.Db)
	admins, err := repo.List(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	res := struct {
		Status string         `json:"status,omitempty"`
		Count  int            `json:"count,omitempty"`
		Data   []models.Admin `json:"data"`
	}{
		Status: "success",
		Count:  len(admins),
		Data:   admins,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *AdminHandler) CreateAdmin(w http.ResponseWriter, r *http.Request) {
	var admin models.Admin

	if err := utils.DecodeBody(r, &admin); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validator.Validate(&admin); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	repo := repos.NewAdminRepository(h.App.Db)
	admin, err := repo.Create(admin)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	res := struct {
		Status string       `json:"status,omitempty"`
		Data   models.Admin `json:"data,omitempty"`
	}{
		Status: "success",
		Data:   admin,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
}
