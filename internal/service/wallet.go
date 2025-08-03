// Package service provides the business logic for the wallet endpoint.
package service

import (
	"errors"

	"github.com/fardinabir/digital-wallet-demo/internal/model"
	"github.com/fardinabir/digital-wallet-demo/internal/repository"
	"github.com/fardinabir/digital-wallet-demo/internal/utils"
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

// NewWallet creates a new Wallet service.
func NewWallet(wr repository.Wallet) Wallet {
	return &wallet{wr}
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

	// Find user wallet
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

	// Find or get provider wallet
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
	if err := t.walletRepository.InsertTransaction(tx, debitTxn); err != nil {
		utils.LogError("Failed to insert debit transaction for deposit", err)
		tx.Rollback()
		return nil, err
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
	if err := t.walletRepository.InsertTransaction(tx, creditTxn); err != nil {
		utils.LogError("Failed to insert credit transaction for deposit", err)
		tx.Rollback()
		return nil, err
	}

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

	// Return the credit transaction for the user
	return creditTxn, nil
}

func (t *wallet) Withdraw(userID string, amount int, providerID *string) (*model.Transaction, error) {
	// Validate amount
	if amount <= 0 {
		return nil, errors.New("invalid amount")
	}
	amountCents := int64(amount)

	// Find user wallet
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

	// Find or get provider wallet
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
	if err := t.walletRepository.InsertTransaction(tx, debitTxn); err != nil {
		utils.LogError("Failed to insert debit transaction for withdraw", err)
		tx.Rollback()
		return nil, err
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
	if err := t.walletRepository.InsertTransaction(tx, creditTxn); err != nil {
		utils.LogError("Failed to insert credit transaction for withdraw", err)
		tx.Rollback()
		return nil, err
	}

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

	// Return the debit transaction for the user
	return debitTxn, nil
}

func (t *wallet) Transfer(fromUserID string, toUserID string, amount int) (*model.Transaction, error) {
	// Validate amount
	if amount <= 0 {
		return nil, errors.New("invalid amount")
	}
	amountCents := int64(amount)

	// Find sender wallet to check balance
	fromWallet, err := t.walletRepository.FindByUserID(fromUserID)
	if err != nil {
		utils.LogError("Sender wallet not found for transfer", err)
		return nil, err
	}

	// Check balance
	if fromWallet.Balance < amountCents {
		return nil, model.ErrInsufficientFunds
	}

	// Find receiver wallet
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
	if err := t.walletRepository.InsertTransaction(tx, debitTxn); err != nil {
		utils.LogError("Failed to insert debit transaction for transfer", err)
		tx.Rollback()
		return nil, err
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
	if err := t.walletRepository.InsertTransaction(tx, creditTxn); err != nil {
		utils.LogError("Failed to insert credit transaction for transfer", err)
		tx.Rollback()
		return nil, err
	}

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

	// Get transactions
	filters := map[string]interface{}{
		"subject_wallet_id": wallet.UserID,
	}
	transactions, err := t.walletRepository.FindAllTransactions(filters)
	if err != nil {
		utils.LogError("Failed to retrieve transactions", err)
		return nil, nil, err
	}

	return wallet, transactions, nil
}
