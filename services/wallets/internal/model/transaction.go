package model

import "time"

// Transaction represents a wallet transaction for API communication
// This is used for communication with the transaction microservice
type Transaction struct {
	ID              int               `json:"id"`
	SubjectWalletID string            `json:"subject_wallet_id"`
	ObjectWalletID  string            `json:"object_wallet_id,omitempty"`
	TransactionType TransactionType   `json:"transaction_type"`
	OperationType   OperationType     `json:"operation_type"`
	Amount          int64             `json:"amount"` // Amount in cents
	Status          TransactionStatus `json:"status"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
}

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
