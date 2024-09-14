package parser

import (
	"github.com/Drynok/tx-parser/internal/model"
)

type Parser interface {
	GetCurrentBlock() int
	Subscribe(address string) bool
	GetTransactions(address string) []model.Transaction
}
