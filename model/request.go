package model

type (
	WithdrawRequest struct {
		UserID int   `json:"user_id"`
		Amount int64 `json:"amount"`
	}
)
