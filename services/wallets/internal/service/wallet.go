// Package service provides the business logic for the wallet endpoint.
package service

import (
	"context"
	"errors"
	"github.com/fardinabir/digital-wallet-demo/services/wallets/internal/cache"
	"github.com/fardinabir/digital-wallet-demo/services/wallets/internal/client"
	"github.com/fardinabir/digital-wallet-demo/services/wallets/internal/model"
	"github.com/fardinabir/digital-wallet-demo/services/wallets/internal/repository"
	"github.com/fardinabir/digital-wallet-demo/services/wallets/internal/utils"
)

// Wallet is the service for the wallet endpoint.
type Wallet interface {
	Create(wallet *model.Wallet) error
	Deposit(userID string, amount int, providerID *string) (*model.Transaction, error)
	Withdraw(userID string, amount int, providerID *string) (*model.Transaction, error)
	Transfer(fromUserID string, toUserID string, amount int) (*model.Transaction, error)
	GetWalletWithTransactions(userID string) (*model.Wallet, []model.Transaction, error)
}

type wallet struct {
	walletRepository repository.Wallet
}

// NewWalletService creates a new Wallet service.
func NewWalletService(wr repository.Wallet) Wallet {
	return &wallet{
		walletRepository: wr,
	}
}

func (t *wallet) Create(wallet *model.Wallet) error {
	err := t.walletRepository.Create(wallet)
	if err != nil {
		utils.LogError("Failed to create wallet", err)
		return err
	}
	return nil
}

func (t *wallet) Deposit(userID string, amount int, providerID *string) (*model.Transaction, error) {
	// Validate amount
	if amount <= 0 {
		return nil, errors.New("invalid amount")
	}
	amountCents := int64(amount)

	// FetchTransactions user wallet
	userWallet, err := t.walletRepository.FindByUserID(userID)
	if err != nil {
		utils.LogError("User wallet not found for deposit", err)
		return nil, err
	}

	// Set default provider if not provided
	defaultProviderID := "deposit-provider-master"
	if providerID == nil {
		providerID = &defaultProviderID
	}

	// FetchTransactions or get provider wallet
	providerWallet, err := t.walletRepository.FindProviderWallet(*providerID)
	if err != nil {
		utils.LogError("Provider wallet not found for deposit", err)
		return nil, errors.New("deposit provider wallet not found")
	}

	// Begin database transaction
	tx := t.walletRepository.BeginTransaction()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return nil, err
	}

	// Create debit transaction for provider
	debitTxn := &model.Transaction{
		SubjectWalletID: providerWallet.UserID,
		ObjectWalletID:  userWallet.UserID,
		TransactionType: model.Deposit,
		OperationType:   model.Debit,
		Amount:          amountCents,
		Status:          model.Completed,
	}

	// Create credit transaction for user
	creditTxn := &model.Transaction{
		SubjectWalletID: userWallet.UserID,
		ObjectWalletID:  providerWallet.UserID,
		TransactionType: model.Deposit,
		OperationType:   model.Credit,
		Amount:          amountCents,
		Status:          model.Completed,
	}

	// Create transaction pair via microservice asynchronously
	go func() {
		if err := client.NewTxnClient().CreateTransactionPair(debitTxn, creditTxn); err != nil {
			utils.LogError("Failed to create transaction pair for deposit", err)
		}
	}()

	// Update wallet balances
	if err := t.walletRepository.UpdateWalletBalance(tx, providerWallet.ID, amountCents, false); err != nil {
		utils.LogError("Failed to update provider wallet balance for deposit", err)
		tx.Rollback()
		return nil, err
	}

	if err := t.walletRepository.UpdateWalletBalance(tx, userWallet.ID, amountCents, true); err != nil {
		utils.LogError("Failed to update user wallet balance for deposit", err)
		tx.Rollback()
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		utils.LogError("Failed to commit deposit transaction", err)
		return nil, err
	}

	// Invalidate cache for both user and provider
	ctx := context.Background()
	redisClient := cache.NewRedisClient()
	if err := redisClient.DeleteTransactionHistory(ctx, userWallet.UserID); err != nil {
		utils.LogError("Failed to invalidate user cache after deposit", err)
	}
	if err := redisClient.DeleteTransactionHistory(ctx, providerWallet.UserID); err != nil {
		utils.LogError("Failed to invalidate provider cache after deposit", err)
	}

	// Return the credit transaction for the user
	return creditTxn, nil
}

