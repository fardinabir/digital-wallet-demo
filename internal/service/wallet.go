// Package service provides the business logic for the wallet endpoint.
package service

import (
	"net/url"

	"github.com/fardinabir/digital-wallet-demo/internal/model"
	"github.com/fardinabir/digital-wallet-demo/internal/repository"
)

// Wallet is the service for the wallet endpoint.
type Wallet interface {
	Create(task string, priority model.Priority) (*model.Wallet, error)
	Update(id int, task string, priority model.Priority, status model.Status) (*model.Wallet, error)
	Delete(id int) error
	Find(id int) (*model.Wallet, error)
	FindAll(qry url.Values) ([]*model.Wallet, error)
}

type wallet struct {
	walletRepository repository.Wallet
}

// NewWallet creates a new Wallet service.
func NewWallet(r repository.Wallet) Wallet {
	return &wallet{r}
}

func (t *wallet) Create(task string, priority model.Priority) (*model.Wallet, error) {
	wallet := model.NewWallet(task, priority)
	if err := t.walletRepository.Create(wallet); err != nil {
		return nil, err
	}
	return wallet, nil
}

func (t *wallet) Update(id int, task string, priority model.Priority, status model.Status) (*model.Wallet, error) {
	wallet := model.NewUpdateWallet(id, task, priority, status)
	// 現在の値を取得
	currentWallet, err := t.Find(id)
	if err != nil {
		return nil, err
	}
	// 空文字列の場合、現在の値を使用
	if wallet.Task == "" {
		wallet.Task = currentWallet.Task
	}
	if wallet.Status == "" {
		wallet.Status = currentWallet.Status
	}
	if wallet.Priority == 0 {
		wallet.Priority = currentWallet.Priority
	}
	if err := t.walletRepository.Update(wallet); err != nil {
		return nil, err
	}
	return wallet, nil
}

func (t *wallet) Delete(id int) error {
	if err := t.walletRepository.Delete(id); err != nil {
		return err
	}
	return nil
}

func (t *wallet) Find(id int) (*model.Wallet, error) {
	wallet, err := t.walletRepository.Find(id)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}

func (t *wallet) FindAll(qry url.Values) ([]*model.Wallet, error) {
	processedQry := map[string]interface{}{}
	if val, ok := qry["task"]; ok {
		processedQry["task"] = val[0]
	}
	if val, ok := qry["status"]; ok {
		processedQry["status"] = val[0]
	}
	wallet, err := t.walletRepository.FindAll(processedQry)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}
