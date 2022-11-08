package repository

import (
	"context"
	"cryptocurrency/internal/models"
	"cryptocurrency/internal/models/mongo_models"
	"cryptocurrency/internal/observers"
	"cryptocurrency/internal/utils/apperrors"
	"github.com/dgraph-io/badger/v3"
	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"os"
	"strconv"
)

type mongoObserverRepository struct {
	DB            *mongo.Client
	BadgerDB      *badger.DB
	CronObservers *cron.Cron
}

func NewObserverRepository(db *mongo.Client, badgerDB *badger.DB, cronObservers *cron.Cron) models.ObserverRepository {
	return &mongoObserverRepository{
		DB:            db,
		BadgerDB:      badgerDB,
		CronObservers: cronObservers,
	}
}

func (m mongoObserverRepository) CreatePriceObserver(ctx context.Context, observer *mongo_models.Observer) (string, bool, error) {
	// get tier of observer
	var cryptocurrency *mongo_models.Cryptocurrency
	err := m.DB.Database(os.Getenv("MONGO_DATABASE")).Collection("currencies").FindOne(ctx, bson.D{
		{"slug", observer.CryptoID},
	}).Decode(&cryptocurrency)
	if err != nil {
		panic(err)
	}
	if cryptocurrency.Rank <= 500 {
		observer.Tier = 1
	} else {
		observer.Tier = 2
	}

	// get limits of user for observers
	var user *mongo_models.User
	err = m.DB.Database(os.Getenv("MONGO_DATABASE")).Collection("users").FindOne(ctx, bson.D{
		{"telegram_user_id", observer.TelegramUserID},
	}).Decode(&user)
	if err != nil {
		log.Printf("Could not get a user from BD with telegram_user_id: %v.Reason: %v\n", observer.TelegramUserID, err)
		return "", false, apperrors.NewInternal()
	}
	var limitOfCountObservers int64
	if user.Status == "free" {
		limitOfCountObservers, err = strconv.ParseInt(os.Getenv("FreeSubscription_CountOfObserversTier"+strconv.Itoa(observer.Tier)), 10, 64)
		if err != nil {
			panic(err)
		}
	} else {
		limitOfCountObservers, err = strconv.ParseInt(os.Getenv("PaidSubscription_CountOfObserversTier"+strconv.Itoa(observer.Tier)), 10, 64)
		if err != nil {
			panic(err)
		}
	}

	// get current count of active observers
	filterTier := bson.D{{"telegram_user_id", observer.TelegramUserID}, {"tier", observer.Tier}, {"is_active", true}}
	countOfObservers, err := m.DB.Database(os.Getenv("MONGO_DATABASE")).Collection("observers").CountDocuments(ctx, filterTier)
	if err != nil {
		log.Printf("Could not get a observers with telegram_user_id: %v.Reason: %v\n", observer.TelegramUserID, err)
		return "", false, apperrors.NewInternal()
	}
	if countOfObservers < limitOfCountObservers {
		observer.IsActive = true
	} else {
		observer.IsActive = false
	}

	// save to bd
	_, err = m.DB.Database(os.Getenv("MONGO_DATABASE")).Collection("observers").InsertOne(ctx, observer)
	if err != nil {
		log.Printf("Could not create a new observer with telegram user id: %v.Reason: %v\n", observer.TelegramUserID, err)
		return "", false, apperrors.NewInternal()
	}

	return user.Status, observer.IsActive, nil
}

func (m mongoObserverRepository) DeletePriceObserver(ctx context.Context, observerId primitive.ObjectID) error {
	filter := bson.D{{"_id", observerId}}

	_, err := m.DB.Database(os.Getenv("MONGO_DATABASE")).Collection("observers").DeleteOne(ctx, filter)
	if err != nil {
		log.Printf("Could not delete a observer with _id: %v.Reason: %v\n", observerId, err)
		return apperrors.NewInternal()
	}

	return nil
}

