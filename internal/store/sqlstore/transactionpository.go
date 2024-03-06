package sqlstore

import (
	"btcwallet/internal/model"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"strings"
	"time"
)

type TransactionRepository struct {
	store *Store
}

func (t TransactionRepository) Get() (*[]model.Transaction, error) {
	query := `SELECT id, amount, spent, created_at FROM transactions`
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
	tx, err := t.store.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	rows, err := tx.Query("SELECT id, amount FROM transactions WHERE spent = 0 ORDER BY created_at")
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

	placeholders := make([]string, len(transactionIDs))
	values := make([]interface{}, len(transactionIDs))
	for i, id := range transactionIDs {
		placeholders[i] = "?"
		values[i] = id
	}

	placeholdersStr := strings.Join(placeholders, ", ")
	updateQuery := fmt.Sprintf("UPDATE transactions SET spent = 1 WHERE id IN (%s)", placeholdersStr)
	if _, err = tx.Exec(updateQuery, values...); err != nil {
		return err
	}

	leftover := totalUnspent - amount
	if leftover >= 0.00001 {
		id := uuid.New().String()
		if _, err = tx.Exec("INSERT INTO transactions (id, amount, spent, created_at) VALUES (?, ?, 0, ?)",
			id, leftover, time.Now().Format("2006-01-02 15:04:05")); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (t TransactionRepository) Balance() (float64, error) {
	var totalAmount float64
	if err := t.store.db.QueryRow(`SELECT IFNULL(SUM(amount), 0) AS total_amount FROM transactions WHERE spent = 0`).Scan(&totalAmount); err != nil {
		return 0, err
	}

	return totalAmount, nil
}

func (t TransactionRepository) Add(transaction *model.Transaction, amount float64) error {
	_, err := t.store.db.Exec(`INSERT INTO transactions (id, amount, spent, created_at) VALUES (?, ?, ?, datetime('now'))`,
		transaction.ID,
		amount,
		transaction.Spent,
	)
	if err != nil {
		return err
	}

	return nil
}
