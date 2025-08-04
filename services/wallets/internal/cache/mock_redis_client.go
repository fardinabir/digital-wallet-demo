// Package cache provides mock Redis client for testing
package cache

import (
	"context"
	"github.com/fardinabir/digital-wallet-demo/services/wallets/internal/model"
)

// MockRedisClient implements RedisClient interface for testing
type MockRedisClient struct {
	Transactions map[string][]model.Transaction
}

// NewMockRedisClient creates a new mock Redis client
func NewMockRedisClient() *MockRedisClient {
	return &MockRedisClient{
		Transactions: make(map[string][]model.Transaction),
	}
}

// GetTransactionHistory returns mock transaction history
func (m *MockRedisClient) GetTransactionHistory(ctx context.Context, userID string) ([]model.Transaction, error) {
	if transactions, exists := m.Transactions[userID]; exists {
		return transactions, nil
	}
	return nil, nil // Cache miss
}

// SaveTransactionHistory saves mock transaction history
func (m *MockRedisClient) SaveTransactionHistory(ctx context.Context, userID string, transactions []model.Transaction) error {
	m.Transactions[userID] = transactions
	return nil
}

// DeleteTransactionHistory deletes mock transaction history
func (m *MockRedisClient) DeleteTransactionHistory(ctx context.Context, userID string) error {
	delete(m.Transactions, userID)
	return nil
}

// Close does nothing for mock client
func (m *MockRedisClient) Close() error {
	return nil
}
