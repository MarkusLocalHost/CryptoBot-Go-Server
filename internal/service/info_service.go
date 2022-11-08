package service

import (
	"context"
	"cryptocurrency/internal/collectors"
	"cryptocurrency/internal/models"
	"cryptocurrency/internal/models/mongo_models"
	"cryptocurrency/internal/models/response_models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type infoService struct {
	InfoRepository models.InfoRepository
}

type ISConfig struct {
	InfoRepository models.InfoRepository
}

func NewInfoService(c *ISConfig) models.InfoService {
	return &infoService{
		InfoRepository: c.InfoRepository,
	}
}

func (i infoService) GetAllUsersLanguages(ctx context.Context) (map[int64]string, error) {
	usersLanguages, err := i.InfoRepository.GetAllUsersLanguages(ctx)
	if err != nil {
		return nil, err
	}

	return usersLanguages, nil
}

func (i infoService) GetAllUsersAdmins(ctx context.Context) ([]int64, error) {
	usersAdmins, err := i.InfoRepository.GetAllUsersAdmins(ctx)
	if err != nil {
		return nil, err
	}

	return usersAdmins, nil
}

func (i infoService) GetBasicCurrencyInfo(ctx context.Context, currencyID string) (response_models.BasicCurrencyInfoFullVersion, error) {
	URL := "https://api.coingecko.com/api/v3/coins/" + currencyID

	resp, err := http.Get(URL)
	if err != nil {
		log.Fatalf("Error to fetch data: %s", err)
	}
	defer resp.Body.Close()

	var TCFGResp response_models.BasicCurrencyInfoFullVersion

	err = json.NewDecoder(resp.Body).Decode(&TCFGResp)
	if err != nil {
		return response_models.BasicCurrencyInfoFullVersion{}, err
	}

	return TCFGResp, nil
}

func (i infoService) GetBasicCurrencyInfoShortVersion(ctx context.Context, currencyID string) (response_models.BasicCurrencyInfoShortVersion, error) {
	var BSISVResp response_models.BasicCurrencyInfoShortVersion

	//log.Println(currencySymbol)
	//currencyId, err := i.InfoRepository.ConvertSymbolToIDForCoinGecko(ctx, currencySymbol)
	//log.Println(currencyId)
	//if err != nil {
	//	switch err {
	//	case mongo.ErrNoDocuments:
	//		link := "https://coinmarketcap.com/currencies/" + strings.ToLower(currencySymbol)
	//		c := colly.NewCollector(
	//			colly.AllowedDomains("coinmarketcap.com"),
	//		)
	//
	//		var price string
	//
	//		c.OnHTML(".priceValue", func(e *colly.HTMLElement) {
	//			price = e.ChildText("span")
	//		})
	//
	//		c.OnError(func(r *colly.Response, err error) {
	//			log.Println("Failed with response:", r, "\nError:", err)
	//		})
	//
	//		err := c.Visit(link)
	//		if err != nil {
	//			return response_models.BasicCurrencyInfoShortVersion{}, err
	//		}
	//
	//		BSISVResp.MarketData.CurrentPrice = make(map[string]float64)
	//		BSISVResp.MarketData.CurrentPrice["usd"], _ = strconv.ParseFloat(price, 64)
	//		BSISVResp.MarketData.CurrentPrice["rub"] = 1.5
	//
	//		return BSISVResp, nil
	//	default:
	//		return response_models.BasicCurrencyInfoShortVersion{}, err
	//	}
	//}

	URL := "https://api.coingecko.com/api/v3/coins/" + currencyID

	resp, err := http.Get(URL)
	if err != nil {
		log.Fatalf("Error to fetch data: %s", err)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&BSISVResp)
	if err != nil {
		return response_models.BasicCurrencyInfoShortVersion{}, err
	}

	return BSISVResp, nil
}