func (t *wallet) Withdraw(userID string, amount int, providerID *string) (*model.Transaction, error) {
	// Validate amount
	if amount <= 0 {
		return nil, errors.New("invalid amount")
	}
	amountCents := int64(amount)

	// FetchTransactions user wallet
	userWallet, err := t.walletRepository.FindByUserID(userID)
	if err != nil {
		utils.LogError("User wallet not found for withdraw", err)
		return nil, err
	}

	// Check balance
	if userWallet.Balance < amountCents {
		return nil, model.ErrInsufficientFunds
	}

	// Set default provider if not provided
	defaultProviderID := "withdraw-provider-master"
	if providerID == nil {
		providerID = &defaultProviderID
	}

	// FetchTransactions or get provider wallet
	providerWallet, err := t.walletRepository.FindProviderWallet(*providerID)
	if err != nil {
		utils.LogError("Provider wallet not found for withdraw", err)
		return nil, errors.New("withdraw provider wallet not found")
	}

	// Begin database transaction
	tx := t.walletRepository.BeginTransaction()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return nil, err
	}

	// Create debit transaction for user
	debitTxn := &model.Transaction{
		SubjectWalletID: userWallet.UserID,
		ObjectWalletID:  providerWallet.UserID,
		TransactionType: model.Withdraw,
		OperationType:   model.Debit,
		Amount:          amountCents,
		Status:          model.Completed,
	}

	// Create credit transaction for provider
	creditTxn := &model.Transaction{
		SubjectWalletID: providerWallet.UserID,
		ObjectWalletID:  userWallet.UserID,
		TransactionType: model.Withdraw,
		OperationType:   model.Credit,
		Amount:          amountCents,
		Status:          model.Completed,
	}

	// Create transaction pair via microservice asynchronously
	go func() {
		if err := client.NewTxnClient().CreateTransactionPair(debitTxn, creditTxn); err != nil {
			utils.LogError("Failed to create transaction pair for withdraw", err)
		}
	}()

	// Update wallet balances
	if err := t.walletRepository.UpdateWalletBalance(tx, userWallet.ID, amountCents, false); err != nil {
		utils.LogError("Failed to update user wallet balance for withdraw", err)
		tx.Rollback()
		return nil, err
	}

	if err := t.walletRepository.UpdateWalletBalance(tx, providerWallet.ID, amountCents, true); err != nil {
		utils.LogError("Failed to update provider wallet balance for withdraw", err)
		tx.Rollback()
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		utils.LogError("Failed to commit withdraw transaction", err)
		return nil, err
	}

	// Invalidate cache for both user and provider
	ctx := context.Background()
	redisClient := cache.NewRedisClient()
	if err := redisClient.DeleteTransactionHistory(ctx, userWallet.UserID); err != nil {
		utils.LogError("Failed to invalidate user cache after withdraw", err)
	}
	if err := redisClient.DeleteTransactionHistory(ctx, providerWallet.UserID); err != nil {
		utils.LogError("Failed to invalidate provider cache after withdraw", err)
	}

	// Return the debit transaction for the user
	return debitTxn, nil
}

