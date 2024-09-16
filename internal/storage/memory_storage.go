package storage

import (
	"sync"

	"github.com/Drynok/tx-parser/internal/model"
)

type MemoryStorage struct {
	subscribers  map[string]bool
	transactions map[string][]model.Transaction
	mu           sync.RWMutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		transactions: make(map[string][]model.Transaction),
		subscribers:  make(map[string]bool),
	}
}

func (m *MemoryStorage) AddTransaction(address string, tx model.Transaction) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.transactions[address] = append(m.transactions[address], tx)
	return nil
}

func (m *MemoryStorage) Transactions(address string) []model.Transaction {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.transactions[address]
}
