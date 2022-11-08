package response_models

type BasicCurrencyInfoShortVersion struct {
	ID            string     `json:"id"`
	Symbol        string     `json:"symbol"`
	Name          string     `json:"name"`
	MarketCapRank int        `json:"market_cap_rank"`
	MarketData    MarketData `json:"market_data"`
}

type MarketData struct {
	CurrentPrice map[string]float64 `json:"current_price"`

	High24H map[string]float64 `json:"high_24h"`
	Low24H  map[string]float64 `json:"low_24h"`

	PriceChangePercentage1HInCurrency  map[string]float64 `json:"price_change_percentage_1h_in_currency"`
	PriceChangePercentage24HInCurrency map[string]float64 `json:"price_change_percentage_24h_in_currency"`
	PriceChangePercentage7DInCurrency  map[string]float64 `json:"price_change_percentage_7d_in_currency"`
	PriceChangePercentage14DInCurrency map[string]float64 `json:"price_change_percentage_14d_in_currency"`
	PriceChangePercentage30DInCurrency map[string]float64 `json:"price_change_percentage_30d_in_currency"`
	PriceChangePercentage1YInCurrency  map[string]float64 `json:"price_change_percentage_1y_in_currency"`
}
