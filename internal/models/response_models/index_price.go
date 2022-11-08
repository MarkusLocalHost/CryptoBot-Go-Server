package response_models

type CurrencyInfoIndexVersion struct {
	ID            string  `json:"id"`
	Symbol        string  `json:"symbol"`
	Name          string  `json:"name"`
	Image         string  `json:"image"`
	CurrentPrice  float64 `json:"current_price"`
	MarketCapRank int     `json:"market_cap_rank"`

	PriceChangeInCurrency1h  float64 `json:"price_change_percentage_1h_in_currency"`
	PriceChangeInCurrency24h float64 `json:"price_change_percentage_24h_in_currency"`
	PriceChangeInCurrency7d  float64 `json:"price_change_percentage_7d_in_currency"`
	PriceChangeInCurrency30d float64 `json:"price_change_percentage_30d_in_currency"`
	PriceChangeInCurrency1y  float64 `json:"price_change_percentage_1y_in_currency"`
}

type CurrencyInfoIndexShortVersionForPortfolio struct {
	ID           string  `json:"id"`
	Symbol       string  `json:"symbol"`
	Name         string  `json:"name"`
	CurrentPrice float64 `json:"current_price"`
}
