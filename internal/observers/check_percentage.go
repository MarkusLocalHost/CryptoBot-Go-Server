package observers

import (
	"context"
	"cryptocurrency/internal/models/mongo_models"
	"cryptocurrency/internal/utils/telegram_api"
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func MakeObserverForPercentage(badgerDB *badger.DB, mongoDB *mongo.Client, restartMinutes time.Duration) {
	// get timestamp
	var resultsTimestamp []struct {
		Id        primitive.ObjectID `json:"_id"       bson:"_id"`
		Timestamp int64              `json:"timestamp" bson:"timestamp"`
	}

	cursor, err := mongoDB.Database(os.Getenv("MONGO_DATABASE")).Collection("price_collector_timestamps").Find(context.Background(), bson.D{
		{"timestamp", bson.D{{"$lte", time.Now().Add(-1 * restartMinutes).Unix()}}},
	})
	if err != nil {
		panic(err)
	}

	if err = cursor.All(context.TODO(), &resultsTimestamp); err != nil {
		panic(err)
	}

	// find distinct value and index
	var distinctValue int64
	var distinctIndex int
	for i, timestampValue := range resultsTimestamp {
		if i == 0 {
			distinctValue = time.Now().Add(-1*restartMinutes).Unix() - timestampValue.Timestamp
			distinctIndex = i
		}

		if time.Now().Add(-1*restartMinutes).Unix()-timestampValue.Timestamp < distinctValue {
			distinctValue = time.Now().Add(-1*restartMinutes).Unix() - timestampValue.Timestamp
			distinctIndex = i
		}
	}

	// send to admin message if distinct value too big
	if restartMinutes == time.Minute*20 {
		if distinctValue > 60*20 {
			err := telegram_api.SendMessageToNotifyAdminAboutTooBigDistinctValue(distinctValue)
			if err != nil {
				panic(err)
			}
		}
	} else if restartMinutes == time.Minute*60 {
		if distinctValue > 60*60 {
			err := telegram_api.SendMessageToNotifyAdminAboutTooBigDistinctValue(distinctValue)
			if err != nil {
				panic(err)
			}
		}
	}

	timestampForSearch := resultsTimestamp[distinctIndex].Timestamp

	// get currencies
	cursor, err = mongoDB.Database(os.Getenv("MONGO_DATABASE")).Collection("price_collector_currencies").Find(context.Background(), bson.D{})
	if err != nil {
		return
	}

	var resultsCurrencies []struct {
		Name string `json:"name"`
	}
	if err = cursor.All(context.TODO(), &resultsCurrencies); err != nil {
		panic(err)
	}

	// get last and current value
	var valPriceCurrent float64
	var valPriceCurrentCopy []byte
	var valPriceLast float64
	var valPriceLastCopy []byte
	currenciesData := make(map[string]map[string]string)
	for _, currency := range resultsCurrencies {
		opts := badger.DefaultOptions("./../../tmp/badger")
		opts.ReadOnly = true
		opts.BypassLockGuard = true
		opts.Logger = nil

		var valTimestampCopy []byte

		err := badgerDB.View(func(txn *badger.Txn) error {
			item, err := txn.Get([]byte("lastCurrencyTimestamp"))
			if err != nil {
				log.Fatal(err)
			}
			err = item.Value(func(val []byte) error {
				valTimestampCopy = append([]byte{}, val...)

				return nil
			})

			item, err = txn.Get([]byte(fmt.Sprintf(currency.Name + "-" + string(valTimestampCopy[:]))))
			if err != nil {
				if err == badger.ErrKeyNotFound {
					return nil
				} else {
					log.Fatal(err)
				}
			}
			err = item.Value(func(val []byte) error {
				valPriceCurrentCopy = append([]byte{}, val...)

				return nil
			})

			item, err = txn.Get([]byte(fmt.Sprintf(currency.Name + "-" + strconv.FormatInt(timestampForSearch, 10))))
			if err != nil {
				if err == badger.ErrKeyNotFound {
					return nil
				} else {
					log.Fatal(err)
				}
			}
			err = item.Value(func(val []byte) error {
				valPriceLastCopy = append([]byte{}, val...)

				return nil
			})

			return nil
		})
		if err != nil {
			log.Fatal(err)
		}

		valPriceCurrentString := strings.ReplaceAll(string(valPriceCurrentCopy[:]), "$", "")
		valPriceCurrentString = strings.ReplaceAll(valPriceCurrentString, ",", "")
		valPriceCurrent, _ = strconv.ParseFloat(valPriceCurrentString, 64)

		valPriceLastString := strings.ReplaceAll(string(valPriceLastCopy[:]), "$", "")
		valPriceLastString = strings.ReplaceAll(valPriceLastString, ",", "")
		valPriceLast, _ = strconv.ParseFloat(valPriceLastString, 64)

		if valPriceCurrent != 0 && valPriceLast != 0 {
			var percentChangeString string
			var percentChange float64
			percent := valPriceCurrent / valPriceLast

			if percent > 1 {
				percentChangeString = "⬆"
				percentChange = percent - 1
			} else if percent < 1 {
				percentChangeString = "⬇"
				percentChange = 1 - percent
			} else {
				percentChangeString = ""
			}

			currenciesData[currency.Name] = make(map[string]string)
			currenciesData[currency.Name]["currencyName"] = currency.Name
			currenciesData[currency.Name]["lastPrice"] = valPriceLastString
			currenciesData[currency.Name]["currentPrice"] = valPriceCurrentString
			currenciesData[currency.Name]["percentSignString"] = percentChangeString
			currenciesData[currency.Name]["percent"] = fmt.Sprintf("%f", percentChange)
		}

	}

	var filter bson.D
	if restartMinutes == time.Minute*20 {
		filter = bson.D{{"observe_20_minutes", true}}
	} else if restartMinutes == time.Minute*60 {
		filter = bson.D{{"observe_60_minutes", true}}
	}
	cursor, err = mongoDB.Database(os.Getenv("MONGO_DATABASE")).Collection("percentage_observers").Find(
		context.Background(),
		filter,
	)
	if err != nil {
		panic(err)
	}

	var resultsFromBD []mongo_models.PercentageObserver
	if err = cursor.All(context.TODO(), &resultsFromBD); err != nil {
		panic(err)
	}

	for _, result := range resultsFromBD {
		if result.FirstFilterType == "" && result.SecondFilterType == "" {
			err := telegram_api.SendMessageToNotifyAboutChangePriceInPercent(result.TelegramUserID, currenciesData)
			if err != nil {
				panic(err)
			}

		} else if result.FirstFilterType != "" && result.SecondFilterType != "" {
			currenciesDataFiltered := make(map[string]map[string]string)
			var valueFilterMustBigger float64
			var valueFilterMustSmaller float64

			if result.FirstFilterType == "percent_bigger" {
				valueFilterMustBigger = result.FirstFilterAmount
				valueFilterMustSmaller = result.SecondFilterAmount
			} else if result.FirstFilterType == "percent_smaller" {
				valueFilterMustBigger = result.SecondFilterAmount
				valueFilterMustSmaller = result.FirstFilterAmount
			} else {
				panic(err)
			}

			for _, currencyData := range currenciesData {
				percentFromData, err := strconv.ParseFloat(currencyData["percent"], 64)
				if err != nil {
					panic(err)
				}

				if percentFromData >= valueFilterMustBigger && percentFromData <= valueFilterMustSmaller {
					currenciesDataFiltered[currencyData["currencyName"]] = make(map[string]string)
					currenciesDataFiltered[currencyData["currencyName"]]["currencyName"] =
						currencyData["currencyName"]
					currenciesDataFiltered[currencyData["currencyName"]]["lastPrice"] =
						currencyData["lastPrice"]
					currenciesDataFiltered[currencyData["currencyName"]]["currentPrice"] =
						currencyData["currentPrice"]
					currenciesDataFiltered[currencyData["currencyName"]]["percentSignString"] =
						currencyData["percentSignString"]
					currenciesDataFiltered[currencyData["currencyName"]]["percent"] =
						currencyData["percent"]
				}
			}

			err := telegram_api.SendMessageToNotifyAboutChangePriceInPercent(result.TelegramUserID, currenciesDataFiltered)
			if err != nil {
				panic(err)
			}

		} else if result.FirstFilterType != "" {
			currenciesDataFiltered := make(map[string]map[string]string)
			var valueFilterMustBigger float64
			var valueFilterMustSmaller float64

			for _, currencyData := range currenciesData {
				percentFromData, err := strconv.ParseFloat(currencyData["percent"], 64)
				if err != nil {
					panic(err)
				}

				if result.FirstFilterType == "percent_bigger" {
					valueFilterMustBigger = result.FirstFilterAmount

					if percentFromData >= valueFilterMustBigger {
						currenciesDataFiltered[currencyData["currencyName"]] = make(map[string]string)
						currenciesDataFiltered[currencyData["currencyName"]]["currencyName"] =
							currencyData["currencyName"]
						currenciesDataFiltered[currencyData["currencyName"]]["lastPrice"] =
							currencyData["lastPrice"]
						currenciesDataFiltered[currencyData["currencyName"]]["currentPrice"] =
							currencyData["currentPrice"]
						currenciesDataFiltered[currencyData["currencyName"]]["percentSignString"] =
							currencyData["percentSignString"]
						currenciesDataFiltered[currencyData["currencyName"]]["percent"] =
							currencyData["percent"]
					}
				} else if result.FirstFilterType == "percent_smaller" {
					valueFilterMustBigger = result.SecondFilterAmount

					if percentFromData <= valueFilterMustSmaller {
						currenciesDataFiltered[currencyData["currencyName"]] = make(map[string]string)
						currenciesDataFiltered[currencyData["currencyName"]]["currencyName"] =
							currencyData["currencyName"]
						currenciesDataFiltered[currencyData["currencyName"]]["lastPrice"] =
							currencyData["lastPrice"]
						currenciesDataFiltered[currencyData["currencyName"]]["currentPrice"] =
							currencyData["currentPrice"]
						currenciesDataFiltered[currencyData["currencyName"]]["percentSignString"] =
							currencyData["percentSignString"]
						currenciesDataFiltered[currencyData["currencyName"]]["percent"] =
							currencyData["percent"]
					}
				} else {
					panic(err)
				}
			}

			err := telegram_api.SendMessageToNotifyAboutChangePriceInPercent(result.TelegramUserID, currenciesDataFiltered)
			if err != nil {
				panic(err)
			}

		} else {
			panic(err)
		}
	}
}
