package model

import (
	"time"

	"github.com/go-playground/validator/v10"
)

// Transaction represents a wallet transaction
type Transaction struct {
	ID              int               `gorm:"primaryKey" json:"id"`
	SubjectWalletID string            `gorm:"not null;" json:"subject_wallet_id"`
	ObjectWalletID  string            `gorm:"not null;" json:"object_wallet_id,omitempty"`
	TransactionType TransactionType   `gorm:"not null" json:"transaction_type"`
	OperationType   OperationType     `gorm:"not null" json:"operation_type"`
	Amount          int64             `gorm:"not null" json:"amount"` // Amount in cents
	Status          TransactionStatus `gorm:"default:'pending'" json:"status"`
	CreatedAt       time.Time         `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time         `gorm:"autoUpdateTime" json:"updated_at"`
}

// NewTransaction returns a new instance of the Transaction model.
func NewTransaction(subjectWalletID, objectWalletID string, transactionType TransactionType, operationType OperationType, amount int64) *Transaction {
	return &Transaction{
		SubjectWalletID: subjectWalletID,
		ObjectWalletID:  objectWalletID,
		TransactionType: transactionType,
		OperationType:   operationType,
		Amount:          amount,
		Status:          Pending,
	}
}

// Provider wallet constants for external wallet service communication
const (
	// DepositProviderID is the UserID for the deposit provider wallet
	// This is a master account that acts as the source for all deposit transactions
	DepositProviderID = "deposit-provider-master"
	// WithdrawProviderID is the UserID for the withdraw provider wallet
	// This is a master account that acts as the destination for all withdraw transactions
	WithdrawProviderID = "withdraw-provider-master"
)

// OperationType represents the operation type for transactions
type OperationType string

const (
	// Debit operation type
	Debit = OperationType("debit")
	// Credit operation type
	Credit = OperationType("credit")
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

// IsValidOperationType checks if the operation type is valid
func IsValidOperationType(fl validator.FieldLevel) bool {
	if fl.Field().IsZero() {
		return true
	}
	operationType := fl.Field().Interface().(OperationType)
	return operationType == Debit || operationType == Credit
}
