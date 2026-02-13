package model

import "time"

type (
	Transaction struct {
		ID        int64     `db:"id" json:"id"`
		UserID    int64     `db:"user_id" json:"user_id"`
		Amount    int64     `db:"amount" json:"amount"`
		Type      string    `db:"type" json:"type"`
		CreatedAt time.Time `db:"created_at" json:"created_at"`
	}
)
