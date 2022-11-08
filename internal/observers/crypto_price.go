package observers

import (
	"context"
	"cryptocurrency/internal/models/mongo_models"
	"cryptocurrency/internal/utils/coingecko_api"
	"cryptocurrency/internal/utils/observer_func"
	"cryptocurrency/internal/utils/scrapers"
	"cryptocurrency/internal/utils/telegram_api"
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"os"
	"strconv"
	"time"
)

func MakeObserver(observer *mongo_models.Observer, badgerDB *badger.DB, mongoDB *mongo.Client, cronObservers *cron.Cron, userTypeSubscription string) {
	var valPrice float64
	var valPriceCopy []byte
	var valConversation float64
	var valConversationCopy []byte

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

		item, err = txn.Get([]byte(fmt.Sprintf(observer.CryptoID + "-" + string(valTimestampCopy[:]))))
		if err != nil {
			if err == badger.ErrKeyNotFound {
				return nil
			} else {
				log.Fatal(err)
			}
		}
		err = item.Value(func(val []byte) error {
			valPriceCopy = append([]byte{}, val...)

			return nil
		})

		// get data for conversion
		if observer.CurrencyOfValue == "btc" {
			item, err = txn.Get([]byte(fmt.Sprintf("bitcoin-" + string(valTimestampCopy[:]))))
			if err != nil {
				if err == badger.ErrKeyNotFound {
					return nil
				} else {
					log.Fatal(err)
				}
			}
			err = item.Value(func(val []byte) error {
				valConversationCopy = append([]byte{}, val...)

				return nil
			})
		} else if observer.CurrencyOfValue == "eth" {
			item, err = txn.Get([]byte(fmt.Sprintf("etherium-" + string(valTimestampCopy[:]))))
			if err != nil {
				if err == badger.ErrKeyNotFound {
					return nil
				} else {
					log.Fatal(err)
				}
			}
			err = item.Value(func(val []byte) error {
				valConversationCopy = append([]byte{}, val...)

				return nil
			})
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	// check if currency found in badger
	if valPriceCopy != nil {
		// get info from badger
		valPrice, err = strconv.ParseFloat(string(valPriceCopy[:]), 64)
		if err != nil {
			log.Fatal(err)
		}

		valConversation, err = strconv.ParseFloat(string(valConversationCopy[:]), 64)
		if err != nil {
			log.Fatal(err)
		}

		if observer.CurrencyOfValue == "rub" {
			valPrice *= coingecko_api.GetPriceInCurrency("tether", "rub")
		} else if observer.CurrencyOfValue == "btc" {
			valPrice /= valConversation
		} else if observer.CurrencyOfValue == "eth" {
			valPrice /= valConversation
		}
		log.Printf("Найденное значение в badger: %f, для валюты %s", valPrice, observer.CryptoName)
	} else if userTypeSubscription == "premium" {
		// get info from cmc
		valPrice = scrapers.GetPriceInCurrencyScraper(observer.CryptoName, observer.CurrencyOfValue)
		log.Println("Не нашел знаничение в badger")
	} else if userTypeSubscription == "free" {
		// get info from coingecko
		valPrice = coingecko_api.GetPriceInCurrency(observer.CryptoID, observer.CurrencyOfValue)
		log.Println("Не нашел знаничение в badger")
	}

	// логирование запроса в бд
	_, err = mongoDB.Database(os.Getenv("MONGO_DATABASE")).Collection("log_price_observers").InsertOne(context.Background(), bson.D{
		{"observer_id", observer.Id},
		{"observed_at", time.Now()},
	})
	if err != nil {
		log.Printf("Could not log to db observer with id: %v.Reason: %v\n", observer.Id, err)
	}

	if observer.IsUpDirection {
		if valPrice >= observer.ExpectedValue {
			//log.Println("Выше! Сигнал")

			// send message
			err := telegram_api.SendMessageToNotifyAboutSignal(observer)
			if err != nil {
				log.Printf("Could not send message to user with id: %v.Reason: %v\n", observer.TelegramUserID, err)
			}

			// update observer to deactivate
			err = observer_func.DeleteObserverFromBD(mongoDB, observer)
			if err != nil {
				log.Printf("Could not update a observer from BD with _id: %v.Reason: %v\n", observer.Id, err)
			}

			// update cron tasks
			err = observer_func.UpdateCronTaskAfterSignal(mongoDB, observer, cronObservers)
			if err != nil {
				log.Printf("Could not find and delete a observer with _id: %v.Reason: %v\n", observer.Id, err)
			}
		} else {
			//log.Println("Не выше")
		}
	} else {
		if valPrice <= observer.ExpectedValue {
			//log.Println("Ниже! Сигнал")

			// send message
			err := telegram_api.SendMessageToNotifyAboutSignal(observer)
			if err != nil {
				log.Printf("Could not send message to user with id: %v.Reason: %v\n", observer.TelegramUserID, err)
			}

			// update observer to deactivate
			err = observer_func.DeleteObserverFromBD(mongoDB, observer)
			if err != nil {
				log.Printf("Could not update a observer from BD with _id: %v.Reason: %v\n", observer.Id, err)
			}

			// update cron tasks
			err = observer_func.UpdateCronTaskAfterSignal(mongoDB, observer, cronObservers)
			if err != nil {
				log.Printf("Could not find and delete a observer with _id: %v.Reason: %v\n", observer.Id, err)
			}
		} else {
			//log.Println("Не ниже")
		}
	}
}
