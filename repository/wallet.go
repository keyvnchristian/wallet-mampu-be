package repository

import (
	"context"
	"database/sql"
	"errors"
	"mampu/model"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type WalletRepository interface {
	GetWallet(ctx context.Context, userID int) (*model.Wallet, error)
	GetWalletWithLock(ctx context.Context, tx *sqlx.Tx, userID int) (*model.Wallet, error)
	UpdateBalance(ctx context.Context, tx *sqlx.Tx, userID int, newBalance int64) error
}

type walletRepository struct {
	db *sqlx.DB
}

func NewWalletRepository(db *sqlx.DB) WalletRepository {
	return &walletRepository{db: db}
}

func (r *walletRepository) GetWallet(ctx context.Context, userID int) (*model.Wallet, error) {
	query, args, err := sq.
		Select("user_id", "balance").
		From("wallets").
		Where(sq.Eq{"user_id": userID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var wallet model.Wallet
	err = r.db.GetContext(ctx, &wallet, r.db.Rebind(query), args...)
	if err == sql.ErrNoRows {
		return nil, errors.New("wallet not found")
	}
	if err != nil {
		return nil, err
	}

	return &wallet, nil
}

func (r *walletRepository) GetWalletWithLock(ctx context.Context, tx *sqlx.Tx, userID int) (*model.Wallet, error) {
	query, args, err := sq.
		Select("user_id", "balance").
		From("wallets").
		Where(sq.Eq{"user_id": userID}).
		Suffix("FOR UPDATE").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	var wallet model.Wallet
	if tx != nil {
		err = tx.GetContext(ctx, &wallet, query, args...)
	} else {
		err = r.db.GetContext(ctx, &wallet, r.db.Rebind(query), args...)
	}
	if err == sql.ErrNoRows {
		return nil, errors.New("wallet not found")
	}
	if err != nil {
		return nil, err
	}

	return &wallet, nil
}

func (r *walletRepository) UpdateBalance(ctx context.Context, tx *sqlx.Tx, userID int, newBalance int64) error {
	query, args, err := sq.
		Update("wallets").
		Set("balance", newBalance).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.Eq{"user_id": userID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	var execErr error
	if tx != nil {
		_, execErr = tx.ExecContext(ctx, query, args...)
	} else {
		_, execErr = r.db.ExecContext(ctx, r.db.Rebind(query), args...)
	}
	return execErr
}
