package response_models

type CurrentPriceForObserver struct {
	ID            string  `json:"id"`
	Symbol        string  `json:"symbol"`
	Name          string  `json:"name"`
	MarketCapRank int     `json:"market_cap_rank"`
	CurrentPrice  float64 `json:"current_price"`
}
