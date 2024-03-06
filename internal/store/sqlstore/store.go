package sqlstore

import (
	"btcwallet/internal/store"
	"database/sql"
	_ "github.com/lib/pq"
)

type Store struct {
	db              *sql.DB
	transactionRepo *TransactionRepository
	walletRepo      *WalletRepository
}

func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) Transaction() store.TransactionRepository {
	if s.transactionRepo != nil {
		return s.transactionRepo
	}

	return &TransactionRepository{store: s}
}

func (s *Store) Wallet() store.WalletRepository {
	if s.walletRepo != nil {
		return s.walletRepo
	}

	return &WalletRepository{store: s}
}