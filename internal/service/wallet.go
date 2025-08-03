// Package service provides the business logic for the wallet endpoint.
package service

import (
	"errors"

	"github.com/fardinabir/digital-wallet-demo/internal/model"
	"github.com/fardinabir/digital-wallet-demo/internal/repository"
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
	return t.walletRepository.Create(wallet)
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
		return nil, err
	}

	// Find or get provider wallet
	providerWallet, err := t.walletRepository.FindProviderWallet(*providerID)
	if err != nil {
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
		ObjectWalletID:  &userWallet.UserID,
		TransactionType: model.Deposit,
		OperationType:   model.Debit,
		Amount:          amountCents,
		Status:          model.Completed,
	}
	if err := t.walletRepository.CreateTransactionWithTx(tx, debitTxn); err != nil {
		tx.Rollback()
		return nil, err
	}

	// Create credit transaction for user
	creditTxn := &model.Transaction{
		SubjectWalletID: userWallet.UserID,
		ObjectWalletID:  &providerWallet.UserID,
		TransactionType: model.Deposit,
		OperationType:   model.Credit,
		Amount:          amountCents,
		Status:          model.Completed,
	}
	if err := t.walletRepository.CreateTransactionWithTx(tx, creditTxn); err != nil {
		tx.Rollback()
		return nil, err
	}

	// Update wallet balances
	if err := t.walletRepository.UpdateWalletBalance(tx, providerWallet.ID, amountCents, false); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := t.walletRepository.UpdateWalletBalance(tx, userWallet.ID, amountCents, true); err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
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
		return nil, err
	}

	// Check balance
	if userWallet.Balance < amountCents {
		return nil, model.ErrInsufficientFunds
	}

	// Find or get provider wallet
	providerWallet, err := t.walletRepository.FindProviderWallet(*providerID)
	if err != nil {
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
		ObjectWalletID:  &providerWallet.UserID,
		TransactionType: model.Withdraw,
		OperationType:   model.Debit,
		Amount:          amountCents,
		Status:          model.Completed,
	}
	if err := t.walletRepository.CreateTransactionWithTx(tx, debitTxn); err != nil {
		tx.Rollback()
		return nil, err
	}

	// Create credit transaction for provider
	creditTxn := &model.Transaction{
		SubjectWalletID: providerWallet.UserID,
		ObjectWalletID:  &userWallet.UserID,
		TransactionType: model.Withdraw,
		OperationType:   model.Credit,
		Amount:          amountCents,
		Status:          model.Completed,
	}
	if err := t.walletRepository.CreateTransactionWithTx(tx, creditTxn); err != nil {
		tx.Rollback()
		return nil, err
	}

	// Update wallet balances
	if err := t.walletRepository.UpdateWalletBalance(tx, userWallet.ID, amountCents, false); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := t.walletRepository.UpdateWalletBalance(tx, providerWallet.ID, amountCents, true); err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
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
		return nil, err
	}

	// Find receiver wallet to get UserID
	toWallet, err := t.walletRepository.FindByUserID(toUserID)
	if err != nil {
		return nil, err
	}

	// Check balance
	if fromWallet.Balance < amountCents {
		return nil, model.ErrInsufficientFunds
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
		ObjectWalletID:  &toWallet.UserID,
		TransactionType: model.Transfer,
		OperationType:   model.Debit,
		Amount:          amountCents,
		Status:          model.Completed,
	}
	if err := t.walletRepository.CreateTransactionWithTx(tx, debitTxn); err != nil {
		tx.Rollback()
		return nil, err
	}

	// Create credit transaction for receiver
	creditTxn := &model.Transaction{
		SubjectWalletID: toWallet.UserID,
		ObjectWalletID:  &fromWallet.UserID,
		TransactionType: model.Transfer,
		OperationType:   model.Credit,
		Amount:          amountCents,
		Status:          model.Completed,
	}
	if err := t.walletRepository.CreateTransactionWithTx(tx, creditTxn); err != nil {
		tx.Rollback()
		return nil, err
	}

	// Update wallet balances
	if err := t.walletRepository.UpdateWalletBalance(tx, fromWallet.ID, amountCents, false); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := t.walletRepository.UpdateWalletBalance(tx, toWallet.ID, amountCents, true); err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// Return the debit transaction for the sender
	return debitTxn, nil
}

func (t *wallet) GetWalletWithTransactions(userID string) (*model.Wallet, []model.Transaction, error) {
	// Find wallet by user ID
	wallet, err := t.walletRepository.FindByUserID(userID)
	if err != nil {
		return nil, nil, err
	}

	// Get transactions for this wallet using UserID
	filters := map[string]interface{}{
		"subject_wallet_id": wallet.UserID,
	}
	transactions, err := t.walletRepository.FindAllTransactions(filters)
	if err != nil {
		return nil, nil, err
	}

	return wallet, transactions, nil
}
