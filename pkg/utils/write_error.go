package utils

import (
	"encoding/json"
	"net/http"
)

func WriteError(w http.ResponseWriter, status int, msg string) {
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
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
}
