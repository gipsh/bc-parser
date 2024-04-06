package storage

import (
	"testing"

	"github.com/gipsh/bc-parser/internal/client"
	"github.com/stretchr/testify/assert"
)

func TestNewStorageStore(t *testing.T) {

	s := NewStorage()
	s.Store("0x123", client.Transaction{From: "0x123", To: "0x456"})
	s.Store("0x123", client.Transaction{From: "0x123", To: "0x456"})

	txs := s.GetTx("0x123")

	assert.Equal(t, 2, len(txs))
}

func TestNewStorageHasKey(t *testing.T) {

	s := NewStorage()
	s.InitAddress("0x123")
	b := s.HasAddress("0x123")
	assert.Equal(t, true, b)

	err := s.InitAddress("0x123")
	assert.NotNil(t, err)
	assert.Error(t, err, "address already exists")

}
