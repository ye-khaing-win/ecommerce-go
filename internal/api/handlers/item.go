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
	App            *app.Application
	AllowedFilters map[string]struct{}
	AllowedSorts   map[string]struct{}
}

func NewItemHandler(app *app.Application) *ItemHandler {
	return &ItemHandler{
		App: app,
		AllowedFilters: map[string]struct{}{
			"name":        {},
			"description": {},
			"category_id": {},
		},
		AllowedSorts: map[string]struct{}{
			"name":       {},
			"price":      {},
			"created_at": {},
		},
	}
}

func (h *ItemHandler) ListItems(w http.ResponseWriter, r *http.Request) {
	repo := repos.NewItemRepository(h.App.Db)
	items, err := repo.List(r.Context())

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	res := struct {
		Status string        `json:"status,omitempty"`
		Count  int           `json:"count,omitempty"`
		Data   []models.Item `json:"data"`
	}{
		Status: "success",
		Count:  len(items),
		Data:   items,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
}
func (h *ItemHandler) GetItems(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}

	repo := repos.NewItemRepository(h.App.Db)
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

	repo := repos.NewItemRepository(h.App.Db)
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
func (h *ItemHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var item models.Item
	if err := utils.DecodeBody(r, &item); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	repo := repos.NewItemRepository(h.App.Db)
	item, err = repo.Update(id, item)
	if err != nil {
		if errors.Is(err, repos.ErrItemNotFound) {
			utils.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
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
	if err = json.NewEncoder(w).Encode(res); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
}
func (h *ItemHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}

	repo := repos.NewItemRepository(h.App.Db)

	if err = repo.Delete(id); err != nil {
		if errors.Is(err, repos.ErrItemNotFound) {
			utils.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
func (h *ItemHandler) ListItemByCategory(w http.ResponseWriter, r *http.Request) {
	catIDStr := r.PathValue("id")
	catID, err := strconv.Atoi(catIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}

	repo := repos.NewItemRepository(h.App.Db)
	items, err := repo.ListByCategory(r.Context(), catID)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	res := struct {
		Status string        `json:"status,omitempty"`
		Count  int           `json:"count,omitempty"`
		Data   []models.Item `json:"data"`
	}{
		Status: "success",
		Count:  len(items),
		Data:   items,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
}
