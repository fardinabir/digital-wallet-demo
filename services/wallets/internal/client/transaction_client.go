package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/fardinabir/digital-wallet-demo/services/wallets/internal/config"
	"github.com/fardinabir/digital-wallet-demo/services/wallets/internal/model"
	"github.com/fardinabir/digital-wallet-demo/services/wallets/internal/utils"
)

// NewTransaction interface for communicating with transactions microservice
type NewTransaction interface {
	CreateTransactionPair(debitTxn, creditTxn *model.Transaction) error
	FetchTransactions(subjectWalletID string) ([]model.Transaction, error)
}

type transactionClient struct {
	client  *http.Client
	baseURL string
}

var (
	instance NewTransaction
	once     sync.Once
)

// ResetClient resets the singleton instance for testing purposes
func ResetClient() {
	once = sync.Once{}
	instance = nil
}

// NewTxnClient is a factory method that returns NewTransaction interface with singleton pattern
func NewTxnClient() NewTransaction {
	once.Do(func() {
		baseURL := "http://localhost:8082" // default fallback
		globalConfig := config.GetGlobalConfig()
		baseURL = globalConfig.Services.Transaction.BaseURL

		instance = &transactionClient{
			client: &http.Client{
				Timeout: 30 * time.Second,
			},
			baseURL: baseURL,
		}
	})
	return instance
}

// TransactionPairRequest represents the request payload for creating transaction pairs
type TransactionPairRequest struct {
	DebitTransaction  TransactionRequest `json:"debit_transaction"`
	CreditTransaction TransactionRequest `json:"credit_transaction"`
}

// TransactionRequest represents a single transaction in the request
type TransactionRequest struct {
	SubjectWalletID string                  `json:"subject_wallet_id"`
	ObjectWalletID  string                  `json:"object_wallet_id"`
	TransactionType model.TransactionType   `json:"transaction_type"`
	OperationType   model.OperationType     `json:"operation_type"`
	Amount          int64                   `json:"amount"`
	Status          model.TransactionStatus `json:"status"`
}

// TransactionResponse represents the API response wrapper for transactions
type TransactionResponse struct {
	Data []model.Transaction `json:"data"`
}

// FetchTransactions retrieves transactions for a specific wallet from the transaction service
func (tc *transactionClient) FetchTransactions(subjectWalletID string) ([]model.Transaction, error) {
	// Create HTTP request
	url := fmt.Sprintf("%s/api/v1/transactions/%s", tc.baseURL, subjectWalletID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		utils.LogError("Failed to create HTTP request for fetching transactions", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := tc.client.Do(req)
	if err != nil {
		utils.LogError("Failed to send fetch transactions request", err)
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		utils.LogError(fmt.Sprintf("Transaction microservice returned status %d", resp.StatusCode), nil)
		return nil, fmt.Errorf("transaction service returned status %d", resp.StatusCode)
	}

	// Parse response
	var response TransactionResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		utils.LogError("Failed to decode transactions response", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return response.Data, nil
}

// CreateTransactionPair sends both debit and credit transactions to the transactions microservice
func (tc *transactionClient) CreateTransactionPair(debitTxn, creditTxn *model.Transaction) error {
	// Prepare the request payload
	request := TransactionPairRequest{
		DebitTransaction: TransactionRequest{
			SubjectWalletID: debitTxn.SubjectWalletID,
			ObjectWalletID:  debitTxn.ObjectWalletID,
			TransactionType: debitTxn.TransactionType,
			OperationType:   debitTxn.OperationType,
			Amount:          debitTxn.Amount,
			Status:          debitTxn.Status,
		},
		CreditTransaction: TransactionRequest{
			SubjectWalletID: creditTxn.SubjectWalletID,
			ObjectWalletID:  creditTxn.ObjectWalletID,
			TransactionType: creditTxn.TransactionType,
			OperationType:   creditTxn.OperationType,
			Amount:          creditTxn.Amount,
			Status:          creditTxn.Status,
		},
	}

	// Marshal the request to JSON
	jsonData, err := json.Marshal(request)
	if err != nil {
		utils.LogError("Failed to marshal transaction pair request", err)
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	url := fmt.Sprintf("%s/api/v1/transactions", tc.baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		utils.LogError("Failed to create HTTP request for transaction pair", err)
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := tc.client.Do(req)
	if err != nil {
		utils.LogError("Failed to send transaction pair request", err)
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusCreated {
		utils.LogError(fmt.Sprintf("Transaction microservice returned status %d", resp.StatusCode), nil)
		return fmt.Errorf("transaction service returned status %d", resp.StatusCode)
	}

	// Successfully created transaction pair
	return nil
}