func (i infoService) GetTrendingCurrencies(ctx context.Context, source string) ([]response_models.TrendingCurrency, error) {
	var cryptocurrencies []models.TrendingCryptocurrencyCMC
	var data []response_models.TrendingCurrency

	switch source {
	case "coingecko_trending":
		URL := "https://api.coingecko.com/api/v3/search/trending"

		resp, err := http.Get(URL)
		if err != nil {
			log.Fatalf("Error to fetch data: %s", err)
		}
		defer resp.Body.Close()

		var TCFGResp models.TrendingCryptocurrenciesResponse

		err = json.NewDecoder(resp.Body).Decode(&TCFGResp)
		if err != nil {
			return nil, err
		}

		for _, currency := range TCFGResp.Coins {
			data = append(data, response_models.TrendingCurrency{
				ID:     currency.Item.ID,
				Symbol: currency.Item.Symbol,
				Name:   currency.Item.Name,
			})
		}
		log.Println(data)
	case "coinmarketcap_trending":
		cryptocurrencies = collectors.GetTrendingFromCMC()
	case "coinmarketcap_gainers":
		cryptocurrencies = collectors.GetTrendingFromCMCGainersAndLosers()
	case "coinmarketcap_losers":
		cryptocurrencies = collectors.GetTrendingFromCMCGainersAndLosers()
	case "coinmarketcap_most_visited":
		cryptocurrencies = collectors.GetTrendingFromCMCMostVisited()
	case "coinmarketcap_recently_added":
		cryptocurrencies = collectors.GetTrendingFromCMCRecentlyAdded()
	}

	if source != "coingecko_trending" {
		for num, currency := range cryptocurrencies {
			if num >= 30 && source == "coinmarketcap_gainers" {
				break
			}
			if num < 30 && source == "coinmarketcap_losers" {
				continue
			}

			idString := strings.ReplaceAll(currency.Name, " ", "-")
			idString = strings.ToLower(idString)

			data = append(data, response_models.TrendingCurrency{
				ID:     idString,
				Symbol: currency.Symbol,
				Name:   currency.Name,
			})
		}
	}

	return data, nil
}

func (i infoService) TryFindCurrencyByNameOrSlug(ctx context.Context, name string, slug string) (mongo_models.TryFindCurrency, error) {
	//var data []mongo_models.TryFindCurrency
	//var err error
	//
	//if name != "" {
	//	data, err = i.InfoRepository.TryFindCurrency(ctx, strings.ToLower(name))
	//	if err != nil {
	//		log.Fatalf("Error to find data: %s", err)
	//	}
	//} else if slug != "" {
	//	data, err = i.InfoRepository.TryFindCurrencyBySlug(ctx, strings.ToUpper(slug))
	//	if err != nil {
	//		log.Fatalf("Error to find data: %s", err)
	//	}
	//}
	//
	//return data, nil

	// Get from CoinGecko
	var URL string
	if name != "" {
		URL = "https://api.coingecko.com/api/v3/search?query=" + name
	} else if slug != "" {
		URL = "https://api.coingecko.com/api/v3/search?query=" + slug
	}

	resp, err := http.Get(URL)
	if err != nil {
		log.Fatalf("Error to fetch data: %s", err)
	}
	defer resp.Body.Close()

	var data mongo_models.TryFindCurrency

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return mongo_models.TryFindCurrency{}, err
	}

	return data, nil
}

func (i infoService) GetSupportedVSCurrencies(ctx context.Context) ([]string, error) {
	data, err := i.InfoRepository.GetSupportedVSCurrencies(ctx)
	if err != nil {
		log.Fatalf("Error to find data: %s", err)
	}

	return data, nil
}

func (i infoService) GetIndexPriceByPage(ctx context.Context, page int, currency string) ([]response_models.CurrencyInfoIndexVersion, error) {
	// Get from CoinGecko
	URL := fmt.Sprintf("https://api.coingecko.com/api/v3/coins/markets?vs_currency=%s&order=market_cap_desc&per_page=250&page=%v&sparkline=false&price_change_percentage=1h,24h,7d,30d,1y",
		strings.ToLower(currency),
		page,
	)

	resp, err := http.Get(URL)
	if err != nil {
		log.Fatalf("Error to fetch data: %s", err)
	}
	defer resp.Body.Close()

	var CIIResp []response_models.CurrencyInfoIndexVersion

	err = json.NewDecoder(resp.Body).Decode(&CIIResp)
	if err != nil {
		return []response_models.CurrencyInfoIndexVersion{}, err
	}

	return CIIResp, nil
}

