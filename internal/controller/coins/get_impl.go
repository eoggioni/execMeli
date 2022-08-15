package coin_controller

import (
	"exec/internal/service"

	"github.com/gin-gonic/gin"
)

type CoinController struct {
	CoinService service.CoinService
}

func NewCoinController(service service.CoinService) *CoinController {
	return &CoinController{
		CoinService: service,
	}
}

func (cr *CoinController) GetCoins(c *gin.Context) {
	coins := c.QueryMap("coin")
	resp, err := cr.CoinService.GetCoins(coins)
	if err != nil {
		return
	}
	c.JSON(200, resp)
}

func (cr *CoinController) GetCoinsSummary(c *gin.Context) {
	coins := c.QueryMap("coin")
	resp, err := cr.CoinService.GetSummaryCoins(coins)
	if err != nil {
		return
	}
	c.JSON(200, resp)
}
