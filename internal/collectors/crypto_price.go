package collectors

import (
	"context"
	"fmt"
	badger "github.com/dgraph-io/badger/v3"
	"github.com/gocolly/colly/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func GetPrice(dbBadger *badger.DB, dbMongo *mongo.Client) {
	startTime := time.Now()

	link := "https://www.coingecko.com"
	c := colly.NewCollector(
		colly.AllowedDomains("www.coingecko.com"),
	)
	err := c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 10})
	if err != nil {
		panic(err)
	}

	timestampOfRequest := time.Now().Unix()

	// old code for CMC
	//cryptosPrice := make(map[string]string)
	//priceOfTenFirstCoins := make([]string, 0)
	//cryptoNamesForAdditionalInfo := make([]string, 0)
	//
	//c.OnHTML("tbody tr", func(e *colly.HTMLElement) {
	//	isTenFirstCoins := false
	//	symbolOfCurrentCoin := ""
	//
	//	e.ForEach("span", func(i int, element *colly.HTMLElement) {
	//		if strings.HasPrefix(element.Text, "$") && i == 3 && len(priceOfTenFirstCoins) <= 10 {
	//			isTenFirstCoins = true
	//			priceOfTenFirstCoins = append(priceOfTenFirstCoins, strings.TrimSpace(element.Text))
	//		} else if i == 0 {
	//			isTenFirstCoins = false
	//		} else if i == 3 && !isTenFirstCoins {
	//			symbolOfCurrentCoin = strings.ReplaceAll(strings.ToLower(element.Text), " ", "-")
	//		} else if i == 5 && !isTenFirstCoins {
	//			priceFloat, err := strconv.ParseFloat(strings.ReplaceAll(strings.ReplaceAll(element.Text, "$", ""), ",", ""), 64)
	//			if err != nil {
	//				log.Fatalf(err.Error())
	//			}
	//			if priceFloat < 0.1 {
	//				cryptoNamesForAdditionalInfo = append(cryptoNamesForAdditionalInfo, symbolOfCurrentCoin)
	//			}
	//			cryptosPrice[symbolOfCurrentCoin] = element.Text
	//		}
	//	})
	//
	//	isTenFirstCoins = true
	//
	//	e.ForEach("p", func(i int, element *colly.HTMLElement) {
	//		if i == 1 && isTenFirstCoins {
	//			cryptosPrice[strings.ReplaceAll(strings.ToLower(element.Text), " ", "-")] = priceOfTenFirstCoins[0]
	//			priceOfTenFirstCoins = priceOfTenFirstCoins[1:]
	//		}
	//	})
	//})

	cryptosPrice := make(map[string]string)
	c.OnHTML("tbody tr", func(e *colly.HTMLElement) {
		cryptoName := ""

		e.ForEach("td", func(i int, element *colly.HTMLElement) {
			if i == 2 {
				cryptoName = strings.ReplaceAll(strings.ToLower(element.Attr("data-sort")), " ", "-")
			} else if i == 3 {
				cryptosPrice[cryptoName] = element.Attr("data-sort")
			}
		})
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Failed with response:", r, "\nError:", err)
	})

	err = c.Visit(link)
	if err != nil {
		panic(err)
	}
	err = c.Visit(link + "?page=2")
	if err != nil {
		panic(err)
	}
	err = c.Visit(link + "?page=3")
	if err != nil {
		panic(err)
	}
	err = c.Visit(link + "?page=4")
	if err != nil {
		panic(err)
	}
	err = c.Visit(link + "?page=5")
	if err != nil {
		panic(err)
	}
	err = c.Visit(link + "?page=6")
	if err != nil {
		panic(err)
	}
	err = c.Visit(link + "?page=7")
	if err != nil {
		panic(err)
	}
	err = c.Visit(link + "?page=8")
	if err != nil {
		panic(err)
	}
	err = c.Visit(link + "?page=9")
	if err != nil {
		panic(err)
	}
	err = c.Visit(link + "?page=10")
	if err != nil {
		panic(err)
	}
	c.Wait()

	err = writeLogToMongoDB(time.Since(startTime), len(cryptosPrice), dbMongo)
	if err != nil {
		log.Println("Failed to write in mongoDB on request:", "\nError:", err)
	}

	err = writeToBadgerDB(cryptosPrice, timestampOfRequest, dbBadger)
	if err != nil {
		log.Println("Failed to write in badgerDB on request:", "\nError:", err)
	}

	err = writeTimestampToMongoDB(timestampOfRequest, dbMongo)
	if err != nil {
		log.Println("Failed to write in mongoDB on request:", "\nError:", err)
	}
}