func (t *wallet) Transfer(fromUserID string, toUserID string, amount int) (*model.Transaction, error) {
	// Validate amount
	if amount <= 0 {
		return nil, errors.New("invalid amount")
	}
	amountCents := int64(amount)

	// FetchTransactions sender wallet to check balance
	fromWallet, err := t.walletRepository.FindByUserID(fromUserID)
	if err != nil {
		utils.LogError("Sender wallet not found for transfer", err)
		return nil, err
	}

	// Check balance
	if fromWallet.Balance < amountCents {
		return nil, model.ErrInsufficientFunds
	}

	// FetchTransactions receiver wallet
	toWallet, err := t.walletRepository.FindByUserID(toUserID)
	if err != nil {
		utils.LogError("Receiver wallet not found for transfer", err)
		return nil, err
	}

	// Begin database transaction
	tx := t.walletRepository.BeginTransaction()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return nil, err
	}

	// Create debit transaction for sender
	debitTxn := &model.Transaction{
		SubjectWalletID: fromWallet.UserID,
		ObjectWalletID:  toWallet.UserID,
		TransactionType: model.Transfer,
		OperationType:   model.Debit,
		Amount:          amountCents,
		Status:          model.Completed,
	}

	// Create credit transaction for receiver
	creditTxn := &model.Transaction{
		SubjectWalletID: toWallet.UserID,
		ObjectWalletID:  fromWallet.UserID,
		TransactionType: model.Transfer,
		OperationType:   model.Credit,
		Amount:          amountCents,
		Status:          model.Completed,
	}

	// Create transaction pair via microservice asynchronously
	go func() {
		if err := client.NewTxnClient().CreateTransactionPair(debitTxn, creditTxn); err != nil {
			utils.LogError("Failed to create transaction pair for transfer", err)
		}
	}()

	// Update wallet balances
	if err := t.walletRepository.UpdateWalletBalance(tx, fromWallet.ID, amountCents, false); err != nil {
		utils.LogError("Failed to update sender wallet balance for transfer", err)
		tx.Rollback()
		return nil, err
	}

	if err := t.walletRepository.UpdateWalletBalance(tx, toWallet.ID, amountCents, true); err != nil {
		utils.LogError("Failed to update receiver wallet balance for transfer", err)
		tx.Rollback()
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		utils.LogError("Failed to commit transfer transaction", err)
		return nil, err
	}

	// Invalidate cache for both sender and receiver
	ctx := context.Background()
	redisClient := cache.NewRedisClient()
	if err := redisClient.DeleteTransactionHistory(ctx, fromWallet.UserID); err != nil {
		utils.LogError("Failed to invalidate sender cache after transfer", err)
	}
	if err := redisClient.DeleteTransactionHistory(ctx, toWallet.UserID); err != nil {
		utils.LogError("Failed to invalidate receiver cache after transfer", err)
	}

	// Return the debit transaction for the sender
	return debitTxn, nil
}

func (t *wallet) GetWalletWithTransactions(userID string) (*model.Wallet, []model.Transaction, error) {
	// Get wallet
	wallet, err := t.walletRepository.FindByUserID(userID)
	if err != nil {
		utils.LogError("Wallet not found", err)
		return nil, nil, err
	}

	ctx := context.Background()
	redisClient := cache.NewRedisClient()

	// Try to get transactions from Redis cache first
	transactions, err := redisClient.GetTransactionHistory(ctx, wallet.UserID)
	if err != nil {
		utils.LogError("Failed to get transactions from cache", err)
		// Continue to fetch from transaction service
	}

	// If cache miss or error, fetch from transaction microservice
	if transactions == nil {
		transactions, err = client.NewTxnClient().FetchTransactions(wallet.UserID)
		if err != nil {
			utils.LogError("Failed to retrieve transactions from transaction service", err)
			return nil, nil, err
		}

		// Save to cache for future requests
		if err := redisClient.SaveTransactionHistory(ctx, wallet.UserID, transactions); err != nil {
			utils.LogError("Failed to save transactions to cache", err)
			// Continue without caching - not a critical error
		}
	}

	return wallet, transactions, nil
}
