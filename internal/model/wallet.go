package model

import (
	"time"

	"github.com/go-playground/validator/v10"
)

// Wallet is the model for the wallet endpoint.
type Wallet struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	UserID    int       `gorm:"not null;index" json:"user_id"`
	AcntType  AcntType  `gorm:"not null" json:"acnt_type"`
	Balance   int64     `gorm:"default:0" json:"balance"` // Balance in cents
	Status    Status    `json:"status"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// Transaction represents a wallet transaction
type Transaction struct {
	ID              int               `gorm:"primaryKey" json:"id"`
	WalletID        int               `gorm:"not null;index" json:"wallet_id"`
	SubjectWalletID int               `gorm:"not null;index" json:"subject_wallet_id"`
	ObjectWalletID  *int              `gorm:"index" json:"object_wallet_id,omitempty"`
	TransactionType TransactionType   `gorm:"not null" json:"transaction_type"`
	OperationType   OperationType     `gorm:"not null" json:"operation_type"`
	Amount          int64             `gorm:"not null" json:"amount"` // Amount in cents
	Status          TransactionStatus `gorm:"default:'pending'" json:"status"`
	CreatedAt       time.Time         `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time         `gorm:"autoUpdateTime" json:"updated_at"`
}

// NewWallet returns a new instance of the wallet model.
func NewWallet(userID int, acntType AcntType) *Wallet {
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

// OperationType represents the operation type for transactions
type OperationType string

const (
	// Debit operation type
	Debit = OperationType("debit")
	// Credit operation type
	Credit = OperationType("credit")
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

// TransactionType represents the type of transaction
type TransactionType string

const (
	// Deposit transaction type
	Deposit = TransactionType("deposit")
	// Withdraw transaction type
	Withdraw = TransactionType("withdraw")
	// Transfer transaction type
	Transfer = TransactionType("transfer")
)

// TransactionStatus represents the status of a transaction
type TransactionStatus string

const (
	// Pending transaction status
	Pending = TransactionStatus("pending")
	// Completed transaction status
	Completed = TransactionStatus("completed")
	// Failed transaction status
	Failed = TransactionStatus("failed")
	// Cancelled transaction status
	Cancelled = TransactionStatus("cancelled")
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

// IsValidTransactionType checks if the transaction type is valid
func IsValidTransactionType(fl validator.FieldLevel) bool {
	if fl.Field().IsZero() {
		return true
	}
	txnType := fl.Field().Interface().(TransactionType)
	return txnType == Deposit || txnType == Withdraw || txnType == Transfer
}

// IsValidTransactionStatus checks if the transaction status is valid
func IsValidTransactionStatus(fl validator.FieldLevel) bool {
	if fl.Field().IsZero() {
		return true
	}
	status := fl.Field().Interface().(TransactionStatus)
	return status == Pending || status == Completed || status == Failed || status == Cancelled
}

// IsValidAcntType checks if the account type is valid
func IsValidAcntType(fl validator.FieldLevel) bool {
	if fl.Field().IsZero() {
		return true
	}
	acntType := fl.Field().Interface().(AcntType)
	return acntType == User || acntType == Provider
}

// IsValidOperationType checks if the operation type is valid
func IsValidOperationType(fl validator.FieldLevel) bool {
	if fl.Field().IsZero() {
		return true
	}
	operationType := fl.Field().Interface().(OperationType)
	return operationType == Debit || operationType == Credit
}
