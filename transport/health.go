package transport

import (
	"context"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
)

type HealthHandler struct {
	db *sqlx.DB
}

func NewHealthHandler(db *sqlx.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

type HealthResponse struct {
	Status string `json:"status"`
	DB     string `json:"database"`
}

func (h *HealthHandler) Ping(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	err := h.db.PingContext(ctx)

	resp := HealthResponse{
		Status: "ok",
		DB:     "up",
	}

	if err != nil {
		resp.Status = "not ok =("
		resp.DB = "down"
		writeJSON(w, http.StatusServiceUnavailable, resp)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}
