package storage

import (
	"github.com/Drynok/tx-parser/internal/model"
)

type Storage interface {
	AddTransaction(address string, tx model.Transaction) error
	GetTransactions(address string) []model.Transaction
}
