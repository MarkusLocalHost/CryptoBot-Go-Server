package collectors

import (
	"cryptocurrency/internal/models"
	"github.com/gocolly/colly/v2"
	"log"
	"strconv"
	"strings"
)

func GetTrendingFromCMC() []models.TrendingCryptocurrencyCMC {
	c := colly.NewCollector(
		colly.AllowedDomains("coinmarketcap.com"),
	)

	var cryptoData []models.TrendingCryptocurrencyCMC

	c.OnHTML("tbody tr", func(e *colly.HTMLElement) {
		cryptoCurrencyData := models.TrendingCryptocurrencyCMC{}

		e.ForEach("span", func(i int, element *colly.HTMLElement) {
			if i == 2 {
				cryptoCurrencyData.Price, _ = strconv.ParseFloat(strings.ReplaceAll(element.Text, "$", ""), 64)
			} else if i == 3 {
				cryptoCurrencyData.Price24Hour, _ = strconv.ParseFloat(strings.ReplaceAll(element.Text, "%", ""), 64)
			} else if i == 5 {
				cryptoCurrencyData.Price7D, _ = strconv.ParseFloat(strings.ReplaceAll(element.Text, "%", ""), 64)
			} else if i == 7 {
				cryptoCurrencyData.Price30D, _ = strconv.ParseFloat(strings.ReplaceAll(element.Text, "%", ""), 64)
			} else if i == 9 {
				cryptoData = append(cryptoData, cryptoCurrencyData)
			}
		})
	})

	scoreNumber := 0

	c.OnHTML("tbody tr", func(e *colly.HTMLElement) {
		e.ForEach("p", func(i int, element *colly.HTMLElement) {
			if i == 1 {
				cryptoData[scoreNumber].Name = element.Text
			} else if i == 2 {
				cryptoData[scoreNumber].Symbol = element.Text

				cryptoData[scoreNumber].Score = scoreNumber
				scoreNumber++
			}
		})
	})

	//c.OnScraped(func(r *colly.Response) {
	//	log.Println(cryptoData)
	//})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Failed with response:", r, "\nError:", err)
	})

	err := c.Visit("https://coinmarketcap.com/trending-cryptocurrencies/")
	if err != nil {
		panic(err)
	}

	return cryptoData
}

func GetTrendingFromCMCGainersAndLosers() []models.TrendingCryptocurrencyCMC {
	c := colly.NewCollector(
		colly.AllowedDomains("coinmarketcap.com"),
	)

	var cryptoData []models.TrendingCryptocurrencyCMC

	c.OnHTML("tbody tr", func(e *colly.HTMLElement) {
		cryptoCurrencyData := models.TrendingCryptocurrencyCMC{}

		e.ForEach("span", func(i int, element *colly.HTMLElement) {
			if i == 0 {
				cryptoCurrencyData.Price, _ = strconv.ParseFloat(strings.ReplaceAll(element.Text, "$", ""), 64)
			} else if i == 1 {
				cryptoCurrencyData.Price24Hour, _ = strconv.ParseFloat(strings.ReplaceAll(element.Text, "%", ""), 64)
			} else if i == 2 {
				cryptoData = append(cryptoData, cryptoCurrencyData)
			}
		})
	})

	scoreNumber := 0

	c.OnHTML("tbody tr", func(e *colly.HTMLElement) {
		e.ForEach("p", func(i int, element *colly.HTMLElement) {
			if i == 1 {
				cryptoData[scoreNumber].Name = element.Text
			} else if i == 2 {
				cryptoData[scoreNumber].Symbol = element.Text

				cryptoData[scoreNumber].Score = scoreNumber
				scoreNumber++
			}
		})
	})

	//c.OnScraped(func(r *colly.Response) {
	//	log.Println(cryptoData)
	//})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Failed with response:", r, "\nError:", err)
	})

	err := c.Visit("https://coinmarketcap.com/gainers-losers/")
	if err != nil {
		panic(err)
	}

	return cryptoData
}

