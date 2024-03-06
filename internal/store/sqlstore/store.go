package sqlstore

import (
	"btcwallet/internal/store"
	"database/sql"
)

type Store struct {
	db              *sql.DB
	transactionRepo *TransactionRepository
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
