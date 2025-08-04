package controller

import (
	"github.com/labstack/echo/v4"
)

func InitRoutes(api *echo.Group, controller WalletHandler) {
	wallet := api.Group("/wallets")
	{
		wallet.POST("", controller.Create)
		wallet.POST("/deposit", controller.Deposit)
		wallet.POST("/withdraw", controller.Withdraw)
		wallet.POST("/transfer", controller.Transfer)
		wallet.GET("/:user_id", controller.FetchTransactions)
	}
}
