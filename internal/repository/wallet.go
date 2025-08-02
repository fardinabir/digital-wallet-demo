// Package repository provides the database operations for the wallet endpoint.
package repository

import (
	"github.com/fardinabir/digital-wallet-demo/internal/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Wallet is the repository for the wallet endpoint.
type Wallet interface {
	Create(t *model.Wallet) error
	Delete(id int) error
	Update(t *model.Wallet) error
	Find(id int) (*model.Wallet, error)
	FindAll(qry map[string]interface{}) ([]*model.Wallet, error)
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

func (td *wallet) Update(t *model.Wallet) error {
	if err := td.db.Save(t).Error; err != nil {
		return err
	}
	return nil
}

func (td *wallet) Delete(id int) error {
	result := td.db.Where("id = ?", id).Delete(&model.Wallet{})
	if result.RowsAffected == 0 {
		return model.ErrNotFound
	}
	if result.Error != nil {
		return result.Error
	}
	log.Info("Deleted wallet with id: ", id)
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

func (td *wallet) FindAll(qry map[string]interface{}) ([]*model.Wallet, error) {
	var wallets []*model.Wallet
	tx := td.db

	if val, ok := qry["task"].(string); ok {
		tx = tx.Where("task LIKE ?", "%"+val+"%")
		delete(qry, "task")
	}
	if len(qry) > 0 {
		tx = tx.Where(qry)
	}
	err := tx.Order("priority desc").Order("created_at desc").Find(&wallets).Error
	if err != nil {
		return nil, err
	}
	return wallets, nil
}
