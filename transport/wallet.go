package transport

import (
	"context"
	"encoding/json"
	"mampu/model"
	"mampu/usecase"
	"net/http"
	"strconv"
	"time"
)

type (
	WalletHandler struct {
		uc usecase.WalletUsecase
	}
)

func NewWalletHandler(uc usecase.WalletUsecase) *WalletHandler {
	return &WalletHandler{uc: uc}
}

func (h *WalletHandler) Withdraw(w http.ResponseWriter, r *http.Request) {
	var req model.WithdrawRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.UserID <= 0 || req.Amount <= 0 {
		writeError(w, http.StatusBadRequest, "invalid input")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	balance, err := h.uc.Withdraw(ctx, req.UserID, req.Amount)
	if err != nil {
		usecaseErrors(w, err)
		return
	}

	resp := model.WithdrawResponse{
		RemainingBalance: balance.Balance,
	}

	writeJSON(w, http.StatusOK, model.CommonResponse{
		Message: "Successfully Withdraw",
		Data:    resp,
	})
}

func (h *WalletHandler) GetWallet(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		writeError(w, http.StatusBadRequest, "user_id required")
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID <= 0 {
		writeError(w, http.StatusBadRequest, "invalid user_id")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	balance, err := h.uc.GetWallet(ctx, userID)
	if err != nil {
		usecaseErrors(w, err)
		return
	}

	writeJSON(w, http.StatusOK, model.CommonResponse{
		Message: "Successfully Retrieved Wallet",
		Data:    balance.ID,
	})
}
