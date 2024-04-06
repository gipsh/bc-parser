package storage

import (
	"fmt"
	"sync"

	"github.com/gipsh/bc-parser/internal/client"
)

// in-memory stoarge for transactions
// - store transactions
// - thread safe access
// - extendible

type Storer interface {
	Store(address string, tx client.Transaction)
	GetTx(address string) []client.Transaction
	HasAddress(address string) bool
	InitAddress(address string) error // used for subscriptions
}

type Storage struct {
	transactions map[string][]client.Transaction
	mutex        sync.Mutex
}

// NewStorage creates new storage
func NewStorage() *Storage {
	return &Storage{
		transactions: make(map[string][]client.Transaction),
	}
}

// Store transaction
func (s *Storage) Store(address string, tx client.Transaction) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.transactions[address] = append(s.transactions[address], tx)
}

// GetTx returns transactions for address
func (s *Storage) GetTx(address string) []client.Transaction {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.transactions[address]
}

// HasAddress checks if address is in storage
func (s *Storage) HasAddress(address string) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	_, ok := s.transactions[address]
	return ok
}

// InitAddress initializes address in storage
func (s *Storage) InitAddress(address string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if _, ok := s.transactions[address]; !ok {
		s.transactions[address] = make([]client.Transaction, 0)
		return nil
	}
	return fmt.Errorf("address %s already exists", address)
}
