package main

import (
	"net/http"

	coin_controller "exec/internal/controller/coins"
	coin_service "exec/internal/service/coin"
	basictypes "exec/pkg/types"

	"github.com/gin-gonic/gin"
)

func routes(r *gin.Engine) {
	r.GET("/myapp", func(c *gin.Context) {
		resp := basictypes.QueryString{Data: c.DefaultQuery("data", "Olá")}
		c.JSON(http.StatusOK, resp)
	})
	coinService := coin_service.NewCoinService("https://api.coingecko.com/api/v3/coins/")
	coinController := coin_controller.NewCoinController(coinService)
	r.GET("/crypto", coinController.GetCoins)
	r.GET("/summary", coinController.GetCoinsSummary)
}
