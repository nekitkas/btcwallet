package store

type Store interface {
	Transaction() TransactionRepository
}
