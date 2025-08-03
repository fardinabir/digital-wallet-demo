// Package service provides the business logic for the wallet endpoint.
package service

import (
	"errors"
	"strconv"

	"github.com/fardinabir/digital-wallet-demo/internal/model"
	"github.com/fardinabir/digital-wallet-demo/internal/repository"
)

// Wallet is the service for the wallet endpoint.
type Wallet interface {
	Create(wallet *model.Wallet) error
	Deposit(walletID int, amount string, providerID *int) (*model.Transaction, error)
	Withdraw(walletID int, amount string, providerID *int) (*model.Transaction, error)
	Transfer(fromWalletID int, toWalletID int, amount string) (*model.Transaction, error)
	GetWalletWithTransactions(userID int) (*model.Wallet, []model.Transaction, error)
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

func (t *wallet) Deposit(walletID int, amount string, providerID *int) (*model.Transaction, error) {
	// Parse amount
	amountFloat, err := strconv.ParseFloat(amount, 64)
	if err != nil || amountFloat <= 0 {
		return nil, errors.New("invalid amount")
	}
	amountCents := int64(amountFloat * 100)

	// Find user wallet
	userWallet, err := t.walletRepository.Find(walletID)
	if err != nil {
		return nil, err
	}

	// Find or get provider wallet
	providerWallet, err := t.walletRepository.FindProviderWallet(*providerID)
	if err != nil {
		return nil, errors.New("provider wallet not found")
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
		SubjectWalletID: providerWallet.ID,
		ObjectWalletID:  &userWallet.ID,
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
		SubjectWalletID: userWallet.ID,
		ObjectWalletID:  &providerWallet.ID,
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

func (t *wallet) Withdraw(walletID int, amount string, providerID *int) (*model.Transaction, error) {
	// Parse amount
	amountFloat, err := strconv.ParseFloat(amount, 64)
	if err != nil || amountFloat <= 0 {
		return nil, errors.New("invalid amount")
	}
	amountCents := int64(amountFloat * 100)

	// Find user wallet
	userWallet, err := t.walletRepository.Find(walletID)
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
		return nil, errors.New("provider wallet not found")
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
		SubjectWalletID: userWallet.ID,
		ObjectWalletID:  &providerWallet.ID,
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
		SubjectWalletID: providerWallet.ID,
		ObjectWalletID:  &userWallet.ID,
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

func (t *wallet) Transfer(fromWalletID int, toWalletID int, amount string) (*model.Transaction, error) {
	// Parse amount
	amountFloat, err := strconv.ParseFloat(amount, 64)
	if err != nil || amountFloat <= 0 {
		return nil, errors.New("invalid amount")
	}
	amountCents := int64(amountFloat * 100)

	// Find sender wallet to check balance
	fromWallet, err := t.walletRepository.Find(fromWalletID)
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
		SubjectWalletID: fromWalletID,
		ObjectWalletID:  &toWalletID,
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
		SubjectWalletID: toWalletID,
		ObjectWalletID:  &fromWalletID,
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
	if err := t.walletRepository.UpdateWalletBalance(tx, fromWalletID, amountCents, false); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := t.walletRepository.UpdateWalletBalance(tx, toWalletID, amountCents, true); err != nil {
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

func (t *wallet) GetWalletWithTransactions(userID int) (*model.Wallet, []model.Transaction, error) {
	// Find wallet by user ID
	wallet, err := t.walletRepository.FindByUserID(userID)
	if err != nil {
		return nil, nil, err
	}

	// Get transactions for this wallet
	filters := map[string]interface{}{
		"wallet_id": wallet.ID,
	}
	transactions, err := t.walletRepository.FindAllTransactions(filters)
	if err != nil {
		return nil, nil, err
	}

	return wallet, transactions, nil
}
