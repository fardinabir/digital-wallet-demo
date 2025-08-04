package client

import (
	"github.com/fardinabir/digital-wallet-demo/services/wallets/internal/model"
)

// MockTransactionClient implements the NewTransaction interface for testing
type MockTransactionClient struct{}

func (m *MockTransactionClient) CreateTransactionPair(debitTxn, creditTxn *model.Transaction) error {
	// Mock successful transaction creation
	// In a real scenario, this would make HTTP calls to the transaction service
	// But for testing, we just return success
	return nil
}

func (m *MockTransactionClient) FetchTransactions(subjectWalletID string) ([]model.Transaction, error) {
	// For test-user-001, return some sample transactions
	if subjectWalletID == "test-user-001" {
		return []model.Transaction{
			{
				SubjectWalletID: "test-user-001",
				ObjectWalletID:  "deposit-provider-master",
				TransactionType: model.Deposit,
				OperationType:   model.Credit,
				Amount:          5000,
				Status:          model.Completed,
			},
			{
				SubjectWalletID: "test-user-001",
				ObjectWalletID:  "withdraw-provider-master",
				TransactionType: model.Withdraw,
				OperationType:   model.Debit,
				Amount:          2000,
				Status:          model.Completed,
			},
		}, nil
	}
	// For other wallet IDs, return empty list
	return []model.Transaction{}, nil
}
