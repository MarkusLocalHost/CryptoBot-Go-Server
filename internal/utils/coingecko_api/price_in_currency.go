package coingecko_api

import (
	"cryptocurrency/internal/models/response_models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func GetPriceInCurrency(cryptocurrencyID string, vsCurrency string) float64 {
	URL := fmt.Sprintf("https://api.coingecko.com/api/v3/coins/markets?vs_currency=%s&ids=%s&order=market_cap_desc&per_page=100&page=1&sparkline=false",
		vsCurrency, cryptocurrencyID)

	resp, err := http.Get(URL)
	if err != nil {
		log.Fatalf("Error to fetch data: %s", err)
	}
	defer resp.Body.Close()

	var CPFOResp []response_models.CurrentPriceForObserver

	err = json.NewDecoder(resp.Body).Decode(&CPFOResp)
	if err != nil {
		log.Fatalf("Error to decode data: %s", err)
	}

	return CPFOResp[0].CurrentPrice
}