func GetTrendingFromCMCMostVisited() []models.TrendingCryptocurrencyCMC {
	c := colly.NewCollector(
		colly.AllowedDomains("coinmarketcap.com"),
	)

	var cryptoData []models.TrendingCryptocurrencyCMC

	c.OnHTML("tbody tr", func(e *colly.HTMLElement) {
		cryptoCurrencyData := models.TrendingCryptocurrencyCMC{}

		e.ForEach("span", func(i int, element *colly.HTMLElement) {
			if i == 2 {
				cryptoCurrencyData.Price, _ = strconv.ParseFloat(strings.ReplaceAll(element.Text, "$", ""), 64)
			} else if i == 3 {
				cryptoCurrencyData.Price24Hour, _ = strconv.ParseFloat(strings.ReplaceAll(element.Text, "%", ""), 64)
			} else if i == 5 {
				cryptoCurrencyData.Price7D, _ = strconv.ParseFloat(strings.ReplaceAll(element.Text, "%", ""), 64)
			} else if i == 7 {
				cryptoCurrencyData.Price30D, _ = strconv.ParseFloat(strings.ReplaceAll(element.Text, "%", ""), 64)
			} else if i == 9 {
				cryptoData = append(cryptoData, cryptoCurrencyData)
			}
		})
	})

	scoreNumber := 0

	c.OnHTML("tbody tr", func(e *colly.HTMLElement) {
		e.ForEach("p", func(i int, element *colly.HTMLElement) {
			if i == 1 {
				cryptoData[scoreNumber].Name = element.Text
			} else if i == 2 {
				cryptoData[scoreNumber].Symbol = element.Text

				cryptoData[scoreNumber].Score = scoreNumber
				scoreNumber++
			}
		})
	})

	//c.OnScraped(func(r *colly.Response) {
	//	log.Println(cryptoData)
	//})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Failed with response:", r, "\nError:", err)
	})

	err := c.Visit("https://coinmarketcap.com/most-viewed-pages/")
	if err != nil {
		panic(err)
	}

	return cryptoData
}

func GetTrendingFromCMCRecentlyAdded() []models.TrendingCryptocurrencyCMC {
	c := colly.NewCollector(
		colly.AllowedDomains("coinmarketcap.com"),
	)

	var cryptoData []models.TrendingCryptocurrencyCMC

	c.OnHTML("tbody tr", func(e *colly.HTMLElement) {
		cryptoCurrencyData := models.TrendingCryptocurrencyCMC{}

		e.ForEach("span", func(i int, element *colly.HTMLElement) {
			if i == 2 {
				cryptoCurrencyData.Price, _ = strconv.ParseFloat(strings.ReplaceAll(element.Text, "$", ""), 64)
			} else if i == 3 {
				cryptoCurrencyData.Price1Hour, _ = strconv.ParseFloat(strings.ReplaceAll(element.Text, "%", ""), 64)
			} else if i == 5 {
				cryptoCurrencyData.Price24Hour, _ = strconv.ParseFloat(strings.ReplaceAll(element.Text, "%", ""), 64)
			} else if i == 7 {
				cryptoData = append(cryptoData, cryptoCurrencyData)
			}
		})
	})

	scoreNumber := 0

	c.OnHTML("tbody tr", func(e *colly.HTMLElement) {
		e.ForEach("p", func(i int, element *colly.HTMLElement) {
			if i == 1 {
				cryptoData[scoreNumber].Name = element.Text
			} else if i == 2 {
				cryptoData[scoreNumber].Symbol = element.Text

				cryptoData[scoreNumber].Score = scoreNumber
				scoreNumber++
			}
		})
	})

	//c.OnScraped(func(r *colly.Response) {
	//	log.Println(cryptoData)
	//})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Failed with response:", r, "\nError:", err)
	})

	err := c.Visit("https://coinmarketcap.com/new/")
	if err != nil {
		panic(err)
	}

	return cryptoData
}
