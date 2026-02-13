package usecase

import (
	"context"
	"mampu/model"
	"mampu/repository"

	"github.com/jmoiron/sqlx"
)

type WalletUsecase interface {
	GetWallet(ctx context.Context, userID int) (*model.Wallet, error)
	Withdraw(ctx context.Context, userID int, amount int64) (*model.Wallet, error)
}

type walletUsecase struct {
	db    *sqlx.DB
	repos *repository.Repositories
}

func NewWalletService(db *sqlx.DB, repos *repository.Repositories) WalletUsecase {
	return &walletUsecase{db: db, repos: repos}
}

func (u *walletUsecase) GetWallet(ctx context.Context, userID int) (*model.Wallet, error) {
	return u.repos.Wallet.GetWallet(ctx, userID)
}

func (u *walletUsecase) Withdraw(ctx context.Context, userID int, amount int64) (*model.Wallet, error) {
	tx, err := u.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	wallet, err := u.repos.Wallet.GetWalletWithLock(ctx, tx, userID)
	if err != nil {
		return nil, err
	}

	if wallet.Balance < amount {
		return nil, ErrInsufficientBalance
	}

	newBalance := wallet.Balance - amount
	err = u.repos.Wallet.UpdateBalance(ctx, tx, userID, newBalance)
	if err != nil {
		return nil, err
	}

	err = u.repos.Transaction.Insert(ctx, tx, userID, amount)
	if err != nil {
		return nil, err
	}

	wallet.Balance = newBalance
	return wallet, tx.Commit()
}
