// Package repository provides the database operations for the wallet endpoint.
package repository

import (
	"github.com/fardinabir/digital-wallet-demo/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Wallet provides database operations for wallet and transaction management.
type Wallet interface {
	// Wallet operations
	Create(t *model.Wallet) error
	FindByUserID(userID string) (*model.Wallet, error)
	FindProviderWallet(providerID string) (*model.Wallet, error)

	// Transaction operations
	InsertTransaction(tx *gorm.DB, t *model.Transaction) error
	FindAllTransactions(qry map[string]interface{}) ([]model.Transaction, error)

	// Atomic operations
	BeginTransaction() *gorm.DB
	UpdateWalletBalance(tx *gorm.DB, walletID int, amount int64, isCredit bool) error
}

type wallet struct {
	db *gorm.DB
}

// NewWallet creates a new wallet repository instance.
func NewWallet(db *gorm.DB) Wallet {
	return &wallet{
		db: db,
	}
}

// Create inserts a new wallet record into the database.
func (td *wallet) Create(t *model.Wallet) error {
	if err := td.db.Create(t).Error; err != nil {
		return err
	}
	return nil
}

// FindByUserID retrieves a wallet by user ID, returns ErrNotFound if not exists.
func (td *wallet) FindByUserID(userID string) (*model.Wallet, error) {
	var wallet *model.Wallet
	err := td.db.Where("user_id = ?", userID).Take(&wallet).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, model.ErrNotFound
		}
		return nil, err
	}
	return wallet, nil
}

// InsertTransaction creates a transaction record within the given database transaction.
func (td *wallet) InsertTransaction(tx *gorm.DB, t *model.Transaction) error {
	if err := tx.Create(t).Error; err != nil {
		return err
	}
	return nil
}

// FindAllTransactions retrieves transactions matching the query filters, ordered by creation date desc.
func (td *wallet) FindAllTransactions(qry map[string]interface{}) ([]model.Transaction, error) {
	var transactions []model.Transaction
	tx := td.db

	if len(qry) > 0 {
		tx = tx.Where(qry)
	}
	err := tx.Order("created_at desc").Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

// FindProviderWallet retrieves a provider wallet by provider ID for system operations.
func (td *wallet) FindProviderWallet(providerID string) (*model.Wallet, error) {
	var wallet *model.Wallet
	err := td.db.Where("user_id = ? AND acnt_type = ?", providerID, model.Provider).Take(&wallet).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, model.ErrNotFound
		}
		return nil, err
	}
	return wallet, nil
}

// BeginTransaction starts a new database transaction for atomic operations.
func (td *wallet) BeginTransaction() *gorm.DB {
	return td.db.Begin()
}

// UpdateWalletBalance atomically updates wallet balance
// Used row-level Exclusive Locking to ensure single transaction can update the wallet balance at a time
func (td *wallet) UpdateWalletBalance(tx *gorm.DB, walletID int, amount int64, isCredit bool) error {
	var wallet model.Wallet

	// Acquire row-level lock
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", walletID).First(&wallet).Error; err != nil {
		return err
	}

	if isCredit {
		wallet.Balance += amount
	} else {
		wallet.Balance -= amount
		if wallet.Balance < 0 {
			return model.ErrInsufficientFunds
		}
	}

	return tx.Save(&wallet).Error
}
