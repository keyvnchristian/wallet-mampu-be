package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type (
	TransactionRepository interface {
		Insert(ctx context.Context, tx *sqlx.Tx, userID int, amount int64) error
	}

	transactionRepository struct {
		db *sqlx.DB
	}
)

func NewTransactionRepository(db *sqlx.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Insert(ctx context.Context, tx *sqlx.Tx, userID int, amount int64) error {
	query, args, err := sq.
		Insert("transactions").
		Columns("user_id", "amount", "type").
		Values(userID, amount, "withdraw").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, query, args...)
	return err
}
