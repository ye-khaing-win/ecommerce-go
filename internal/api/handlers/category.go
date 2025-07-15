package handlers

import (
	"ecommerce-go/internal/app"
	"ecommerce-go/internal/models"
	"ecommerce-go/internal/repos"
	"ecommerce-go/internal/validator"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
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
		writeError(w, http.StatusInternalServerError, err.Error())
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
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

}
func (h *CategoryHandler) GetCategory(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid ID")
		return
	}
	repo := repos.NewCategoryRepository(h.app.Db)
	cat, err := repo.Get(id)
	switch {
	case errors.Is(err, repos.ErrCategoryNotFound):
		writeError(w, http.StatusNotFound, "Category not found")
		return
	case err != nil:
		writeError(w, http.StatusInternalServerError, err.Error())
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
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

}
func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var cat models.Category
	if err := json.NewDecoder(r.Body).Decode(&cat); err != nil {
		var ute *json.UnmarshalTypeError
		if errors.As(err, &ute) {
			msg := fmt.Sprintf("%s must be %s", strings.ToLower(ute.Field), ute.Type.Kind())
			writeError(w, http.StatusBadRequest, msg)
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validator.Validate(&cat); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	repo := repos.NewCategoryRepository(h.app.Db)
	cat, err := repo.Create(cat)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
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
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

}
func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	var cat models.Category
	if err := json.NewDecoder(r.Body).Decode(&cat); err != nil {
		var ute *json.UnmarshalTypeError
		if errors.As(err, &ute) {
			msg := fmt.Sprintf("%s must be %s", strings.ToLower(ute.Field), ute.Type.Kind())
			writeError(w, http.StatusBadRequest, msg)
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	repo := repos.NewCategoryRepository(h.app.Db)
	cat, err = repo.Update(id, cat)
	if err != nil {
		if errors.Is(err, repos.ErrCategoryNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
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
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

}
func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	repo := repos.NewCategoryRepository(h.app.Db)
	if err := repo.Delete(id); err != nil {
		if errors.Is(err, repos.ErrCategoryNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	res := struct {
		Status  string `json:"status,omitempty"`
		Message string `json:"message,omitempty"`
	}{
		Status:  "error",
		Message: msg,
	}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
}
