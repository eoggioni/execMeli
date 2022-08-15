package service

type GetCoinResponse struct {
	Id          string     `json:"id"`
	Symbol      string     `json:"symbol,omitempty"`
	Name        string     `json:"name,omitempty"`
	Description string     `json:"description,omitempty"`
	Images      []string   `json:"image,omitempty"`
	Price       coinMarket `json:"market_data,omitempty"`
	Partial     bool       `json:"partial"`
}

type GetCoinSummaryResponse struct {
	Id      string     `json:"id"`
	Name    string     `json:"name,omitempty"`
	Price   coinMarket `json:"market_data,omitempty"`
	Partial bool       `json:"partial"`
}

type cashType struct {
	USD float32 `json:"usd"`
	BRL float32 `json:"brl"`
}

type coinMarket struct {
	Value cashType `json:"current_price"`
}

type CoinService interface {
	GetCoins(ids map[string]string) ([]GetCoinSummaryResponse, error)
	GetSummaryCoins(ids map[string]string) ([]GetCoinSummaryResponse, error)
}
