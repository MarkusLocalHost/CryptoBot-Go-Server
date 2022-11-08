package main

import (
	"context"
	"cryptocurrency/internal/models/mongo_models"
	"cryptocurrency/internal/observers"
	"fmt"
	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"math"
	"os"
	"time"
)

func initCronTasks(ds *dataSource) (*cron.Cron, error) {
	//Initialize Cron tasks observers
	cronObservers := cron.New()

	// Initialize cron subscription tasks
	var docsSubscription []interface{}

	collectionUsersSubscriptionTime := ds.MongoDBClient.Database(os.Getenv("MONGO_DATABASE")).Collection("users_subscription_time")
	cursor, err := collectionUsersSubscriptionTime.Find(context.Background(), bson.D{{}})
	if err != nil {
		return nil, err
	}

	var subscriptionResults []*mongo_models.UserSubscriptionTime
	if err = cursor.All(context.TODO(), &subscriptionResults); err != nil {
		return nil, err
	}
	for _, subscriptionResult := range subscriptionResults {
		if subscriptionResult.ActiveBefore.Before(time.Now()) {
			// поменять подписку пользователя на free
			_, err = ds.MongoDBClient.Database(os.Getenv("MONGO_DATABASE")).Collection("users").
				UpdateOne(
					context.Background(),
					bson.D{{"telegram_user_id", subscriptionResult.TelegramUserId}},
					bson.D{{
						"$set",
						bson.D{{"status", "free"}},
					}},
				)
			if err != nil {
				return nil, err
			}

			// удалить запись о его подписке в users_subscription_time
			err = ds.MongoDBClient.Database(os.Getenv("MONGO_DATABASE")).Collection("users_subscription_time").
				FindOneAndDelete(
					context.Background(),
					bson.D{{"telegram_user_id", subscriptionResult.TelegramUserId}},
				).Err()
			if err != nil {
				return nil, err
			}

			// todo сообщить пользователю
		} else {
			// todo тест на 1 год 10 лет
			// создать cron на изменение если меньше
			var timeToObserve string

			hoursToObserve := math.Trunc(time.Now().Sub(subscriptionResult.ActiveBefore).Hours()) * -1
			minutesToObserve := math.Trunc(time.Now().Sub(subscriptionResult.ActiveBefore).Minutes()+hoursToObserve*60) * -1
			if hoursToObserve <= 23 {
				timeToObserve = fmt.Sprintf("%d %d * * *", int(minutesToObserve), int(hoursToObserve))
				log.Println(timeToObserve)

				cronScheduleId, err := cronObservers.AddFunc(timeToObserve, func() {
					observers.ChangeSubscription(ds.MongoDBClient, ds.BadgerClient, cronObservers, subscriptionResult.TelegramUserId)
				})
				if err != nil {
					return nil, err
				}

				docsSubscription = append(docsSubscription, bson.D{
					{"cron_id", cronScheduleId},
					{"telegram_user_id", subscriptionResult.TelegramUserId},
					{"time_to_restart", timeToObserve},
				})
			}
		}
	}

	//Save cronToObserverMap to MongoDB
	collectionCronToSubscription := ds.MongoDBClient.Database(os.Getenv("MONGO_DATABASE")).Collection("cron_subscription_ids")
	if err = collectionCronToSubscription.Drop(context.Background()); err != nil {
		log.Fatalf("Failure to delete collection in MongoDB: %v\n", err)
	}
	if docsSubscription != nil {
		_, err = collectionCronToSubscription.InsertMany(context.Background(), docsSubscription)
		if err != nil {
			log.Fatalf("Failure to save to MongoDB cron ID to observer ID: %v\n", err)
		}
	}

	//Make cron to observe subscriptions
	_, err = cronObservers.AddFunc("@every 12h", func() {
		observers.ViewSubscription(ds.MongoDBClient, ds.BadgerClient, cronObservers)
	})

	//Initialize of cron ID and observer ID
	var docsObserver []interface{}

	collection := ds.MongoDBClient.Database(os.Getenv("MONGO_DATABASE")).Collection("observers")
	cursor, err = collection.Find(context.Background(), bson.D{{}})
	if err != nil {
		return nil, err
	}

	var observerResults []bson.M
	if err = cursor.All(context.TODO(), &observerResults); err != nil {
		return nil, err
	}
	for _, observerResult := range observerResults {
		observer := &mongo_models.Observer{
			Id:              observerResult["_id"].(primitive.ObjectID),
			CryptoID:        observerResult["crypto_id"].(string),
			CryptoName:      observerResult["crypto_name"].(string),
			CryptoSymbol:    observerResult["crypto_symbol"].(string),
			TelegramUserID:  observerResult["telegram_user_id"].(int64),
			CurrencyOfValue: observerResult["currency_of_value"].(string),
			ExpectedValue:   observerResult["expected_value"].(float64),
			IsUpDirection:   observerResult["is_up_direction"].(bool),
			IsActive:        observerResult["is_active"].(bool),
			Tier:            int(observerResult["tier"].(int32)),
		}

		if observer.IsActive {
			// get type of subscription for owner observer
			var user *mongo_models.User
			err = ds.MongoDBClient.Database(os.Getenv("MONGO_DATABASE")).Collection("users").FindOne(context.Background(), bson.D{
				{"telegram_user_id", observer.TelegramUserID},
			}).Decode(&user)

			var cronId cron.EntryID
			var timeToRestart string
			if user.Status == "free" {
				var err error
				timeToRestart = os.Getenv("FreeSubscription_TimeToObserve_Seconds")
				cronId, err = cronObservers.AddFunc("@every "+timeToRestart+"s", func() {
					observers.MakeObserver(observer, ds.BadgerClient, ds.MongoDBClient, cronObservers, "free")
				})
				if err != nil {
					log.Fatalf("Failure to start cron tasks: %v\n", err)
				}
			} else if user.Status == "premium" {
				var err error
				timeToRestart = os.Getenv("PaidSubscription_TimeToObserve_Seconds")
				cronId, err = cronObservers.AddFunc("@every "+timeToRestart+"s", func() {
					observers.MakeObserver(observer, ds.BadgerClient, ds.MongoDBClient, cronObservers, "premium")
				})
				if err != nil {
					log.Fatalf("Failure to start cron tasks: %v\n", err)
				}
			}

			docsObserver = append(docsObserver, bson.D{
				{"cron_id", cronId},
				{"observer_id", observer.Id},
				{"time_to_restart", timeToRestart},
			})
		}
	}

	//Save cronToObserverMap to MongoDB
	collectionCronToObserver := ds.MongoDBClient.Database(os.Getenv("MONGO_DATABASE")).Collection("cron_observer_ids")
	if err = collectionCronToObserver.Drop(context.Background()); err != nil {
		log.Fatalf("Failure to delete collection in MongoDB: %v\n", err)
	}
	if docsObserver != nil {
		_, err = collectionCronToObserver.InsertMany(context.Background(), docsObserver)
		if err != nil {
			log.Fatalf("Failure to save to MongoDB cron ID to observer ID: %v\n", err)
		}
	}

	cronObservers.Start()

	return cronObservers, nil
}
