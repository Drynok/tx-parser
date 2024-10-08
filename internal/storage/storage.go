package storage

import (
	"sync"

	"github.com/Drynok/tx-parser/internal/model"
)

type Storage interface {
	AddTransaction(address string, tx model.Transaction) error
	Transactions(address string) []model.Transaction
	Subscribe(address string) bool
	IsSubscribed(address string) bool
}

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

func (m *MemoryStorage) IsSubscribed(address string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, ok := m.subscribers[address]
	return ok
}

func (m *MemoryStorage) Subscribe(address string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.subscribers[address]; ok {
		return false
	}
	m.subscribers[address] = true
	return true
}
