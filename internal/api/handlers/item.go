package handlers

import (
	"ecommerce-go/internal/app"
	"ecommerce-go/internal/models"
	"ecommerce-go/internal/repos"
	"ecommerce-go/internal/validator"
	"ecommerce-go/pkg/utils"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type ItemHandler struct {
	app *app.Application
}

func NewItemHandler(app *app.Application) *ItemHandler {
	return &ItemHandler{app: app}
}

func (h *ItemHandler) ListItems(w http.ResponseWriter, r *http.Request) {

	res := struct {
		Status string `json:"status,omitempty"`
		//Count   int    `json:"count,omitempty"`
		Message string `json:"message"`
		//Data   []models.Category `json:"data"`
	}{
		Status: "success",
		//Count:  len(cats),
		//Data:   cats,
		Message: "Good job buddy",
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *ItemHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}

	repo := repos.NewItemRepository(h.app.Db)
	item, err := repo.Get(id)

	switch {
	case errors.Is(err, repos.ErrItemNotFound):
		utils.WriteError(w, http.StatusNotFound, err.Error())
		return
	case err != nil:
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	res := struct {
		Status string      `json:"status,omitempty"`
		Data   models.Item `json:"data,omitempty"`
	}{
		Status: "success",
		Data:   item,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *ItemHandler) CreateItem(w http.ResponseWriter, r *http.Request) {
	var item models.Item

	if err := utils.DecodeBody(r, &item); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validator.Validate(&item); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	repo := repos.NewItemRepository(h.app.Db)
	item, err := repo.Create(item)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	res := struct {
		Status string      `json:"status,omitempty"`
		Data   models.Item `json:"data,omitempty"`
	}{
		Status: "success",
		Data:   item,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

}
