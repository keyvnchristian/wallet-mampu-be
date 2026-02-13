package model

type (
	CommonResponse struct {
		Message string `json:"message"`
		Data    any    `json:"data,omitempty"`
	}

	WithdrawResponse struct {
		RemainingBalance int64 `json:"remaining_balance"`
	}

	BalanceResponse struct {
		UserID  int   `json:"user_id"`
		Balance int64 `json:"balance"`
	}
)
