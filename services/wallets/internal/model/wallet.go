package model

import (
	"time"

	"github.com/go-playground/validator/v10"
)

// Wallet is the model for the wallet endpoint.
type Wallet struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	UserID    string    `gorm:"not null;uniqueIndex" json:"user_id"`
	AcntType  AcntType  `gorm:"not null" json:"acnt_type"`
	Balance   int64     `gorm:"default:0" json:"balance"` // Balance in cents
	Status    Status    `json:"status"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// NewWallet returns a new instance of the wallet model.
func NewWallet(userID string, acntType AcntType) *Wallet {
	return &Wallet{
		UserID:   userID,
		AcntType: acntType,
		Balance:  0,
		Status:   Active,
	}
}

// AcntType represents the account type
type AcntType string

const (
	// User account type
	User = AcntType("user")
	// Provider account type
	Provider = AcntType("provider")
)

// Provider wallet constants for master accounts
const (
	// DepositProviderID is the UserID for the deposit provider wallet
	// This is a master account that acts as the source for all deposit transactions
	DepositProviderID = "deposit-provider-master"
	// WithdrawProviderID is the UserID for the withdraw provider wallet
	// This is a master account that acts as the destination for all withdraw transactions
	WithdrawProviderID = "withdraw-provider-master"
)

// Status is the status of the wallet.
type Status string

const (
	// Active is the status for an active wallet.
	Active = Status("active")
	// Inactive is the status for an inactive wallet.
	Inactive = Status("inactive")
	// Suspended is the status for a suspended wallet.
	Suspended = Status("suspended")
)

// StatusMap is a map of wallet status.
var StatusMap = map[Status]bool{
	Active:    true,
	Inactive:  true,
	Suspended: true,
}

// IsValidStatus checks if the status is valid (Active, Inactive, Suspended)
func IsValidStatus(fl validator.FieldLevel) bool {
	if fl.Field().IsZero() {
		return true // Skip validation for empty or nil fields
	}
	status := fl.Field().Interface().(Status)
	return status == Active || status == Inactive || status == Suspended
}

// IsValidAcntType checks if the account type is valid
func IsValidAcntType(fl validator.FieldLevel) bool {
	if fl.Field().IsZero() {
		return true
	}
	acntType := fl.Field().Interface().(AcntType)
	return acntType == User || acntType == Provider
}
