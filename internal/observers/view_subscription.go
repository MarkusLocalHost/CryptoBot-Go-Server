package observers

import (
	"context"
	"cryptocurrency/internal/models/mongo_models"
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"math"
	"os"
	"time"
)

func ViewSubscription(mongoDB *mongo.Client, badgerDB *badger.DB, cronObservers *cron.Cron) {
	var docsSubscription []interface{}

	collectionUsersSubscriptionTime := mongoDB.Database(os.Getenv("MONGO_DATABASE")).Collection("users_subscription_time")
	cursor, err := collectionUsersSubscriptionTime.Find(context.Background(), bson.D{{}})
	if err != nil {
		panic(err)
	}

	collectionCronToSubscription := mongoDB.Database(os.Getenv("MONGO_DATABASE")).Collection("cron_subscription_ids")
	var subscriptionResults []*mongo_models.UserSubscriptionTime
	if err = cursor.All(context.TODO(), &subscriptionResults); err != nil {
		panic(err)
	}
	for _, subscriptionResult := range subscriptionResults {
		// создать cron на изменение если меньше
		var timeToObserve string

		hoursToObserve := math.Trunc(time.Now().Sub(subscriptionResult.ActiveBefore).Hours()) * -1
		minutesToObserve := math.Trunc(time.Now().Sub(subscriptionResult.ActiveBefore).Minutes()+hoursToObserve*60) * -1
		if hoursToObserve <= 23 {
			// проверим если этот пользователь уже в списке cron на смену подписки
			count, err := collectionCronToSubscription.CountDocuments(context.Background(), bson.D{{"telegram_user_id", subscriptionResult.TelegramUserId}})
			if err != nil {
				panic(err)
			}

			if count == 0 {
				timeToObserve = fmt.Sprintf("%d %d * * *", int(minutesToObserve), int(hoursToObserve))

				cronScheduleId, err := cronObservers.AddFunc(timeToObserve, func() {
					ChangeSubscription(mongoDB, badgerDB, cronObservers, subscriptionResult.TelegramUserId)
				})
				if err != nil {
					panic(err)
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
	if err = collectionCronToSubscription.Drop(context.Background()); err != nil {
		log.Fatalf("Failure to delete collection in MongoDB: %v\n", err)
	}
	if docsSubscription != nil {
		_, err = collectionCronToSubscription.InsertMany(context.Background(), docsSubscription)
		if err != nil {
			log.Fatalf("Failure to save to MongoDB cron ID to observer ID: %v\n", err)
		}
	}
}
