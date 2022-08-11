package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	cryptodata "exec/internal/sections/crypto"
	basictypes "exec/pkg/types"

	"github.com/gin-gonic/gin"
)

func routes(r *gin.Engine) {
	r.GET("/myapp", func(c *gin.Context) {
		resp := basictypes.QueryString{Data: c.DefaultQuery("data", "Olá")}
		c.JSON(http.StatusOK, resp)
	})

	r.GET("/crypto", func(c *gin.Context) {
		resp, err := http.Get("https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd&ids=babybitcoin%2Cbinance-bitcoin%2Cbitcoin")
		if err != nil {
			log.Fatalln(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		var jsonData []basictypes.JsonResponse
		error := json.Unmarshal(body, &jsonData)
		if error != nil {
			log.Fatalln(error)
		}
		c.JSON(http.StatusOK, jsonData)
	})

	r.GET("/cryptosumary", func(c *gin.Context) {
		coins := c.QueryMap("coin")
		dataChan := make(chan basictypes.CoinRespose)
		for _, element := range coins {
			cryptodata.GetMarketData(element, dataChan)
		}
		var response []basictypes.CoinRespose
		status := http.StatusOK
		for i := 0; i < 3; i++ {
			dataCoin := <-dataChan
			if dataCoin.Parcial {
				status = http.StatusPartialContent
			}
			response = append(response, dataCoin)
		}
		c.JSON(status, response)
	})
}