func (m mongoObserverRepository) ChangeStatusPriceObserver(ctx context.Context, observerId primitive.ObjectID) (status string, error error) {
	var observer *mongo_models.Observer
	err := m.DB.Database(os.Getenv("MONGO_DATABASE")).Collection("observers").FindOne(ctx, bson.D{
		{"_id", observerId},
	}).Decode(&observer)
	if err != nil {
		log.Printf("Could not get a observer from BD with _id: %v.Reason: %v\n", observerId, err)
		return "", apperrors.NewInternal()
	}

	if observer.IsActive {
		err = m.RemoveCronTask(ctx, observerId)
		if err != nil {
			log.Printf("Could not remove cron task with observer _id: %v.Reason: %v\n", observerId, err)
			return "", apperrors.NewInternal()
		}

		_, err = m.DB.Database(os.Getenv("MONGO_DATABASE")).Collection("observers").UpdateOne(ctx, bson.D{
			{"_id", observerId},
		}, bson.D{
			{"$set", bson.D{{"is_active", false}}},
		})
		if err != nil {
			log.Printf("Could not update a observer from BD with _id: %v.Reason: %v\n", observerId, err)
			return "", apperrors.NewInternal()
		}

		status = "observer stopped"
	} else {
		// get limits of user for observers
		var user *mongo_models.User
		err = m.DB.Database(os.Getenv("MONGO_DATABASE")).Collection("users").FindOne(ctx, bson.D{
			{"telegram_user_id", observer.TelegramUserID},
		}).Decode(&user)
		if err != nil {
			log.Printf("Could not get a user from BD with telegram_user_id: %v.Reason: %v\n", observer.TelegramUserID, err)
			return "", apperrors.NewInternal()
		}
		var limitOfCountObservers int64
		if user.Status == "free" {
			limitOfCountObservers, err = strconv.ParseInt(os.Getenv("FreeSubscription_CountOfObserversTier"+strconv.Itoa(observer.Tier)), 10, 64)
			if err != nil {
				panic(err)
			}
		} else {
			limitOfCountObservers, err = strconv.ParseInt(os.Getenv("PaidSubscription_CountOfObserversTier"+strconv.Itoa(observer.Tier)), 10, 64)
			if err != nil {
				panic(err)
			}
		}

		// get current count of active observers
		filterTier := bson.D{{"telegram_user_id", observer.TelegramUserID}, {"tier", observer.Tier}, {"is_active", true}}
		countOfObservers, err := m.DB.Database(os.Getenv("MONGO_DATABASE")).Collection("observers").CountDocuments(ctx, filterTier)
		if err != nil {
			log.Printf("Could not get a observers with telegram_user_id: %v.Reason: %v\n", observer.TelegramUserID, err)
			return "", apperrors.NewInternal()
		}

		// check count and starting observer
		if countOfObservers < limitOfCountObservers {
			// start cron task for observer
			err = m.MakeCronTask(ctx, observer, user.Status)
			if err != nil {
				log.Printf("Could not create cron task with observer _id: %v.Reason: %v\n", observerId, err)
				return "", apperrors.NewInternal()
			}

			_, err = m.DB.Database(os.Getenv("MONGO_DATABASE")).Collection("observers").UpdateOne(ctx, bson.D{
				{"_id", observerId},
			}, bson.D{
				{"$set", bson.D{{"is_active", true}}},
			})
			if err != nil {
				log.Printf("Could not update a observer from BD with _id: %v.Reason: %v\n", observerId, err)
				return "", apperrors.NewInternal()
			}

			status = "observer started"
		} else {
			status = "limit reached"
		}
	}

	return status, nil
}

