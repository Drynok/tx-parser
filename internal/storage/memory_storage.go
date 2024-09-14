package storage

import (
	"sync"

	"github.com/Drynok/tx-parser/internal/model"
)

type MemoryStorage struct {
	Subscribers  map[string]bool
	Transactions map[string][]model.Transaction
	mu           sync.RWMutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		Transactions: make(map[string][]model.Transaction),
	}
}

func (m *MemoryStorage) AddTransaction(address string, tx model.Transaction) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Transactions[address] = append(m.Transactions[address], tx)
	return nil
}

func (m *MemoryStorage) GetTransactions(address string) []model.Transaction {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.Transactions[address]
}
