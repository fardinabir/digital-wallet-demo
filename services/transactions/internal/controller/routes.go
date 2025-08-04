package controller

import (
	"github.com/labstack/echo/v4"
)

func InitRoutes(api *echo.Group, controller TransactionHandler) {
	transactions := api.Group("/transactions")
	{
		transactions.POST("", controller.CreateTransactionPair)
		transactions.GET("/:subject_wallet_id", controller.GetTransactions)
	}
}
