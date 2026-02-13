package model

import "time"

type (
	Wallet struct {
		ID        int64     `db:"id" json:"id"`
		UserID    int64     `db:"user_id" json:"-"`
		Balance   int64     `db:"balance" json:"balance"`
		UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	}
)
