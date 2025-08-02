package controller

import (
	"github.com/labstack/echo/v4"
)

func InitRoutes(api *echo.Group, controller WalletHandler) {
	wallet := api.Group("/wallets")
	{
		wallet.POST("", controller.Create)
		wallet.GET("", controller.FindAll)
		wallet.GET("/:id", controller.Find)
		wallet.PUT("/:id", controller.Update)
		wallet.DELETE("/:id", controller.Delete)
	}
}
