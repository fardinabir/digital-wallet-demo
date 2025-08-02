package controller

import (
	"net/http"

	"github.com/fardinabir/digital-wallet-demo/internal/errors"
	"github.com/fardinabir/digital-wallet-demo/internal/model"
	"github.com/fardinabir/digital-wallet-demo/internal/service"
	"github.com/labstack/echo/v4"
)

// WalletHandler is the request handler for the wallet endpoint.
type WalletHandler interface {
	Create(c echo.Context) error
	Update(c echo.Context) error
	Delete(c echo.Context) error
	Find(c echo.Context) error
	FindAll(c echo.Context) error
}

type walletHandler struct {
	Handler
	service service.Wallet
}

// NewWallet returns a new instance of the wallet handler.
func NewWallet(s service.Wallet) WalletHandler {
	return &walletHandler{service: s}
}

// CreateRequest is the request parameter for creating a new wallet
type CreateRequest struct {
	Task     string         `json:"task" validate:"required"`
	Priority model.Priority `json:"priority" validate:"required,validPriority"`
}

// adding test comments

// @Summary	Create a new wallet
// @Tags		wallets
// @Accept		json
// @Produce	json
// @Param		request	body		CreateRequest	true	"json"
// @Success	201		{object}	ResponseError{data=model.Wallet}
// @Failure	400		{object}	ResponseError
// @Failure	500		{object}	ResponseError
// @Router		/wallets [post]
func (t *walletHandler) Create(c echo.Context) error {
	var req CreateRequest
	if err := t.MustBind(c, &req); err != nil {
		return c.JSON(http.StatusBadRequest,
			ResponseError{Errors: []Error{{Code: errors.CodeBadRequest, Message: err.Error()}}})
	}

	wallet, err := t.service.Create(req.Task, req.Priority)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			ResponseError{Errors: []Error{{Code: errors.CodeInternalServerError, Message: err.Error()}}})
	}

	return c.JSON(http.StatusCreated, ResponseData{Data: wallet})
}

// UpdateRequest is the request parameter for updating a wallet
type UpdateRequest struct {
	UpdateRequestBody
	UpdateRequestPath
}

// UpdateRequestBody is the request body for updating a wallet
type UpdateRequestBody struct {
	Task     string         `json:"task,omitempty"`
	Status   model.Status   `json:"status,omitempty" validate:"validStatus"`
	Priority model.Priority `json:"priority,omitempty" validate:"validPriority"`
}

// UpdateRequestPath is the request parameter for updating a wallet
type UpdateRequestPath struct {
	ID int `param:"id" validate:"required"`
}

// @Summary	Update a wallet
// @Tags		wallets
// @Accept		json
// @Produce	json
// @Param		body	body		UpdateRequestBody	true	"body"
// @Param		path	path		UpdateRequestPath	false	"path"
// @Success	201		{object}	ResponseData{Data=model.Wallet}
// @Failure	400		{object}	ResponseError
// @Failure	500		{object}	ResponseError
// @Router		/wallets/{id} [put]
func (t *walletHandler) Update(c echo.Context) error {
	var req UpdateRequest
	if err := t.MustBind(c, &req); err != nil {
		return c.JSON(http.StatusBadRequest,
			ResponseError{Errors: []Error{{Code: errors.CodeBadRequest, Message: err.Error()}}})
	}

	wallet, err := t.service.Update(req.ID, req.Task, req.Priority, req.Status)
	if err != nil {
		if err == model.ErrNotFound {
			return c.JSON(http.StatusNotFound,
				ResponseError{Errors: []Error{{Code: errors.CodeNotFound, Message: "wallet not found"}}})
		}
		return c.JSON(http.StatusInternalServerError,
			ResponseError{Errors: []Error{{Code: errors.CodeInternalServerError, Message: err.Error()}}})
	}

	return c.JSON(http.StatusOK, ResponseData{Data: wallet})
}

// DeleteRequest is the request parameter for deleting a wallet
type DeleteRequest struct {
	ID int `param:"id" validate:"required"`
}

// @Summary	Delete a wallet
// @Tags		wallets
// @Param		path	path	DeleteRequest	false	"path"
// @Success	204
// @Failure	400	{object}	ResponseError
// @Failure	404	{object}	ResponseError
// @Failure	500	{object}	ResponseError
// @Router		/wallets/{id} [delete]
func (t *walletHandler) Delete(c echo.Context) error {
	var req DeleteRequest
	if err := t.MustBind(c, &req); err != nil {
		return c.JSON(http.StatusBadRequest,
			ResponseError{Errors: []Error{{Code: errors.CodeBadRequest, Message: err.Error()}}})
	}

	if err := t.service.Delete(req.ID); err != nil {
		if err == model.ErrNotFound {
			return c.JSON(http.StatusNotFound,
				ResponseError{Errors: []Error{{Code: errors.CodeNotFound, Message: "wallet not found"}}})
		}
		return c.JSON(http.StatusInternalServerError,
			ResponseError{Errors: []Error{{Code: errors.CodeInternalServerError, Message: err.Error()}}})
	}
	return c.NoContent(http.StatusNoContent)
}

// FindRequest is the request parameter for finding a wallet
type FindRequest struct {
	ID int `param:"id" validate:"required"`
}

// @Summary	Find a wallet
// @Tags		wallets
// @Param		path	path		FindRequest	false	"path"
// @Success	200		{object}	ResponseData{Data=model.Wallet}
// @Failure	400		{object}	ResponseError
// @Failure	404		{object}	ResponseError
// @Failure	500		{object}	ResponseError
// @Router		/wallets/{id} [get]
func (t *walletHandler) Find(c echo.Context) error {
	var req FindRequest
	if err := t.MustBind(c, &req); err != nil {
		return c.JSON(http.StatusBadRequest,
			ResponseError{Errors: []Error{{Code: errors.CodeBadRequest, Message: err.Error()}}})
	}

	res, err := t.service.Find(req.ID)
	if err != nil {
		if err == model.ErrNotFound {
			return c.JSON(http.StatusNotFound,
				ResponseError{Errors: []Error{{Code: errors.CodeNotFound, Message: "wallet not found"}}})
		}
		return c.JSON(http.StatusInternalServerError,
			ResponseError{Errors: []Error{{Code: errors.CodeInternalServerError, Message: err.Error()}}})
	}
	return c.JSON(http.StatusOK, ResponseData{Data: res})
}

// @Summary	Find all wallets
// @Tags		wallets
// @Param		task	query		string	false	"Filter by task name"
// @Param		status	query		string	false	"Filter by task status"
// @Success	200		{object}	ResponseData{Data=[]model.Wallet}
// @Failure	500		{object}	ResponseError
// @Router		/wallets [get]
func (t *walletHandler) FindAll(c echo.Context) error {
	params := c.QueryParams()
	res, err := t.service.FindAll(params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			ResponseError{Errors: []Error{{Code: errors.CodeInternalServerError, Message: err.Error()}}})
	}
	return c.JSON(http.StatusOK, ResponseData{Data: res})
}
