package scrapers

import (
	"github.com/gocolly/colly/v2"
	"log"
	"strconv"
	"strings"
)

func GetPriceInCurrencyScraper(cryptoID string, currencyOfValue string) float64 {
	var price float64
	var priceString string

	c := colly.NewCollector(
		colly.AllowedDomains("coinmarketcap.com"),
	)

	c.OnHTML(".priceValue", func(e *colly.HTMLElement) {
		priceString = e.Text
	})

	// Before make a request
	c.OnRequest(func(r *colly.Request) {
		if currencyOfValue == "btc" {
			// BTC cookie
			r.Headers.Set("Cookie", "currency=%7B%22name%22%3A%22Bitcoin%22%2C%22token%22%3A%22BTC%22%2C%22id%22%3A1%2C%22symbol%22%3A%22BTC%22%7D")
		} else if currencyOfValue == "rub" {
			// RUB cookie
			r.Headers.Set("Cookie", "currency=%7B%22id%22%3A2806%2C%22name%22%3A%22Russian%20Ruble%22%2C%22symbol%22%3A%22rub%22%2C%22token%22%3A%22%E2%82%BD%22%7D")
		} else if currencyOfValue == "eth" {
			// ETH cookie
			r.Headers.Set("Cookie", "currency=%7B%22name%22%3A%22Ethereum%22%2C%22token%22%3A%22ETH%22%2C%22id%22%3A1027%2C%22symbol%22%3A%22ETH%22%7D")
		}
	})

	c.OnScraped(func(r *colly.Response) {
		if currencyOfValue == "usd" {
			priceString = strings.ReplaceAll(priceString[:], "$", "")
			priceString = strings.ReplaceAll(priceString, ",", "")
			price, _ = strconv.ParseFloat(priceString, 64)
		} else if currencyOfValue == "rub" {
			priceString = strings.ReplaceAll(priceString[:], "₽", "")
			priceString = strings.ReplaceAll(priceString, ",", "")
			price, _ = strconv.ParseFloat(priceString, 64)
		} else if currencyOfValue == "btc" {
			priceString = strings.ReplaceAll(priceString[:], "BTC", "")
			priceString = strings.ReplaceAll(priceString, ",", "")
			price, _ = strconv.ParseFloat(priceString, 64)
		} else if currencyOfValue == "eth" {
			priceString = strings.ReplaceAll(priceString[:], "ETH", "")
			priceString = strings.ReplaceAll(priceString, ",", "")
			price, _ = strconv.ParseFloat(priceString, 64)
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Failed with response:", r, "\nError:", err)
	})

	err := c.Visit("https://coinmarketcap.com/currencies/" + cryptoID + "/")
	if err != nil {
		panic(err)
	}

	return price
}
