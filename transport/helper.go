package transport

import (
	"encoding/json"
	"errors"
	"mampu/model"
	"mampu/usecase"
	"net/http"
)

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, model.CommonResponse{
		Message: message,
	})
}

func usecaseErrors(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, usecase.ErrWalletNotFound):
		writeError(w, http.StatusNotFound, err.Error())

	case errors.Is(err, usecase.ErrInsufficientBalance):
		writeError(w, http.StatusBadRequest, err.Error())

	case errors.Is(err, usecase.ErrInvalidAmount):
		writeError(w, http.StatusBadRequest, err.Error())

	default:
		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}
