package store

import "btcwallet/internal/model"

type TransactionRepository interface {
	// Get list all transactions
	Get() (*[]model.Transaction, error)
	// Transfer creates a new spending by providing amount in EUR
	Transfer(amount float64) error
	// Balance Show current balance
	Balance() (float64, error)
	// Add Adds new unspent transaction by providing amount in EUR
	Add(transaction *model.Transaction, amount float64) error
}
