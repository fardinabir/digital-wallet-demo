// Package repository provides the database operations for the wallet endpoint.
package repository

import (
	"github.com/fardinabir/digital-wallet-demo/internal/model"
	"gorm.io/gorm"
)

// Wallet is the repository for the wallet and transaction operations.
type Wallet interface {
	// Wallet operations
	Create(t *model.Wallet) error
	Find(id int) (*model.Wallet, error)
	FindByUserID(userID int) (*model.Wallet, error)
	FindProviderWallet(providerID int) (*model.Wallet, error)

	// Transaction operations
	CreateTransaction(t *model.Transaction) error
	CreateTransactionWithTx(tx *gorm.DB, t *model.Transaction) error
	FindAllTransactions(qry map[string]interface{}) ([]model.Transaction, error)

	// Atomic operations
	BeginTransaction() *gorm.DB
	UpdateWalletBalance(tx *gorm.DB, walletID int, amount int64, isCredit bool) error
}

type wallet struct {
	db *gorm.DB
}

// NewWallet returns a new instance of the wallet repository.
func NewWallet(db *gorm.DB) Wallet {
	return &wallet{
		db: db,
	}
}

func (td *wallet) Create(t *model.Wallet) error {
	if err := td.db.Create(t).Error; err != nil {
		return err
	}
	return nil
}

func (td *wallet) Find(id int) (*model.Wallet, error) {
	var wallet *model.Wallet
	err := td.db.Where("id = ?", id).Take(&wallet).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, model.ErrNotFound
		}
		return nil, err
	}
	return wallet, nil
}

func (td *wallet) FindByUserID(userID int) (*model.Wallet, error) {
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

// Transaction operations
func (td *wallet) CreateTransaction(t *model.Transaction) error {
	if err := td.db.Create(t).Error; err != nil {
		return err
	}
	return nil
}

func (td *wallet) CreateTransactionWithTx(tx *gorm.DB, t *model.Transaction) error {
	if err := tx.Create(t).Error; err != nil {
		return err
	}
	return nil
}

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

func (td *wallet) FindProviderWallet(providerID int) (*model.Wallet, error) {
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

func (td *wallet) BeginTransaction() *gorm.DB {
	return td.db.Begin()
}

func (td *wallet) UpdateWalletBalance(tx *gorm.DB, walletID int, amount int64, isCredit bool) error {
	var wallet model.Wallet
	if err := tx.Where("id = ?", walletID).First(&wallet).Error; err != nil {
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