func (i infoService) GetPriceForSymbolCurrenciesInPortfolio(portfolio []*mongo_models.Portfolio) ([]*mongo_models.Portfolio, error) {
	// Collect symbols
	// todo: problem with binance and heco chains
	var currencySymbols []string
	for _, currencySymbol := range portfolio {
		if contains(currencySymbols, currencySymbol.Cryptocurrency) {
			continue
		}
		currencySymbols = append(currencySymbols, currencySymbol.Cryptocurrency)
	}

	// Get from CoinGecko id and symbols
	UrlIdSymbols := "https://api.coingecko.com/api/v3/coins/list"

	respIdSymbols, err := http.Get(UrlIdSymbols)
	if err != nil {
		log.Fatalf("Error to fetch data: %s", err)
	}
	defer respIdSymbols.Body.Close()

	type CurrencyIdAndSymbolsResp struct {
		ID     string `json:"id"`
		Symbol string `json:"symbol"`
		Name   string `json:"name"`
	}
	var CIASResp []CurrencyIdAndSymbolsResp

	err = json.NewDecoder(respIdSymbols.Body).Decode(&CIASResp)
	if err != nil {
		return []*mongo_models.Portfolio{}, err
	}

	// Convert from symbols to id
	var currencyIds []string
	for _, currency := range CIASResp {
		if contains(currencySymbols, strings.ToUpper(currency.Symbol)) {
			currencyIds = append(currencyIds, currency.ID)
		}
	}

	// Get from CoinGecko prices
	URL := fmt.Sprintf("https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd&ids=%s&order=market_cap_desc&per_page=250&page=1&sparkline=false",
		strings.Join(currencyIds, ","),
	)

	resp, err := http.Get(URL)
	if err != nil {
		log.Fatalf("Error to fetch data: %s", err)
	}
	defer resp.Body.Close()

	var CIISVFPResp []response_models.CurrencyInfoIndexShortVersionForPortfolio

	err = json.NewDecoder(resp.Body).Decode(&CIISVFPResp)
	if err != nil {
		return []*mongo_models.Portfolio{}, err
	}

	// Append to Portfolio actual price
	for _, currencyPortfolio := range portfolio {
		for _, currencyResp := range CIISVFPResp {
			if currencyPortfolio.Price != 0 {
				if currencyPortfolio.Cryptocurrency == strings.ToUpper(currencyResp.Symbol) {
					currencyPortfolio.ActualPrice = currencyResp.CurrentPrice
				}
			}
		}
	}

	return portfolio, nil
}

func (i infoService) GetExchangeRateFromBestchange(ctx context.Context, from string, fromType string, to string, toType string, limitCurrency string, limitValue float64) (map[string][]response_models.ExchangeRateVariant, error) {
	results := make(map[string][]response_models.ExchangeRateVariant)

	variants, merchants := collectors.GetExchangeRateVariants(from, to, fromType)

	for key := range variants {
		for _, variant := range variants[key] {
			rateTo, err := strconv.ParseFloat(variant.RateTo, 64)
			if err != nil {
				log.Fatal(err)
			}
			rateFrom, err := strconv.ParseFloat(variant.RateFrom, 64)
			if err != nil {
				log.Fatal(err)
			}
			rate := rateTo * rateFrom

			min, err := strconv.ParseFloat(variant.Min, 64)
			if err != nil {
				log.Fatal(err)
			}

			max, err := strconv.ParseFloat(variant.Max, 64)
			if err != nil {
				log.Fatal(err)
			}

			reserve, err := strconv.ParseFloat(variant.Reserve, 64)
			if err != nil {
				log.Fatal(err)
			}

			goodReviews, err := strconv.Atoi(variant.GoodReviews)
			if err != nil {
				log.Fatal(err)
			}

			badReviews, err := strconv.Atoi(variant.BadReviews)
			if err != nil {
				log.Fatal(err)
			}

			exchangeVariant := response_models.ExchangeRateVariant{
				Exchanger:     merchants[variant.MerchantID],
				ExchangeAttrs: nil,
				Rate:          rate,
				Min:           min,
				Max:           max,
				Reserve:       reserve,
				GoodReviews:   goodReviews,
				BadReviews:    badReviews,
				Link: "https://www.bestchange.ru/click.php?id=" + variant.MerchantID +
					"&from=" + from + "&to=" + to + "&city=0",
			}

			if fromType == "bank" || fromType == "wallet" {
				if limitCurrency == "rub" {
					if limitValue >= exchangeVariant.Min && limitValue <= exchangeVariant.Max {
						results[key] = append(results[key], exchangeVariant)
					}
				} else if limitCurrency == "crypto" {
					if limitValue*exchangeVariant.Rate >= exchangeVariant.Min && limitValue*exchangeVariant.Rate <= exchangeVariant.Max {
						results[key] = append(results[key], exchangeVariant)
					}
				}
			} else if fromType == "crypto" {
				if limitCurrency == "rub" {
					if limitValue/exchangeVariant.Rate >= exchangeVariant.Min && limitValue/exchangeVariant.Rate <= exchangeVariant.Max {
						results[key] = append(results[key], exchangeVariant)
					}
				} else if limitCurrency == "crypto" {
					if limitValue >= exchangeVariant.Min && limitValue <= exchangeVariant.Max {
						results[key] = append(results[key], exchangeVariant)
					}
				}
			}
		}
	}

	for key := range results {
		sort.Slice(results[key], func(i, j int) bool {
			return results[key][i].Rate <= results[key][j].Rate
		})
	}

	return results, nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
