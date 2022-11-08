package models

type TrendingCryptocurrenciesResponse struct {
	Coins []TrendingCryptocurrency `json:"coins"`
}

type TrendingCryptocurrency struct {
	Item TrendingCryptocurrencyBody `json:"item"`
}

type TrendingCryptocurrencyBody struct {
	ID            string  `json:"id"`
	CoinID        int     `json:"coin_id"`
	Name          string  `json:"name"`
	Symbol        string  `json:"symbol"`
	MarketCapRank int     `json:"market_cap_rank"`
	Thumb         string  `json:"thumb"`
	Slug          string  `json:"slug"`
	PriceBTC      float64 `json:"price_btc"`
	Score         int     `json:"score"`
}

type CryptocurrencyBasicInfoResponse struct {
	ID            string  `json:"id"`
	Symbol        string  `json:"symbol"`
	Name          string  `json:"name"`
	MarketCapRank int     `json:"market_cap_rank"`
	Thumb         string  `json:"thumb"`
	Slug          string  `json:"slug"`
	PriceBTC      float64 `json:"price_btc"`
	Score         int     `json:"score"`

	LastUpdated string `json:"last_updated"`
}

type TrendingCryptocurrencyCMC struct {
	//ID            string  `json:"id"`
	//CoinID        int     `json:"coin_id"`
	Name          string `json:"name"`
	Symbol        string `json:"symbol"`
	MarketCapRank int    `json:"market_cap_rank"`
	//Thumb         string  `json:"thumb"`
	//Slug          string  `json:"slug"`
	Price       float64 `json:"price"`
	Score       int     `json:"score"`
	Price1Hour  float64
	Price24Hour float64
	Price7D     float64
	Price30D    float64
}
