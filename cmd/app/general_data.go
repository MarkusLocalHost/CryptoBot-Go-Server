package main

import (
	"context"
	"cryptocurrency/internal/collectors"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"os"
	"time"
)

func initGeneralData(ds *dataSource) error {
	// todo удалить timestamp которым больше 2ух суток

	client := &http.Client{}

	timestampCollection := ds.MongoDBClient.Database(os.Getenv("MONGO_DATABASE")).Collection("general_data_timestamp")

	// Loading currency from coinmarketcap
	collection := ds.MongoDBClient.Database(os.Getenv("MONGO_DATABASE")).Collection("currencies")

	estCount, err := collection.EstimatedDocumentCount(context.TODO())
	if err != nil {
		return err
	}

	var timestampResult struct {
		Title     string    `bson:"title"`
		CreatedAt time.Time `bson:"created_at"`
	}
	err = timestampCollection.FindOne(context.TODO(), bson.D{
		{"title", "general_data_currency"},
	}).Decode(&timestampResult)
	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			_, err := timestampCollection.InsertOne(context.TODO(), bson.D{
				{"title", "general_data_currency"},
				{"created_at", nil},
			})
			if err != nil {
				return err
			}
		default:
			return err
		}
	}

	if estCount == 0 || timestampResult.CreatedAt.Day() != time.Now().Day() {
		if estCount != 0 {
			err := collection.Drop(context.TODO())
			if err != nil {
				return err
			}
		}
		req, err := http.NewRequest("GET", "https://pro-api.coinmarketcap.com/v1/cryptocurrency/map", nil)
		if err != nil {
			return err
		}
		req.Header.Set("Accepts", "application/json")
		req.Header.Add("X-CMC_PRO_API_KEY", os.Getenv("CMC_API_KEY"))

		resp, err := client.Do(req)
		if err != nil {
			return err
		}

		type CurrencyData struct {
			Id                  int
			Name                string
			Symbol              string
			Slug                string
			Rank                int
			IsActive            int
			FirstHistoricalData string `json:"first_historical_data"`
			LastHistoricalData  string `json:"last_historical_data"`
			Platform            interface{}
		}
		type ResultCMC struct {
			Status interface{}
			Data   []CurrencyData
		}

		var resultCMC ResultCMC

		err = json.NewDecoder(resp.Body).Decode(&resultCMC)
		if err != nil {
			return err
		}

		var docs []interface{}
		for _, currencyData := range resultCMC.Data {
			docs = append(docs, bson.D{
				{"name", currencyData.Name},
				{"symbol", currencyData.Symbol},
				{"slug", currencyData.Slug},
				{"rank", currencyData.Rank},
				{"is_active", currencyData.IsActive},
				{"first_hist_data", currencyData.FirstHistoricalData},
				{"last_hist_data", currencyData.LastHistoricalData},
			})
		}

		_, err = collection.InsertMany(context.TODO(), docs)
		if err != nil {
			return err
		}

		// timestamp of collect general data
		_, err = timestampCollection.UpdateOne(context.TODO(), bson.D{
			{"title", "general_data_currency"},
		}, bson.D{
			{"$set", bson.D{{"created_at", time.Now()}}},
		})
		if err != nil {
			return err
		}
	}

	// Loading symbol to id for coingecko
	collection = ds.MongoDBClient.Database(os.Getenv("MONGO_DATABASE")).Collection("coingecko_currencies")

	estCount, err = collection.EstimatedDocumentCount(context.TODO())
	if err != nil {
		return err
	}

	if estCount == 0 {
		// Make request to supported VS currencies
		resp, err := http.Get("https://api.coingecko.com/api/v3/coins/list")
		if err != nil {
			return err
		}

		var result_CIS []struct {
			Id     string `bson:"id"`
			Symbol string `bson:"symbol"`
			Name   string `bson:"name"`
		}

		err = json.NewDecoder(resp.Body).Decode(&result_CIS)
		if err != nil {
			return err
		}

		// Save to Database
		var docs []interface{}
		for _, currency := range result_CIS {
			docs = append(docs, bson.D{
				{"id", currency.Id},
				{"symbol", currency.Symbol},
				{"name", currency.Name},
			})
		}

		_, err = collection.InsertMany(context.TODO(), docs)
		if err != nil {
			return err
		}
	}

	// Loading supported VS currencies
	collection = ds.MongoDBClient.Database(os.Getenv("MONGO_DATABASE")).Collection("supported_vs_currencies")

	estCount, err = collection.EstimatedDocumentCount(context.TODO())
	if err != nil {
		return err
	}

	if estCount == 0 {
		// Make request to supported VS currencies
		resp, err := http.Get("https://api.coingecko.com/api/v3/simple/supported_vs_currencies")
		if err != nil {
			return err
		}

		var result_SC []string

		err = json.NewDecoder(resp.Body).Decode(&result_SC)
		if err != nil {
			return err
		}

		// Save to Database
		var docs []interface{}
		for _, supportVsCurrency := range result_SC {
			docs = append(docs, bson.D{
				{"currency_name", supportVsCurrency},
				{"is_primary", isPrimaryCurrency(supportVsCurrency)}})
		}

		_, err = collection.InsertMany(context.TODO(), docs)
		if err != nil {
			return err
		}
	}

	// count first 1000 currencies in bdto observe
	collectionPriceCollectorCurrencies := ds.MongoDBClient.Database(os.Getenv("MONGO_DATABASE")).Collection("price_collector_currencies")

	estCount, err = collectionPriceCollectorCurrencies.EstimatedDocumentCount(context.TODO())
	if err != nil {
		return err
	}
	// timestamp of collect price collector data currencies
	err = timestampCollection.FindOne(context.TODO(), bson.D{
		{"title", "price_collector_data_currencies"},
	}).Decode(&timestampResult)
	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			_, err := timestampCollection.InsertOne(context.TODO(), bson.D{
				{"title", "price_collector_data_currencies"},
				{"created_at", nil},
			})
			if err != nil {
				return err
			}
		default:
			return err
		}
	}

	// loading first 1000 currencies to observe
	if estCount == 0 || timestampResult.CreatedAt.Hour() != time.Now().Hour() {
		cryptoNames, err := collectors.GetNameToObserveFirst1000Currency()
		if err != nil {
			return err
		}
		var docs []interface{}
		for _, cryptocurrency := range cryptoNames {
			docs = append(docs, bson.D{{"name", cryptocurrency}})
		}
		err = collectionPriceCollectorCurrencies.Drop(context.TODO())
		if err != nil {
			return err
		}
		_, err = collectionPriceCollectorCurrencies.InsertMany(context.TODO(), docs)
		if err != nil {
			return err
		}
		_, err = timestampCollection.UpdateOne(context.TODO(), bson.D{
			{"title", "price_collector_data_currencies"},
		}, bson.D{
			{"$set", bson.D{{"created_at", time.Now()}}},
		})
	}

	// clear timestamps from price_collector_timestamps
	_, err = ds.MongoDBClient.Database(os.Getenv("MONGO_DATABASE")).Collection("price_collector_timestamps").DeleteMany(context.TODO(), bson.D{
		{"timestamp", bson.D{{"$lt", time.Now().Add(-48 * time.Hour).Unix()}}},
	})
	if err != nil {
		return err
	}

	return nil
}

func isPrimaryCurrency(currency string) bool {
	switch currency {
	case
		"btc",
		"rub",
		"eth",
		"usd":
		return true
	}
	return false
}