func GetNameToObserveFirst1000Currency() ([]string, error) {
	link := "https://coinmarketcap.com/"
	page := 1
	c := colly.NewCollector(
		colly.AllowedDomains("coinmarketcap.com"),
	)

	var cryptoNames []string

	c.OnHTML("tbody tr", func(e *colly.HTMLElement) {
		e.ForEach("span", func(i int, element *colly.HTMLElement) {
			if i == 3 && !strings.Contains(element.Text, "%") {
				cryptoNames = append(cryptoNames, strings.ToLower(strings.ReplaceAll(element.Text, " ", "-")))
			}
		})

		e.ForEach("p", func(i int, element *colly.HTMLElement) {
			if i == 1 {
				cryptoNames = append(cryptoNames, strings.ToLower(strings.ReplaceAll(element.Text, " ", "-")))
			}
		})
	})

	c.OnScraped(func(r *colly.Response) {
		if page < 10 {
			page += 1
			err := r.Request.Visit(fmt.Sprintf("%s?page=%v", link, page))
			if err != nil {
				log.Println("Failed with response on page:", page, "\nError:", err)
			}
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Failed with response:", r, "\nError:", err)
	})

	err := c.Visit(link)
	if err != nil {
		return nil, err
	}

	return cryptoNames, nil
}

//func getPriceCurrentCurrency(subLink string) (string, error) {
//	link := "https://coinmarketcap.com" + subLink
//	c := colly.NewCollector(
//		colly.AllowedDomains("coinmarketcap.com"),
//	)
//
//	var price string
//
//	c.OnHTML(".priceValue", func(e *colly.HTMLElement) {
//		price = e.ChildText("span")
//	})
//
//	c.OnError(func(r *colly.Response, err error) {
//		log.Println("Failed with response:", r, "\nError:", err)
//	})
//
//	err := c.Visit(link)
//	if err != nil {
//		return "", err
//	}
//
//	return price, nil
//}

func writeToBadgerDB(data map[string]string, timestamp int64, db *badger.DB) error {
	for cryptocurrency, price := range data {
		err := db.Update(func(txn *badger.Txn) error {
			e := badger.NewEntry([]byte(strings.ToLower(cryptocurrency)+"-"+strconv.FormatInt(timestamp, 10)), []byte(price)).WithTTL(time.Hour * 48)
			err := txn.SetEntry(e)
			return err
		})
		if err != nil {
			return err
		}
	}

	err := db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte("lastCurrencyTimestamp"), []byte(strconv.FormatInt(timestamp, 10)))
		return err
	})
	if err != nil {
		return err
	}

	return nil
}

func writeTimestampToMongoDB(timestamp int64, db *mongo.Client) error {
	_, err := db.Database(os.Getenv("MONGO_DATABASE")).Collection("price_collector_timestamps").InsertOne(context.Background(), bson.D{
		{"timestamp", timestamp},
	})
	if err != nil {
		return err
	}

	return nil
}

func writeLogToMongoDB(requestTimeLength time.Duration, currenciesLength int, db *mongo.Client) error {
	_, err := db.Database(os.Getenv("MONGO_DATABASE")).Collection("log_collector_prices").InsertOne(context.Background(), bson.D{
		{"request_time_length", requestTimeLength.Seconds()},
		{"currencies_length", currenciesLength},
		{"created_at", time.Now()},
	})
	if err != nil {
		return err
	}

	return nil
}
