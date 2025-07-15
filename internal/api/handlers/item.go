package handlers

import (
	"ecommerce-go/internal/app"
	"encoding/json"
	"net/http"
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
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
}
