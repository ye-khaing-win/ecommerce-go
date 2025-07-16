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

type CategoryHandler struct {
	app *app.Application
}

func NewCategoryHandler(app *app.Application) *CategoryHandler {
	return &CategoryHandler{app: app}
}

func (h *CategoryHandler) ListCategories(w http.ResponseWriter, r *http.Request) {

	repo := repos.NewCategoryRepository(h.app.Db)
	cats, err := repo.List(r.Context())

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	res := struct {
		Status string            `json:"status,omitempty"`
		Count  int               `json:"count,omitempty"`
		Data   []models.Category `json:"data"`
	}{
		Status: "success",
		Count:  len(cats),
		Data:   cats,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

}
func (h *CategoryHandler) GetCategory(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid ID")
		return
	}
	repo := repos.NewCategoryRepository(h.app.Db)
	cat, err := repo.Get(id)
	switch {
	case errors.Is(err, repos.ErrCategoryNotFound):
		utils.WriteError(w, http.StatusNotFound, err.Error())
		return
	case err != nil:
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	res := struct {
		Status string          `json:"status,omitempty"`
		Data   models.Category `json:"data,omitempty"`
	}{
		Status: "success",
		Data:   cat,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

}
func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var cat models.Category

	if err := utils.DecodeBody(r, &cat); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validator.Validate(&cat); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	repo := repos.NewCategoryRepository(h.app.Db)
	cat, err := repo.Create(cat)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	res := struct {
		Status string          `json:"status,omitempty"`
		Data   models.Category `json:"data,omitempty"`
	}{
		Status: "success",
		Data:   cat,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

}
func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	var cat models.Category
	if err = utils.DecodeBody(r, &cat); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	repo := repos.NewCategoryRepository(h.app.Db)
	cat, err = repo.Update(id, cat)
	if err != nil {
		if errors.Is(err, repos.ErrCategoryNotFound) {
			utils.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	res := struct {
		Status string          `json:"status,omitempty"`
		Data   models.Category `json:"data,omitempty"`
	}{
		Status: "success",
		Data:   cat,
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(res); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

}
func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	repo := repos.NewCategoryRepository(h.app.Db)
	if err := repo.Delete(id); err != nil {
		if errors.Is(err, repos.ErrCategoryNotFound) {
			utils.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
