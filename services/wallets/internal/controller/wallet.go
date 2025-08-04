package controller

import (
	"net/http"

	"github.com/fardinabir/digital-wallet-demo/services/wallets/internal/errors"
	"github.com/fardinabir/digital-wallet-demo/services/wallets/internal/model"
	"github.com/fardinabir/digital-wallet-demo/services/wallets/internal/service"
	"github.com/labstack/echo/v4"
)

// WalletHandler is the request handler for the wallet endpoint.
type WalletHandler interface {
	Create(c echo.Context) error
	Deposit(c echo.Context) error
	Withdraw(c echo.Context) error
	Transfer(c echo.Context) error
	FetchTransactions(c echo.Context) error
}

type walletHandler struct {
	Handler
	service service.Wallet
}

// NewWalletController returns a new instance of the wallet handler.
func NewWalletController(s service.Wallet) WalletHandler {
	return &walletHandler{service: s}
}

// CreateRequest is the request parameter for creating a new wallet
type CreateRequest struct {
	UserID   string         `json:"user_id" validate:"required"`
	AcntType model.AcntType `json:"acnt_type" validate:"required,validAcntType"`
}

// DepositRequest represents the request for deposit operation
type DepositRequest struct {
	UserID     string  `json:"user_id" validate:"required"`
	Amount     int     `json:"amount" validate:"required,gt=0"`
	ProviderID *string `json:"provider_id,omitempty"`
}

// WithdrawRequest represents the request for withdraw operation
type WithdrawRequest struct {
	UserID     string  `json:"user_id" validate:"required"`
	Amount     int     `json:"amount" validate:"required,gt=0"`
	ProviderID *string `json:"provider_id,omitempty"`
}

// TransferRequest represents the request for transfer operation
type TransferRequest struct {
	FromUserID string `json:"from_user_id" validate:"required"`
	ToUserID   string `json:"to_user_id" validate:"required"`
	Amount     int    `json:"amount" validate:"required,gt=0"`
}

// WalletSummary represents essential wallet information for API responses
type WalletSummary struct {
	Balance  int64          `json:"balance"`
	AcntType model.AcntType `json:"acnt_type"`
	Status   model.Status   `json:"status"`
}

// WalletResponse represents wallet with transaction history
type WalletResponse struct {
	Wallet       WalletSummary       `json:"wallet"`
	Transactions []model.Transaction `json:"transactions"`
}

// @Summary	Create a new wallet
// @Tags		wallets
// @Accept		json
// @Produce	json
// @Param		request	body		CreateRequest	true	"json"
// @Success	201		{object}	ResponseData{data=model.Wallet}
// @Failure	400		{object}	ResponseError
// @Failure	500		{object}	ResponseError
// @Router		/wallets [post]
func (t *walletHandler) Create(c echo.Context) error {
	var req CreateRequest
	if err := t.MustBind(c, &req); err != nil {
		return c.JSON(http.StatusBadRequest,
			ResponseError{Errors: []Error{{Code: errors.CodeBadRequest, Message: err.Error()}}})
	}

	wallet := model.NewWallet(req.UserID, req.AcntType)
	if err := t.service.Create(wallet); err != nil {
		return c.JSON(http.StatusInternalServerError,
			ResponseError{Errors: []Error{{Code: errors.CodeInternalServerError, Message: err.Error()}}})
	}

	return c.JSON(http.StatusCreated, ResponseData{Data: wallet})
}

// @Summary	Deposit money to wallet
// @Tags		wallets
// @Accept		json
// @Produce	json
// @Param		request	body		DepositRequest	true	"Deposit request"
// @Success	201		{object}	ResponseData{data=model.Transaction}
// @Failure	400		{object}	ResponseError
// @Failure	404		{object}	ResponseError
// @Failure	500		{object}	ResponseError
// @Router		/wallets/deposit [post]
func (t *walletHandler) Deposit(c echo.Context) error {
	var req DepositRequest
	if err := t.MustBind(c, &req); err != nil {
		return c.JSON(http.StatusBadRequest,
			ResponseError{Errors: []Error{{Code: errors.CodeBadRequest, Message: err.Error()}}})
	}

	transaction, err := t.service.Deposit(req.UserID, req.Amount, req.ProviderID)
	if err != nil {
		if err == model.ErrNotFound {
			return c.JSON(http.StatusNotFound,
				ResponseError{Errors: []Error{{Code: errors.CodeNotFound, Message: "Wallet not found"}}})
		}
		return c.JSON(http.StatusInternalServerError,
			ResponseError{Errors: []Error{{Code: errors.CodeInternalServerError, Message: err.Error()}}})
	}

	return c.JSON(http.StatusCreated, ResponseData{Data: transaction})
}

