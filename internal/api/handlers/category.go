package handlers

import (
	"ecommerce-go/internal/models"
	"ecommerce-go/internal/repositories"
	"ecommerce-go/internal/validator"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func ListCategories(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("List Categories")); err != nil {
		log.Println("Write error: ", err)
	}

}
func GetCategories(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("Get Category")); err != nil {
		log.Println("Write error: ", err)
	}
}
func CreateCategories(w http.ResponseWriter, r *http.Request) {
	var cat models.Category
	err := json.NewDecoder(r.Body).Decode(&cat)
	if err != nil {
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

	cat, err = repositories.CreateCategory(cat)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	res := struct {
		Status string          `json:"status,omitempty"`
		Data   models.Category `json:"data,omitempty"`
	}{
		Status: "success",
		Data:   cat,
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

}
func UpdateCategories(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("Update Category")); err != nil {
		log.Println("Write error: ", err)
	}
}
func DeleteCategories(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("Delete Category")); err != nil {
		log.Println("Write error: ", err)
	}
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
