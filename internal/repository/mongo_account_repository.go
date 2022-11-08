package repository

import (
	"context"
	"cryptocurrency/internal/models"
	"cryptocurrency/internal/models/mongo_models"
	"cryptocurrency/internal/observers"
	"cryptocurrency/internal/utils/apperrors"
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type mongoAccountRepository struct {
	DB            *mongo.Client
	BadgerDB      *badger.DB
	CronObservers *cron.Cron
}

func NewAccountRepository(db *mongo.Client, badgerDB *badger.DB, cronObservers *cron.Cron) models.AccountRepository {
	return &mongoAccountRepository{
		DB:            db,
		BadgerDB:      badgerDB,
		CronObservers: cronObservers,
	}
}

func (m mongoAccountRepository) CreateAccount(ctx context.Context, telegramUserId int64, language string) error {
	filter := bson.D{{"telegram_user_id", telegramUserId}}
	count, err := m.DB.Database(os.Getenv("MONGO_DATABASE")).Collection("users").CountDocuments(ctx, filter)
	if err != nil {
		panic(err)
	}

	if count == 0 {
		_, err = m.DB.Database(os.Getenv("MONGO_DATABASE")).Collection("users").InsertOne(ctx, bson.D{
			{"telegram_user_id", telegramUserId},
			{"status", "free"},
			{"language", language},
			{"created_at", time.Now()},
		})
		if err != nil {
			panic(err)
		}
	}

	return nil
}

func (m mongoAccountRepository) GetUserPriceObservers(ctx context.Context, telegramUserId int64) (observer []*mongo_models.Observer, err error) {
	coll := m.DB.Database(os.Getenv("MONGO_DATABASE")).Collection("observers")
	filter := bson.D{{"telegram_user_id", telegramUserId}}

	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		panic(err)
	}

	var results []*mongo_models.Observer
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	return results, nil
}

func (m mongoAccountRepository) GetUserPercentageObservers(ctx context.Context, telegramUserId int64) (observers []*mongo_models.PercentageObserver, err error) {
	coll := m.DB.Database(os.Getenv("MONGO_DATABASE")).Collection("percentage_observers")
	filter := bson.D{{"telegram_user_id", telegramUserId}}

	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		panic(err)
	}

	var results []*mongo_models.PercentageObserver
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	return results, nil
}

func (m mongoAccountRepository) ViewUserSubscription(ctx context.Context, telegramUserId int64) (string, error) {
	var user *mongo_models.User
	err := m.DB.Database(os.Getenv("MONGO_DATABASE")).Collection("users").FindOne(ctx, bson.D{
		{"telegram_user_id", telegramUserId},
	}).Decode(&user)
	if err != nil {
		log.Printf("Could not get a user from BD with telegram_user_id: %v.Reason: %v\n", telegramUserId, err)
		return "", apperrors.NewInternal()
	}

	return user.Status, nil
}

func (m mongoAccountRepository) ViewCountOfUserActiveObserver(ctx context.Context, telegramUserId int64) (int, int, int, error) {
	filterTier1 := bson.D{{"telegram_user_id", telegramUserId}, {"tier", 1}, {"is_active", true}}
	countOfObserversTier1, err := m.DB.Database(os.Getenv("MONGO_DATABASE")).Collection("observers").CountDocuments(ctx, filterTier1)
	if err != nil {
		log.Printf("Could not get a observers with telegram_user_id: %v.Reason: %v\n", telegramUserId, err)
		return 0, 0, 0, apperrors.NewInternal()
	}

	filterTier2 := bson.D{{"telegram_user_id", telegramUserId}, {"tier", 2}, {"is_active", true}}
	countOfObserversTier2, err := m.DB.Database(os.Getenv("MONGO_DATABASE")).Collection("observers").CountDocuments(ctx, filterTier2)
	if err != nil {
		log.Printf("Could not get a observers with telegram_user_id: %v.Reason: %v\n", telegramUserId, err)
		return 0, 0, 0, apperrors.NewInternal()
	}

	filterPercentage := bson.D{{"telegram_user_id", telegramUserId}}
	countOfObserversPercentage, err := m.DB.Database(os.Getenv("MONGO_DATABASE")).Collection("percentage_observers").CountDocuments(ctx, filterPercentage)
	if err != nil {
		log.Printf("Could not get a observers with telegram_user_id: %v.Reason: %v\n", telegramUserId, err)
		return 0, 0, 0, apperrors.NewInternal()
	}

	return int(countOfObserversTier1), int(countOfObserversTier2), int(countOfObserversPercentage), nil
}