// @Summary	Withdraw money from wallet
// @Tags		wallets
// @Accept		json
// @Produce	json
// @Param		request	body		WithdrawRequest	true	"Withdraw request"
// @Success	201		{object}	ResponseData{data=model.Transaction}
// @Failure	400		{object}	ResponseError
// @Failure	404		{object}	ResponseError
// @Failure	422		{object}	ResponseError
// @Failure	500		{object}	ResponseError
// @Router		/wallets/withdraw [post]
func (t *walletHandler) Withdraw(c echo.Context) error {
	var req WithdrawRequest
	if err := t.MustBind(c, &req); err != nil {
		return c.JSON(http.StatusBadRequest,
			ResponseError{Errors: []Error{{Code: errors.CodeBadRequest, Message: err.Error()}}})
	}

	transaction, err := t.service.Withdraw(req.UserID, req.Amount, req.ProviderID)
	if err != nil {
		if err == model.ErrNotFound {
			return c.JSON(http.StatusNotFound,
				ResponseError{Errors: []Error{{Code: errors.CodeNotFound, Message: "Wallet not found"}}})
		}
		if err == model.ErrInsufficientFunds {
			return c.JSON(http.StatusUnprocessableEntity,
				ResponseError{Errors: []Error{{Code: errors.CodeBadRequest, Message: "Insufficient balance"}}})
		}
		return c.JSON(http.StatusInternalServerError,
			ResponseError{Errors: []Error{{Code: errors.CodeInternalServerError, Message: err.Error()}}})
	}

	return c.JSON(http.StatusCreated, ResponseData{Data: transaction})
}

// @Summary	Transfer money between wallets
// @Tags		wallets
// @Accept		json
// @Produce	json
// @Param		request	body		TransferRequest	true	"Transfer request"
// @Success	201		{object}	ResponseData{data=model.Transaction}
// @Failure	400		{object}	ResponseError
// @Failure	404		{object}	ResponseError
// @Failure	422		{object}	ResponseError
// @Failure	500		{object}	ResponseError
// @Router		/wallets/transfer [post]
func (t *walletHandler) Transfer(c echo.Context) error {
	var req TransferRequest
	if err := t.MustBind(c, &req); err != nil {
		return c.JSON(http.StatusBadRequest,
			ResponseError{Errors: []Error{{Code: errors.CodeBadRequest, Message: err.Error()}}})
	}

	// Validate that from and to wallets are different
	if req.FromUserID == req.ToUserID {
		return c.JSON(http.StatusBadRequest,
			ResponseError{Errors: []Error{{Code: errors.CodeBadRequest, Message: "Cannot transfer to the same wallet"}}})
	}

	transaction, err := t.service.Transfer(req.FromUserID, req.ToUserID, req.Amount)
	if err != nil {
		if err == model.ErrNotFound {
			return c.JSON(http.StatusNotFound,
				ResponseError{Errors: []Error{{Code: errors.CodeNotFound, Message: "Wallet not found"}}})
		}
		if err == model.ErrInsufficientFunds {
			return c.JSON(http.StatusUnprocessableEntity,
				ResponseError{Errors: []Error{{Code: errors.CodeBadRequest, Message: "Insufficient balance"}}})
		}
		return c.JSON(http.StatusInternalServerError,
			ResponseError{Errors: []Error{{Code: errors.CodeInternalServerError, Message: err.Error()}}})
	}

	return c.JSON(http.StatusCreated, ResponseData{Data: transaction})
}

// FindRequest is the request parameter for finding a wallet
type FindRequest struct {
	UserID string `param:"user_id" validate:"required"`
}

// @Summary	View wallet balance & transaction history
// @Tags		wallets
// @Param		user_id	path		string	true	"User ID"
// @Success	200		{object}	ResponseData{data=WalletResponse}
// @Failure	400		{object}	ResponseError
// @Failure	404		{object}	ResponseError
// @Failure	500		{object}	ResponseError
// @Router		/wallets/{user_id} [get]
func (t *walletHandler) FetchTransactions(c echo.Context) error {
	var req FindRequest
	if err := t.MustBind(c, &req); err != nil {
		return c.JSON(http.StatusBadRequest,
			ResponseError{Errors: []Error{{Code: errors.CodeBadRequest, Message: err.Error()}}})
	}

	wallet, transactions, err := t.service.GetWalletWithTransactions(req.UserID)
	if err != nil {
		if err == model.ErrNotFound {
			return c.JSON(http.StatusNotFound,
				ResponseError{Errors: []Error{{Code: errors.CodeNotFound, Message: "wallet not found"}}})
		}
		return c.JSON(http.StatusInternalServerError,
			ResponseError{Errors: []Error{{Code: errors.CodeInternalServerError, Message: err.Error()}}})
	}

	response := WalletResponse{
		Wallet: WalletSummary{
			Balance:  wallet.Balance,
			AcntType: wallet.AcntType,
			Status:   wallet.Status,
		},
		Transactions: transactions,
	}
	return c.JSON(http.StatusOK, ResponseData{Data: response})
}
