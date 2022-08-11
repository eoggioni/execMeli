package basictypes

type JsonResponse struct {
	Id             string  `json:"id"`
	Symbol         string  `json:"symbol"`
	Name           string  `json:"name"`
	Image          string  `json:"image"`
	CurrentPrice   float32 `json:"current_price"`
	PriceChange24h float32 `json:"price_change_24h"`
	LastUpdated    string  `json:"last_update"`
}

type CoinRespose struct {
	Id      string     `json:"id"`
	Content coinMarket `json:"market_data"`
	Parcial bool
}

type cashType struct {
	USD float32 `json:"usd"`
	BRL float32 `json:"brl"`
}

type coinMarket struct {
	Value cashType `json:"current_price"`
}
