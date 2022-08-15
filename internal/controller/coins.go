package controller

import (
	"github.com/gin-gonic/gin"
)

type CryptoController interface {
	GetCoins(c *gin.Context)
	GetCoinsSummary(c *gin.Context)
}
