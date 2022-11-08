package mongo_models

type TryFindCurrency struct {
	Coins []Coin `json:"coins"`
}

type Coin struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Symbol        string `json:"symbol"`
	MarketCapRank int    `json:"market_cap_rank"`
	Thumb         string `json:"thumb"`
	Large         string `json:"large"`
}
