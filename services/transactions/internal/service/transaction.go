package service

import (
	"github.com/fardinabir/digital-wallet-demo/services/transactions/internal/model"
	"github.com/fardinabir/digital-wallet-demo/services/transactions/internal/repository"
)

// TransactionService provides transaction operations
type TransactionService interface {
	CreateTransactionPair(debitTxn, creditTxn *model.Transaction) error
	GetTransactions(subjectWalletID string) ([]model.Transaction, error)
}

type transactionService struct {
	repo repository.TransactionRepository
}

// NewTransactionService creates a new transaction service
func NewTransactionService(repo repository.TransactionRepository) TransactionService {
	return &transactionService{repo: repo}
}

// CreateTransactionPair creates both debit and credit transactions atomically
func (s *transactionService) CreateTransactionPair(debitTxn, creditTxn *model.Transaction) error {
	return s.repo.CreateTransactionPair(debitTxn, creditTxn)
}

// GetTransactions retrieves all transactions for a specific wallet
func (s *transactionService) GetTransactions(subjectWalletID string) ([]model.Transaction, error) {
	filters := map[string]interface{}{
		"subject_wallet_id": subjectWalletID,
	}
	return s.repo.FindAllTransactions(filters)
}