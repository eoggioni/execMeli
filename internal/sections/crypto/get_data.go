package cryptodata

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	basictypes "exec/pkg/types"
)

func GetMarketData(coinId string, dataChan chan basictypes.CoinRespose) {
	go func() {
		url := fmt.Sprintf("https://api.coingecko.com/api/v3/coins/%s?localization=false&tickers=false&market_data=true&community_data=false&developer_data=false&sparkline=false", coinId)
		resp, err := http.Get(url)
		var m basictypes.CoinRespose

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
