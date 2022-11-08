package response_models

type BasicCurrencyInfoFullVersion struct {
	ID            string                `json:"id"`
	Symbol        string                `json:"symbol"`
	Name          string                `json:"name"`
	MarketCapRank int                   `json:"market_cap_rank"`
	MarketData    MarketDataFullVersion `json:"market_data"`
}

type MarketDataFullVersion struct {
	CurrentPrice map[string]float64 `json:"current_price"`

	High24H map[string]float64 `json:"high_24h"`
	Low24H  map[string]float64 `json:"low_24h"`

	PriceChangePercentage1HInCurrency  map[string]float64 `json:"price_change_percentage_1h_in_currency"`
	PriceChangePercentage24HInCurrency map[string]float64 `json:"price_change_percentage_24h_in_currency"`
	PriceChangePercentage7DInCurrency  map[string]float64 `json:"price_change_percentage_7d_in_currency"`
	PriceChangePercentage14DInCurrency map[string]float64 `json:"price_change_percentage_14d_in_currency"`
	PriceChangePercentage30DInCurrency map[string]float64 `json:"price_change_percentage_30d_in_currency"`
	PriceChangePercentage1YInCurrency  map[string]float64 `json:"price_change_percentage_1y_in_currency"`

	CurrentCapitalization map[string]float64 `json:"market_cap"`
	FDV                   map[string]float64 `json:"fully_diluted_valuation"`

	MarketCapitalizationChange24InCurrency   float64 `json:"market_cap_change_24"`
	MarketCapitalizationChange24InPercentage float64 `json:"market_cap_change_percentage_24h"`

	TotalVolume map[string]float64 `json:"total_volume"`

	Sparkline7D map[string][]float64 `json:"sparkline_7d"`
}
