package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type queryString struct {
	Data string `json:"data"`
}

type jsonResponse struct {
	Id             string  `json:"id"`
	Symbol         string  `json:"symbol"`
	Name           string  `json:"name"`
	Image          string  `json:"image"`
	CurrentPrice   float32 `json:"current_price"`
	PriceChange24h float32 `json:"price_change_24h"`
	LastUpdated    string  `json:"last_update"`
}
type cashType struct {
	USD float32 `json:"usd"`
	BRL float32 `json:"brl"`
}

type coinMarket struct {
	Value cashType `json:"current_price"`
}

type coinRespose struct {
	Id      string     `json:"id"`
	Content coinMarket `json:"market_data"`
	Parcial bool
}

func main() {
	r := gin.Default()
	routes(r)
	r.Run()
}

func routes(r *gin.Engine) {
	r.GET("/myapp", func(c *gin.Context) {
		resp := queryString{c.DefaultQuery("data", "Olá")}
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
		var jsonData []jsonResponse
		error := json.Unmarshal(body, &jsonData)
		if error != nil {
			log.Fatalln(error)
		}
		c.JSON(http.StatusOK, jsonData)
	})

	r.GET("/cryptosumary", func(c *gin.Context) {
		coins := c.QueryMap("coin")
		dataChan := make(chan coinRespose)
		for _, element := range coins {
			getMarketData(element, dataChan)
		}
		var response []coinRespose
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

func getMarketData(coinId string, dataChan chan coinRespose) {
	go func() {
		url := fmt.Sprintf("https://api.coingecko.com/api/v3/coins/%s?localization=false&tickers=false&market_data=true&community_data=false&developer_data=false&sparkline=false", coinId)
		resp, err := http.Get(url)
		var m coinRespose

		if err != nil {
			stringify := fmt.Sprintf(`{"id":"%s"}`, coinId)
			json.Unmarshal([]byte(stringify), &m)
			m.Parcial = true
			dataChan <- m
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		unErr := json.Unmarshal(body, &m)
		if unErr != nil {
			log.Fatalln(unErr)
		}
		if m.Id != "" {
			m.Parcial = false
		} else {
			m.Id = coinId
			m.Parcial = true
		}
		dataChan <- m
	}()
}
