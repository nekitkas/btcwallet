package model

import (
	"github.com/google/uuid"
	"time"
)

type Transaction struct {
	ID        string    `db:"id" json:"id"`
	Amount    float64   `db:"amount" json:"amount"`
	Spent     bool      `db:"spent" json:"spent"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

func NewTransaction() *Transaction {
	return &Transaction{
		ID:        uuid.New().String(),
		CreatedAt: time.Now(),
	}
}
