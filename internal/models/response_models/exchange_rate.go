package response_models

type ExchangeRateVariant struct {
	Exchanger     string
	ExchangeAttrs []string
	Rate          float64
	Min           float64
	Max           float64
	Reserve       float64
	GoodReviews   int
	BadReviews    int
	Link          string
}
