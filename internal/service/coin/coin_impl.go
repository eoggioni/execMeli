package coin_service

import (
	"encoding/json"
	"exec/internal/service"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

type CoinService struct {
	Url string
}

type HttpResponse struct {
	StatusCode int
	Body       []byte
}

func NewCoinService(url string) *CoinService {
	return &CoinService{
		Url: url,
	}
}

func (c *CoinService) dataAccess(id string) (*HttpResponse, error) {
	endpoint := fmt.Sprintf(c.Url+"%s?localization=false&tickers=false&market_data=true&community_data=false&developer_data=false&sparkline=false", id)
	resp, err := http.Get(endpoint)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return &HttpResponse{
		Body:       body,
		StatusCode: resp.StatusCode,
	}, nil
}

func (c *CoinService) GetCoins(ids map[string]string) ([]service.GetCoinSummaryResponse, error) {
	concurrency := len(ids)
	channel := make(chan service.GetCoinSummaryResponse, concurrency)
	var coinResponses []service.GetCoinSummaryResponse
	var wg sync.WaitGroup
	wg.Add(concurrency)

	for _, id := range ids {
		go c.worker(id, true, channel, &wg)
	}
	wg.Wait()
	close(channel)

	for c := range channel {
		coinResponses = append(coinResponses, c)
	}
	return coinResponses, nil
}

func (s *CoinService) GetSummaryCoins(ids map[string]string) ([]service.GetCoinSummaryResponse, error) {
	concurrency := len(ids)
	channel := make(chan service.GetCoinSummaryResponse, concurrency)
	var coinResponses []service.GetCoinSummaryResponse
	var wg sync.WaitGroup
	wg.Add(concurrency)

	for _, id := range ids {
		go s.worker(id, true, channel, &wg)
	}
	wg.Wait()
	close(channel)

	for i := 0; i < concurrency; i++ {
		dataCoin := <-channel
		coinResponses = append(coinResponses, dataCoin)
	}
	return coinResponses, nil
}

func (s *CoinService) worker(id string, moreData bool, c chan<- service.GetCoinSummaryResponse, wg *sync.WaitGroup) {
	defer wg.Done()
	defer s.Recover(id, c)
	coinData, err := s.dataAccess(id)
	if err != nil {
		return
	}
	if moreData {
		var d service.GetCoinResponse
		unErr := json.Unmarshal(coinData.Body, &d)
		if unErr != nil {
			return
		}
		// c <- d
	} else {
		var d service.GetCoinSummaryResponse
		unErr := json.Unmarshal(coinData.Body, &d)
		if unErr != nil {
			return
		}
		c <- d
	}
}

func (s *CoinService) Recover(id string, c chan<- service.GetCoinSummaryResponse) {
	if r := recover(); r != nil {
		c <- service.GetCoinSummaryResponse{
			Id:      id,
			Partial: true,
		}
	}
}
