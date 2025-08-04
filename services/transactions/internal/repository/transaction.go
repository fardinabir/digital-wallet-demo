package repository

import (
	"github.com/fardinabir/digital-wallet-demo/services/transactions/internal/model"
	"gorm.io/gorm"
)

// TransactionRepository provides database operations for transactions
type TransactionRepository interface {
	CreateTransactionPair(debitTxn, creditTxn *model.Transaction) error
	FindAllTransactions(filters map[string]interface{}) ([]model.Transaction, error)
}

type transactionRepository struct {
	db *gorm.DB
}

// NewTransactionRepository creates a new transaction repository
func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

// CreateTransactionPair creates both debit and credit transactions atomically
func (r *transactionRepository) CreateTransactionPair(debitTxn, creditTxn *model.Transaction) error {
	// Begin database transaction
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	// Insert debit transaction
	if err := tx.Create(debitTxn).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Insert credit transaction
	if err := tx.Create(creditTxn).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit transaction
	return tx.Commit().Error
}

// FindAllTransactions retrieves transactions matching the query filters
func (r *transactionRepository) FindAllTransactions(filters map[string]interface{}) ([]model.Transaction, error) {
	var transactions []model.Transaction
	tx := r.db

	if len(filters) > 0 {
		tx = tx.Where(filters)
	}

	err := tx.Order("created_at desc").Find(&transactions).Error
	if err != nil {
		return nil, err
	}

	return transactions, nil
}