func (m mongoAccountRepository) ViewUserPortfolio(ctx context.Context, telegramUserId int64) ([]*mongo_models.Portfolio, error) {
	coll := m.DB.Database(os.Getenv("MONGO_DATABASE")).Collection("portfolio")
	filter := bson.D{{"telegram_user_id", telegramUserId}}

	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		panic(err)
	}

	var results []*mongo_models.Portfolio
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	return results, nil
}

func (m mongoAccountRepository) AddToUserPortfolio(ctx context.Context, data *mongo_models.Portfolio) error {
	// save to bd
	_, err := m.DB.Database(os.Getenv("MONGO_DATABASE")).Collection("portfolio").InsertOne(ctx, data)
	if err != nil {
		log.Printf("Could not create a new observer with telegram user id: %v.Reason: %v\n", data.TelegramUserID, err)
		return apperrors.NewInternal()
	}

	return nil
}

func (m mongoAccountRepository) UpdateElementPriceUserPortfolio(ctx context.Context, id primitive.ObjectID, price float64) error {
	_, err := m.DB.Database(os.Getenv("MONGO_DATABASE")).Collection("portfolio").UpdateOne(ctx, bson.D{
		{"_id", id},
	}, bson.D{
		{"$set", bson.D{{"price", price}}},
	})
	if err != nil {
		log.Printf("Could not update a observer from BD with _id: %v.Reason: %v\n", id, err)
		return apperrors.NewInternal()
	}

	return nil
}

func (m mongoAccountRepository) UpdateElementValueUserPortfolio(ctx context.Context, id primitive.ObjectID, value float64) error {
	_, err := m.DB.Database(os.Getenv("MONGO_DATABASE")).Collection("portfolio").UpdateOne(ctx, bson.D{
		{"_id", id},
	}, bson.D{
		{"$set", bson.D{{"value", value}}},
	})
	if err != nil {
		log.Printf("Could not update a observer from BD with _id: %v.Reason: %v\n", id, err)
		return apperrors.NewInternal()
	}

	return nil
}

func (m mongoAccountRepository) DeleteElementUserPortfolio(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.D{{"_id", id}}

	_, err := m.DB.Database(os.Getenv("MONGO_DATABASE")).Collection("portfolio").DeleteOne(ctx, filter)
	if err != nil {
		log.Printf("Could not delete a observer with _id: %v.Reason: %v\n", id, err)
		return apperrors.NewInternal()
	}

	return nil
}

func (m mongoAccountRepository) CheckPromoCode(ctx context.Context, promoCodeValue string, telegramUserId int64) (string, error) {
	// check promo code in db
	var promoCode *mongo_models.PromoCode
	err := m.DB.Database(os.Getenv("MONGO_DATABASE")).Collection("promo_codes").FindOne(ctx, bson.D{
		{"value", promoCodeValue},
	}).Decode(&promoCode)
	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return "no promo code in db", nil
		default:
			return "", err
		}
	}

	// check status by current time
	if promoCode.ActiveBefore.Before(time.Now()) {
		return "time of active is ended", nil
	}

	// check count of activation
	count, err := m.DB.Database(os.Getenv("MONGO_DATABASE")).Collection("promo_codes_activations").CountDocuments(ctx,
		bson.D{
			{"promo_code_id", promoCode.Id},
		})
	if err != nil {
		return "", err
	}

	if count < int64(promoCode.CountOfActivation) {
		// check is activate this user this promo code
		checkCount, err := m.DB.Database(os.Getenv("MONGO_DATABASE")).Collection("promo_codes_activations").CountDocuments(ctx,
			bson.D{
				{"promo_code_id", promoCode.Id},
				{"telegram_user_id", telegramUserId},
			})
		if err != nil {
			return "", err
		}

		if checkCount == 1 {
			return "you already activate this promo code", nil
		} else {
			// input log of activation
			_, err := m.DB.Database(os.Getenv("MONGO_DATABASE")).Collection("promo_codes_activations").InsertOne(context.TODO(), bson.D{
				{"promo_code_id", promoCode.Id},
				{"telegram_user_id", telegramUserId},
				{"activated_at", time.Now()},
			})
			if err != nil {
				return "", err
			}

			// activate premium
			err = m.ChangeUserSubscription(ctx, telegramUserId, "premium")
			if err != nil {
				return "", err
			}

			// add hours to subscription
			timeDuration, err := time.ParseDuration(promoCode.SubscriptionHours)
			if err != nil {
				return "", err
			}
			err = m.AddTimeToUserSubscription(ctx, telegramUserId, timeDuration)
			if err != nil {
				return "", err
			}

			return "promo code activated", nil
		}
	} else {
		return "limit of activation", nil
	}
}

