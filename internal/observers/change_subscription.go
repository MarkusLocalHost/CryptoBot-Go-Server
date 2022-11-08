package observers

import (
	"context"
	"cryptocurrency/internal/models/mongo_models"
	"github.com/dgraph-io/badger/v3"
	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
)

func ChangeSubscription(mongoDB *mongo.Client, badgerDB *badger.DB, cronObservers *cron.Cron, telegramUserId int64) {
	// поменять подписку пользователя
	_, err := mongoDB.Database(os.Getenv("MONGO_DATABASE")).Collection("users").
		UpdateOne(
			context.Background(),
			bson.D{{"telegram_user_id", telegramUserId}},
			bson.D{{
				"$set",
				bson.D{{"status", "free"}},
			}},
		)
	if err != nil {
		panic(err)
	}

	// удалить запись о его подписке в users_subscription_time
	err = mongoDB.Database(os.Getenv("MONGO_DATABASE")).Collection("users_subscription_time").
		FindOneAndDelete(
			context.Background(),
			bson.D{{"telegram_user_id", telegramUserId}},
		).Err()
	if err != nil {
		panic(err)
	}

	// перезапустить cron на обсервер на другое время, а также поменять эти записи
	// получить id активных обсерверов
	cursor, err := mongoDB.Database(os.Getenv("MONGO_DATABASE")).
		Collection("observers").
		Find(
			context.Background(),
			bson.D{
				{"telegram_user_id", telegramUserId},
				{"is_active", true},
			})
	if err != nil {
		panic(err)
	}

	var results []*mongo_models.Observer
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	// пройтись по всем обсерверам и cron tasks
	for _, observer := range results {
		// получить запись
		var cronObserverId *mongo_models.CronObserverIds

		err = mongoDB.Database(os.Getenv("MONGO_DATABASE")).
			Collection("cron_observer_ids").
			FindOne(context.Background(), bson.D{{"observer_id", observer.Id}}).
			Decode(&cronObserverId)

		// удалить cron
		cronObservers.Remove(cron.EntryID(cronObserverId.CronId))

		// запустить cron
		cronId, err := cronObservers.AddFunc("@every "+os.Getenv("FreeSubscription_TimeToObserve_Seconds")+"s", func() {
			MakeObserver(observer, badgerDB, mongoDB, cronObservers, "free")
		})
		if err != nil {
			panic(err)
		}

		// изменить запись
		err = mongoDB.Database(os.Getenv("MONGO_DATABASE")).
			Collection("cron_observer_ids").
			FindOneAndUpdate(
				context.Background(),
				bson.D{{"observer_id", observer.Id}},
				bson.D{{
					"$set",
					bson.D{
						{"cron_id", cronId},
						{"time_to_restart", os.Getenv("FreeSubscription_TimeToObserve_Seconds")},
					},
				}},
			).Err()
		if err != nil {
			panic(err)
		}
	}
}
