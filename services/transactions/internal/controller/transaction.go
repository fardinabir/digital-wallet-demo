package controller

import (
	"net/http"

	"github.com/fardinabir/digital-wallet-demo/services/transactions/internal/errors"
	"github.com/fardinabir/digital-wallet-demo/services/transactions/internal/model"
	"github.com/fardinabir/digital-wallet-demo/services/transactions/internal/service"
	"github.com/labstack/echo/v4"
)

// TransactionHandler is the request handler for the transaction endpoint.
type TransactionHandler interface {
	CreateTransactionPair(c echo.Context) error
	GetTransactions(c echo.Context) error
}

type transactionHandler struct {
	Handler
	service service.TransactionService
}

// NewTransactionHandler returns a new instance of the transaction handler.
func NewTransactionHandler(s service.TransactionService) TransactionHandler {
	return &transactionHandler{service: s}
}

// TransactionPairRequest represents the request for creating a transaction pair
type TransactionPairRequest struct {
	DebitTransaction  TransactionRequest `json:"debit_transaction" validate:"required"`
	CreditTransaction TransactionRequest `json:"credit_transaction" validate:"required"`
}

// TransactionRequest represents a single transaction in the request
type TransactionRequest struct {
	SubjectWalletID string                  `json:"subject_wallet_id" validate:"required"`
	ObjectWalletID  string                  `json:"object_wallet_id" validate:"required"`
	TransactionType model.TransactionType   `json:"transaction_type" validate:"required"`
	OperationType   model.OperationType     `json:"operation_type" validate:"required"`
	Amount          int64                   `json:"amount" validate:"required,gt=0"`
	Status          model.TransactionStatus `json:"status" validate:"required"`
}

// GetTransactionsRequest represents the request for getting transactions
type GetTransactionsRequest struct {
	SubjectWalletID string `param:"subject_wallet_id" validate:"required"`
}

// @Summary	Create a transaction pair (debit and credit)
// @Tags		transactions
// @Accept		json
// @Produce	json
// @Param		request	body		TransactionPairRequest	true	"Transaction pair request"
// @Success	201		{object}	ResponseData{data=string}
// @Failure	400		{object}	ResponseError
// @Failure	500		{object}	ResponseError
// @Router		/transactions [post]
func (h *transactionHandler) CreateTransactionPair(c echo.Context) error {
	var req TransactionPairRequest
	if err := h.MustBind(c, &req); err != nil {
		return c.JSON(http.StatusBadRequest,
			ResponseError{Errors: []Error{{Code: errors.CodeBadRequest, Message: err.Error()}}})
	}

	// Convert request to model transactions
	debitTxn := &model.Transaction{
		SubjectWalletID: req.DebitTransaction.SubjectWalletID,
		ObjectWalletID:  req.DebitTransaction.ObjectWalletID,
		TransactionType: req.DebitTransaction.TransactionType,
		OperationType:   req.DebitTransaction.OperationType,
		Amount:          req.DebitTransaction.Amount,
		Status:          req.DebitTransaction.Status,
	}

	creditTxn := &model.Transaction{
		SubjectWalletID: req.CreditTransaction.SubjectWalletID,
		ObjectWalletID:  req.CreditTransaction.ObjectWalletID,
		TransactionType: req.CreditTransaction.TransactionType,
		OperationType:   req.CreditTransaction.OperationType,
		Amount:          req.CreditTransaction.Amount,
		Status:          req.CreditTransaction.Status,
	}

	// Create transaction pair
	if err := h.service.CreateTransactionPair(debitTxn, creditTxn); err != nil {
		return c.JSON(http.StatusInternalServerError,
			ResponseError{Errors: []Error{{Code: errors.CodeInternalServerError, Message: err.Error()}}})
	}

	return c.JSON(http.StatusCreated, ResponseData{Data: "Transaction pair created successfully"})
}

// @Summary	Get transactions for a wallet
// @Tags		transactions
// @Produce	json
// @Param		subject_wallet_id	path		string	true	"Subject Wallet ID"
// @Success	200					{object}	ResponseData{data=[]model.Transaction}
// @Failure	400					{object}	ResponseError
// @Failure	500					{object}	ResponseError
// @Router		/transactions/{subject_wallet_id} [get]
func (h *transactionHandler) GetTransactions(c echo.Context) error {
	var req GetTransactionsRequest
	if err := h.MustBind(c, &req); err != nil {
		return c.JSON(http.StatusBadRequest,
			ResponseError{Errors: []Error{{Code: errors.CodeBadRequest, Message: err.Error()}}})
	}

	transactions, err := h.service.GetTransactions(req.SubjectWalletID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			ResponseError{Errors: []Error{{Code: errors.CodeInternalServerError, Message: err.Error()}}})
	}

	return c.JSON(http.StatusOK, ResponseData{Data: transactions})
}