func (m mongoAccountRepository) ChangeUserSubscription(ctx context.Context, telegramUserId int64, status string) error {
	// получить пользователя
	var user *mongo_models.User

	err := m.DB.Database(os.Getenv("MONGO_DATABASE")).Collection("users").
		FindOne(ctx, bson.D{{"telegram_user_id", telegramUserId}}).
		Decode(&user)
	if err != nil {
		return err
	}

	// просмотр статуса
	if user.Status == "free" {
		// поменять статус на премиум
		_, err = m.DB.Database(os.Getenv("MONGO_DATABASE")).Collection("users").
			UpdateOne(
				ctx,
				bson.D{{"telegram_user_id", telegramUserId}},
				bson.D{{
					"$set",
					bson.D{{"status", "premium"}},
				}},
			)

		// перезапустить cron на обсервер на другое время, а также поменять эти записи
		// получить id активных обсерверов
		cursor, err := m.DB.Database(os.Getenv("MONGO_DATABASE")).
			Collection("observers").
			Find(
				ctx,
				bson.D{
					{"telegram_user_id", telegramUserId},
					{"is_active", true},
				})
		if err != nil {
			return err
		}

		var results []*mongo_models.Observer
		if err = cursor.All(context.TODO(), &results); err != nil {
			panic(err)
		}

		// пройтись по всем обсерверам и cron tasks
		for _, observer := range results {
			// получить запись
			var cronObserverId *mongo_models.CronObserverIds

			err = m.DB.Database(os.Getenv("MONGO_DATABASE")).
				Collection("cron_observer_ids").
				FindOne(ctx, bson.D{{"observer_id", observer.Id}}).
				Decode(&cronObserverId)

			// удалить cron
			m.CronObservers.Remove(cron.EntryID(cronObserverId.CronId))

			// запустить cron
			cronId, err := m.CronObservers.AddFunc("@every "+os.Getenv("PaidSubscription_TimeToObserve_Seconds")+"s", func() {
				observers.MakeObserver(observer, m.BadgerDB, m.DB, m.CronObservers, "premium")
			})
			if err != nil {
				return err
			}

			// изменить запись
			err = m.DB.Database(os.Getenv("MONGO_DATABASE")).
				Collection("cron_observer_ids").
				FindOneAndUpdate(
					ctx,
					bson.D{{"observer_id", observer.Id}},
					bson.D{{
						"$set",
						bson.D{
							{"cron_id", cronId},
							{"time_to_restart", os.Getenv("PaidSubscription_TimeToObserve_Seconds")},
						},
					}},
				).Err()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (m mongoAccountRepository) AddTimeToUserSubscription(ctx context.Context, telegramUserId int64, hoursToAdd time.Duration) error {
	var insertedId interface{}
	// проверить есть ли записи, если нет, то создать
	var userSubscriptionTime *mongo_models.UserSubscriptionTime

	err := m.DB.Database(os.Getenv("MONGO_DATABASE")).Collection("users_subscription_time").
		FindOne(ctx, bson.D{{"telegram_user_id", telegramUserId}}).
		Decode(&userSubscriptionTime)
	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			insertedId, err = m.DB.Database(os.Getenv("MONGO_DATABASE")).Collection("users_subscription_time").InsertOne(context.TODO(), bson.D{
				{"telegram_user_id", telegramUserId},
				{"active_before", time.Now().Add(hoursToAdd)},
			})
			if err != nil {
				return err
			}
		default:
			return err
		}
	}

	if userSubscriptionTime != nil {
		_, err = m.DB.Database(os.Getenv("MONGO_DATABASE")).Collection("users_subscription_time").
			UpdateOne(
				ctx,
				bson.D{{"_id", userSubscriptionTime.Id}},
				bson.D{{
					"$set",
					bson.D{{"active_before", userSubscriptionTime.ActiveBefore.Add(hoursToAdd)}},
				}},
			)
		if err != nil {
			return err
		}
	} else {
		_, err = m.DB.Database(os.Getenv("MONGO_DATABASE")).Collection("users_subscription_time").
			UpdateOne(
				ctx,
				bson.D{{"_id", insertedId}},
				bson.D{{
					"$set",
					bson.D{{"active_before", time.Now().Add(hoursToAdd)}},
				}},
			)
		if err != nil {
			return err
		}
	}

	// проверить есть ли cron на изменение подписки и если да, то пересоздать, добавив время
	collectionCronToSubscription := m.DB.Database(os.Getenv("MONGO_DATABASE")).Collection("cron_subscription_ids")
	count, err := collectionCronToSubscription.CountDocuments(context.Background(), bson.D{{"telegram_user_id", telegramUserId}})
	if err != nil {
		return err
	}

	if count != 0 {
		var cronSubscriptionId mongo_models.CronSubscriptionIds
		err = collectionCronToSubscription.FindOne(ctx, bson.D{
			{"telegram_user_id", telegramUserId},
		}).Decode(&cronSubscriptionId)
		if err != nil {
			return err
		}

		hoursInTimeToRestart := strings.Split(cronSubscriptionId.TimeToRestart, " ")
		hours, err := strconv.Atoi(hoursInTimeToRestart[1])
		if err != nil {
			return err
		}
		// если часов по итогу меньше чем 23, то пересоздать cron, если больше то просто удалить
		if time.Duration(hours*1000000*60*60).Hours()+hoursToAdd.Hours() <= 23 {
			// пересоздать cron
			// удалить cron
			m.CronObservers.Remove(cron.EntryID(cronSubscriptionId.CronId))

			// запустить cron
			hoursToObserve := time.Duration(hours*1000000*60*60).Hours() + hoursToAdd.Hours()
			minutesToObserve := hoursInTimeToRestart[0]
			timeToObserve := fmt.Sprintf("%s %d * * *", minutesToObserve, int(hoursToObserve))
			log.Println(timeToObserve)

			cronScheduleId, err := m.CronObservers.AddFunc(timeToObserve, func() {
				observers.ChangeSubscription(m.DB, m.BadgerDB, m.CronObservers, telegramUserId)
			})
			if err != nil {
				panic(err)
			}

			// изменить запись
			_, err = collectionCronToSubscription.UpdateOne(ctx, bson.D{
				{"_id", cronSubscriptionId.CronId},
			}, bson.D{
				{"$set", bson.D{
					{"cron_id", cronScheduleId},
				}},
			})
			if err != nil {
				return err
			}
		} else {
			// удалить запись
			_, err = collectionCronToSubscription.DeleteOne(ctx, bson.D{
				{"_id", cronSubscriptionId.CronId},
			})
			if err != nil {
				return err
			}

			// удалить cron
			m.CronObservers.Remove(cron.EntryID(cronSubscriptionId.CronId))
		}
	}

	return nil
}

func (m mongoAccountRepository) GetUsersByFilter(ctx context.Context, filter string) ([]mongo_models.User, error) {
	coll := m.DB.Database(os.Getenv("MONGO_DATABASE")).Collection("users")

	var filterDB bson.D
	if filter == "free" {
		filterDB = bson.D{{"status", "free"}}
	} else if filter == "premium" {
		filterDB = bson.D{{"status", "premium"}}
	} else {
		filterDB = bson.D{{}}
	}

	cursor, err := coll.Find(context.TODO(), filterDB)
	if err != nil {
		panic(err)
	}

	var results []mongo_models.User
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	return results, nil
}
