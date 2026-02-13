package repository

import "github.com/jmoiron/sqlx"

type Repositories struct {
	Wallet      WalletRepository
	Transaction TransactionRepository
}

func NewRepositories(db *sqlx.DB) *Repositories {
	return &Repositories{
		Wallet:      NewWalletRepository(db),
		Transaction: NewTransactionRepository(db),
	}
}
