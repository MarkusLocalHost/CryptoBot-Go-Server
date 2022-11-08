package response_models

type TrendingCurrency struct {
	ID     string `json:"id"`
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
}
