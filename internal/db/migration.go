package db

import (
	"github.com/fardinabir/digital-wallet-demo/internal/model"
	"gorm.io/gorm"
)

// Migrate runs the auto-migration for the database
func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&model.Wallet{}, &model.Transaction{}); err != nil {
		return err
	}
	return nil
}
