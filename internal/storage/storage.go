package storage

import (
	"github.com/Drynok/tx-parser/internal/model"
)

type Storage interface {
	AddTransaction(address string, tx model.Transaction) error
	Transactions(address string) []model.Transaction
	// Subscribes() map[string]bool
}
