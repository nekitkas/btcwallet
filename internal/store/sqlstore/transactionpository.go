package sqlstore

import (
	"btcwallet/internal/model"
	"errors"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type TransactionRepository struct {
	store *Store
}

func (t TransactionRepository) Get() (*[]model.Transaction, error) {
	query := `SELECT * FROM transactions`
	rows, err := t.store.db.Query(query)
	if err != nil {
		return nil, err
	}

	var transactions []model.Transaction

	for rows.Next() {
		var transaction model.Transaction
		if err = rows.Scan(&transaction.ID, &transaction.Amount, &transaction.Spent, &transaction.CreatedAt); err != nil {
			return nil, err
		}

		transactions = append(transactions, transaction)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &transactions, nil
}

func (t TransactionRepository) Transfer(amount float64) error {
	rows, err := t.store.db.Query("SELECT id, amount FROM transactions WHERE spent = FALSE ORDER BY created_at FOR UPDATE")
	if err != nil {
		return err
	}
	defer rows.Close()

	var totalUnspent float64
	var transactionIDs []string

	for rows.Next() {
		var id string
		var value float64
		if err = rows.Scan(&id, &value); err != nil {
			return err
		}
		totalUnspent += value
		transactionIDs = append(transactionIDs, id)
		if totalUnspent >= amount {
			break
		}
	}

	if err = rows.Err(); err != nil {
		return err
	}

	if totalUnspent < amount || amount < 0.00001 {
		return errors.New("not enough balance or transfer amount is too small")
	}

	_, err = t.store.db.Exec("UPDATE transactions SET spent = TRUE WHERE id = ANY($1)", pq.Array(transactionIDs))
	if err != nil {
		return err
	}

	leftover := totalUnspent - amount
	if leftover >= 0.00001 {
		newTransactionID := uuid.New().String()
		_, err = t.store.db.Exec("INSERT INTO transactions (id, amount, spent, created_at) VALUES ($1, $2, FALSE, NOW())", newTransactionID, leftover)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t TransactionRepository) Balance() (float64, error) {
	var totalAmount float64
	if err := t.store.db.QueryRow(
		`SELECT SUM(amount) AS total_amount FROM transactions WHERE spent = false;`).Scan(&totalAmount); err != nil {
		return 0, err
	}

	return totalAmount, nil
}

func (t TransactionRepository) Add(transaction *model.Transaction, amount float64) error {
	_, err := t.store.db.Exec(`INSERT INTO transactions (id, amount, spent, created_at) VALUES ($1,$2,$3,$4)`,
		transaction.ID,
		amount,
		transaction.Spent,
		transaction.CreatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}