func (m mongoObserverRepository) MakeCronTask(ctx context.Context, observer *mongo_models.Observer, userTypeSubscription string) error {
	var cronId cron.EntryID
	var timeToRestart string
	if userTypeSubscription == "free" {
		var err error
		timeToRestart = os.Getenv("FreeSubscription_TimeToObserve_Seconds")
		cronId, err = m.CronObservers.AddFunc("@every "+timeToRestart+"s", func() {
			observers.MakeObserver(observer, m.BadgerDB, m.DB, m.CronObservers, userTypeSubscription)
		})
		if err != nil {
			return err
		}
	} else if userTypeSubscription == "premium" {
		var err error
		timeToRestart = os.Getenv("PaidSubscription_TimeToObserve_Seconds")
		cronId, err = m.CronObservers.AddFunc("@every "+timeToRestart+"s", func() {
			observers.MakeObserver(observer, m.BadgerDB, m.DB, m.CronObservers, userTypeSubscription)
		})
		if err != nil {
			return err
		}
	}

	_, err := m.DB.Database(os.Getenv("MONGO_DATABASE")).Collection("cron_observer_ids").InsertOne(ctx, bson.D{
		{"cron_id", cronId},
		{"observer_id", observer.Id},
		{"time_to_restart", timeToRestart},
	})
	if err != nil {
		log.Printf("Could not delete a observer with _id: %v.Reason: %v\n", observer.Id, err)
		return apperrors.NewInternal()
	}

	return nil
}

func (m mongoObserverRepository) RemoveCronTask(ctx context.Context, observerId primitive.ObjectID) error {
	var deletedDocument bson.M
	err := m.DB.Database(os.Getenv("MONGO_DATABASE")).Collection("cron_observer_ids").FindOneAndDelete(ctx,
		bson.D{{"observer_id", observerId}}).Decode(&deletedDocument)
	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return nil
		default:
			log.Printf("Could not find and delete a observer with _id: %v.Reason: %v\n", observerId, err)
			return apperrors.NewInternal()
		}
	}

	cronId := int(deletedDocument["cron_id"].(int32))
	m.CronObservers.Remove(cron.EntryID(cronId))

	return nil
}

func (m mongoObserverRepository) CreatePercentageObserver(ctx context.Context, percentageObserver *mongo_models.PercentageObserver) (status string, err error) {
	// get limits of user for observers
	var user *mongo_models.User
	err = m.DB.Database(os.Getenv("MONGO_DATABASE")).Collection("users").FindOne(ctx, bson.D{
		{"telegram_user_id", percentageObserver.TelegramUserID},
	}).Decode(&user)
	if err != nil {
		log.Printf("Could not get a user from BD with telegram_user_id: %v.Reason: %v\n", percentageObserver.TelegramUserID, err)
		return "", apperrors.NewInternal()
	}
	var limitOfCountPercentageObservers int64
	if user.Status == "free" {
		limitOfCountPercentageObservers, err = strconv.ParseInt(os.Getenv("FreeSubscription_CountOfPercentObservers"), 10, 64)
		if err != nil {
			panic(err)
		}
	} else {
		limitOfCountPercentageObservers, err = strconv.ParseInt(os.Getenv("PaidSubscription_CountOfPercentObservers"), 10, 64)
		if err != nil {
			panic(err)
		}
	}

	// get current count of active observers
	filter := bson.D{{"telegram_user_id", percentageObserver.TelegramUserID}}
	countOfPercentageObservers, err := m.DB.Database(os.Getenv("MONGO_DATABASE")).Collection("percentage_observers").CountDocuments(ctx, filter)
	if err != nil {
		log.Printf("Could not get a observers with telegram_user_id: %v.Reason: %v\n", percentageObserver.TelegramUserID, err)
		return "", apperrors.NewInternal()
	}
	if countOfPercentageObservers < limitOfCountPercentageObservers {
		// save to bd
		_, err = m.DB.Database(os.Getenv("MONGO_DATABASE")).Collection("percentage_observers").InsertOne(ctx, percentageObserver)
		if err != nil {
			log.Printf("Could not create a new observer with telegram user id: %v.Reason: %v\n", percentageObserver.TelegramUserID, err)
			return "", apperrors.NewInternal()
		}
	} else {
		return "observer not created because limits", nil
	}

	return "observer created and started", nil
}

func (m mongoObserverRepository) DeletePercentageObserver(ctx context.Context, observerId primitive.ObjectID) error {
	filter := bson.D{{"_id", observerId}}

	_, err := m.DB.Database(os.Getenv("MONGO_DATABASE")).Collection("percentage_observers").DeleteOne(ctx, filter)
	if err != nil {
		log.Printf("Could not delete a observer with _id: %v.Reason: %v\n", observerId, err)
		return apperrors.NewInternal()
	}

	return nil
}
