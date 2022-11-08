package repository

import (
	"context"
	"cryptocurrency/internal/models"
	"cryptocurrency/internal/models/mongo_models"
	"cryptocurrency/internal/models/response_models/slack_bot"
	"cryptocurrency/internal/utils/apperrors"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"os"
	"time"
)

type managerRepository struct {
	DB *mongo.Client
}

func NewManagerRepository(db *mongo.Client) models.ManagerRepository {
	return &managerRepository{
		DB: db,
	}
}

func (m managerRepository) CreatePromoCode(ctx context.Context, promoCode *mongo_models.PromoCode) error {
	filter := bson.D{{"value", promoCode.Value}}
	count, err := m.DB.Database(os.Getenv("MONGO_DATABASE")).Collection("promo_codes").CountDocuments(ctx, filter)
	if err != nil {
		panic(err)
	}

	if count == 0 {
		_, err := m.DB.Database(os.Getenv("MONGO_DATABASE")).Collection("promo_codes").InsertOne(ctx, promoCode)
		if err != nil {
			log.Printf("Could not create a new promo code with .Reason: %v\n", err)
			return err
		}
	} else {
		return errors.New("promocode with this value already exists")
	}

	return nil
}

func (m managerRepository) ViewPromoCodes(ctx context.Context) ([]*slack_bot.PromoCodesView, error) {
	promoCodesCollection := m.DB.Database(os.Getenv("MONGO_DATABASE")).Collection("promo_codes")

	cursorPromoCodes, err := promoCodesCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		panic(err)
	}

	var promoCodes []mongo_models.PromoCode
	if err = cursorPromoCodes.All(context.TODO(), &promoCodes); err != nil {
		panic(err)
	}

	var data []*slack_bot.PromoCodesView

	for _, promoCode := range promoCodes {
		data = append(data, &slack_bot.PromoCodesView{
			PromoCode:            promoCode,
			PromoCodeActivations: nil,
		})
	}

	promoCodesActivationsCollection := m.DB.Database(os.Getenv("MONGO_DATABASE")).Collection("promo_codes_activations")

	cursorPromoCodesActivations, err := promoCodesActivationsCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		panic(err)
	}

	var promoCodesActivations []mongo_models.PromoCodeActivations
	if err = cursorPromoCodesActivations.All(context.TODO(), &promoCodesActivations); err != nil {
		panic(err)
	}

	for _, promoCodesActivation := range promoCodesActivations {
		for _, promoCodeView := range data {
			if promoCodesActivation.PromoCodeId == promoCodeView.PromoCode.Id {
				promoCodeView.PromoCodeActivations = append(promoCodeView.PromoCodeActivations, promoCodesActivation)
			}
		}
	}

	return data, nil
}

func (m managerRepository) ViewOnlineUsersStats(ctx context.Context) ([]*slack_bot.OnlineUserStats, error) {
	logRequestsCollection := m.DB.Database(os.Getenv("MONGO_DATABASE")).Collection("log_requests")
	cursorOnlineUsersLogs, err := logRequestsCollection.Find(ctx, bson.D{
		{"telegram_user_id", bson.D{{"$exists", true}}},
	})
	if err != nil {
		return nil, err
	}

	var logRequests []*mongo_models.LogRequest
	if err = cursorOnlineUsersLogs.All(context.TODO(), &logRequests); err != nil {
		panic(err)
	}

	var onlineUsersStats []*slack_bot.OnlineUserStats
	for _, onlineUserLog := range logRequests {
		var isFound = false
		for _, onlineUserStat := range onlineUsersStats {
			if onlineUserStat.Time.Year() == onlineUserLog.CreatedAt.Year() &&
				onlineUserStat.Time.Month() == onlineUserLog.CreatedAt.Month() &&
				onlineUserStat.Time.Day() == onlineUserLog.CreatedAt.Day() &&
				onlineUserStat.Time.Hour() == onlineUserLog.CreatedAt.Hour() &&
				!contains(onlineUserStat.IdsUsers, onlineUserLog.TelegramUserId) {
				// добавить количество пользователей +1
				onlineUserStat.OnlineUsers += 1
				onlineUserStat.IdsUsers = append(onlineUserStat.IdsUsers, onlineUserLog.TelegramUserId)

				isFound = true
			}
		}

		if isFound == false {
			// добавить новую стату
			var idsUsers []int64
			idsUsers = append(idsUsers, onlineUserLog.TelegramUserId)

			onlineUsersStats = append(onlineUsersStats, &slack_bot.OnlineUserStats{
				OnlineUsers: 1,
				IdsUsers:    idsUsers,
				Time:        onlineUserLog.CreatedAt,
			})
		}
	}

	return onlineUsersStats, nil
}

func (m managerRepository) ViewCountRequests(ctx context.Context) ([]*slack_bot.LogRequestStats, error) {
	logRequestsCollection := m.DB.Database(os.Getenv("MONGO_DATABASE")).Collection("log_requests")
	cursorOnlineUsersLogs, err := logRequestsCollection.Find(ctx, bson.D{
		{"telegram_user_id", bson.D{{"$exists", true}}},
	})
	if err != nil {
		return nil, err
	}

	var logRequests []*mongo_models.LogRequest
	if err = cursorOnlineUsersLogs.All(context.TODO(), &logRequests); err != nil {
		panic(err)
	}

	var logRequestStats []*slack_bot.LogRequestStats
	for _, onlineUserLog := range logRequests {
		var isFound = false
		for _, logRequestStat := range logRequestStats {
			if logRequestStat.Time.Year() == onlineUserLog.CreatedAt.Year() &&
				logRequestStat.Time.Month() == onlineUserLog.CreatedAt.Month() &&
				logRequestStat.Time.Day() == onlineUserLog.CreatedAt.Day() &&
				logRequestStat.Time.Hour() == onlineUserLog.CreatedAt.Hour() {
				// добавить количество запросов +1
				logRequestStat.RequestCount += 1

				isFound = true
			}
		}

		if isFound == false {
			// добавить новую стату
			logRequestStats = append(logRequestStats, &slack_bot.LogRequestStats{
				RequestCount: 1,
				Time:         onlineUserLog.CreatedAt,
			})
		}
	}

	return logRequestStats, nil
}

func (m managerRepository) ViewUserActions(ctx context.Context, telegramUserId int64) ([]*mongo_models.LogRequest, error) {
	logRequestsCollection := m.DB.Database(os.Getenv("MONGO_DATABASE")).Collection("log_requests")
	cursorOnlineUsersLogs, err := logRequestsCollection.Find(ctx, bson.D{
		{"telegram_user_id", telegramUserId},
	})
	if err != nil {
		return nil, err
	}

	var logRequests []*mongo_models.LogRequest
	if err = cursorOnlineUsersLogs.All(context.TODO(), &logRequests); err != nil {
		panic(err)
	}

	return logRequests, nil
}

func (m managerRepository) ViewUsersInfo(ctx context.Context) ([]*mongo_models.User, error) {
	usersCollection := m.DB.Database(os.Getenv("MONGO_DATABASE")).Collection("users")

	cursorUsers, err := usersCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		panic(err)
	}

	var users []*mongo_models.User
	if err = cursorUsers.All(context.TODO(), &users); err != nil {
		panic(err)
	}

	return users, nil
}

func (m managerRepository) ViewUserInfo(ctx context.Context, telegramUserId int64) (*slack_bot.UserInfo, error) {
	var userInfo = &slack_bot.UserInfo{}
	// получить базовую информацию
	var user *mongo_models.User
	err := m.DB.Database(os.Getenv("MONGO_DATABASE")).Collection("users").FindOne(ctx, bson.D{
		{"telegram_user_id", telegramUserId},
	}).Decode(&user)
	if err != nil {
		return nil, err
	}

	userInfo.TelegramUserId = user.TelegramUserID
	userInfo.Status = user.Status
	userInfo.CreatedAt = user.CreatedAt

	// если статус премиум, запросить из бд сколько часов осталось
	if user.Status == "premium" {
		var userSubscriptionTime mongo_models.UserSubscriptionTime

		err := m.DB.Database(os.Getenv("MONGO_DATABASE")).Collection("users_subscription_time").FindOne(ctx, bson.D{
			{"telegram_user_id", telegramUserId},
		}).Decode(&userSubscriptionTime)
		if err != nil {
			return nil, err
		}

		userInfo.PremiumDuration = userSubscriptionTime.ActiveBefore.Sub(time.Now()).Hours()
	} else {
		userInfo.PremiumDuration = 0
	}

	// история активации промокодов
	var promoCodeActivations []mongo_models.PromoCodeActivations

	promoCodeActivationsCursor, err := m.DB.Database(os.Getenv("MONGO_DATABASE")).Collection("promo_codes_activations").Find(ctx, bson.D{
		{"telegram_user_id", telegramUserId},
	})
	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			promoCodeActivations = nil
		default:
			log.Printf("Could not find promocode activations for user with _id: %v.Reason: %v\n", user.TelegramUserID, err)
			return nil, apperrors.NewInternal()
		}
	} else {
		if err = promoCodeActivationsCursor.All(context.TODO(), &promoCodeActivations); err != nil {
			panic(err)
		}
	}

	userInfo.PromoCodeActivations = promoCodeActivations

	// обсерверы на цены
	var priceObservers []mongo_models.Observer

	priceObserversCursor, err := m.DB.Database(os.Getenv("MONGO_DATABASE")).Collection("observers").Find(ctx, bson.D{
		{"telegram_user_id", telegramUserId},
	})
	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			priceObservers = nil
		default:
			log.Printf("Could not find price observers for user with _id: %v.Reason: %v\n", user.TelegramUserID, err)
			return nil, apperrors.NewInternal()
		}
	} else {
		if err = priceObserversCursor.All(context.TODO(), &priceObservers); err != nil {
			panic(err)
		}
	}

	userInfo.PriceObservers = priceObservers

	// информация по percentage observers
	var percentageObserver mongo_models.PercentageObserver

	err = m.DB.Database(os.Getenv("MONGO_DATABASE")).Collection("percentage_observers").FindOne(ctx, bson.D{
		{"telegram_user_id", telegramUserId},
	}).Decode(&percentageObserver)
	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			percentageObserver = mongo_models.PercentageObserver{}
		default:
			log.Printf("Could not find percentage observers for user with _id: %v.Reason: %v\n", user.TelegramUserID, err)
			return nil, apperrors.NewInternal()
		}
	}

	userInfo.PercentageObserver = percentageObserver

	// portfolio
	var portfolio []mongo_models.Portfolio

	portfolioCursor, err := m.DB.Database(os.Getenv("MONGO_DATABASE")).Collection("portfolio").Find(ctx, bson.D{
		{"telegram_user_id", telegramUserId},
	})
	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			portfolio = nil
		default:
			log.Printf("Could not find portfolio for user with _id: %v.Reason: %v\n", user.TelegramUserID, err)
			return nil, apperrors.NewInternal()
		}
	} else {
		if err = portfolioCursor.All(context.TODO(), &portfolio); err != nil {
			panic(err)
		}
	}

	userInfo.Portfolio = portfolio

	// количество requests
	userInfo.CountRequestsLast24H = 0
	userInfo.CountRequestsLast7D = 0
	userInfo.CountRequestsLastMonth = 0
	userInfo.CountRequests = 0

	var logRequests []*mongo_models.LogRequest

	logRequestsCollection := m.DB.Database(os.Getenv("MONGO_DATABASE")).Collection("log_requests")
	cursorOnlineUsersLogs, err := logRequestsCollection.Find(ctx, bson.D{
		{"telegram_user_id", telegramUserId},
	})
	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			portfolio = nil
		default:
			log.Printf("Could not find log for user with _id: %v.Reason: %v\n", user.TelegramUserID, err)
			return nil, apperrors.NewInternal()
		}
	} else {
		if err = cursorOnlineUsersLogs.All(context.TODO(), &logRequests); err != nil {
			panic(err)
		}
	}

	for _, logRequest := range logRequests {
		if logRequest.CreatedAt.Sub(time.Now()).Hours()*-1 <= 24 {
			log.Println(logRequest.CreatedAt.Sub(time.Now()).Hours())
			userInfo.CountRequestsLast24H += 1
			userInfo.CountRequestsLast7D += 1
			userInfo.CountRequestsLastMonth += 1
			userInfo.CountRequests += 1
		} else if logRequest.CreatedAt.Sub(time.Now()).Hours()*-1 <= 24*7 {
			userInfo.CountRequestsLast7D += 1
			userInfo.CountRequestsLastMonth += 1
			userInfo.CountRequests += 1
		} else if logRequest.CreatedAt.Sub(time.Now()).Hours()*-1 <= 24*7*30 {
			userInfo.CountRequestsLastMonth += 1
			userInfo.CountRequests += 1
		} else {
			userInfo.CountRequests += 1
		}
	}

	return userInfo, nil
}

func (m managerRepository) ViewActiveObservers(ctx context.Context) ([]*mongo_models.Observer, error) {
	observerRequestsCollection := m.DB.Database(os.Getenv("MONGO_DATABASE")).Collection("observers")
	cursorObservers, err := observerRequestsCollection.Find(ctx, bson.D{
		{"is_active", true},
	})
	if err != nil {
		return nil, err
	}

	var activeObservers []*mongo_models.Observer
	if err = cursorObservers.All(context.TODO(), &activeObservers); err != nil {
		panic(err)
	}

	return activeObservers, nil
}

func (m managerRepository) AddUserToAdminGroup(ctx context.Context, telegramUserId int64) error {
	// save to bd
	_, err := m.DB.Database(os.Getenv("MONGO_DATABASE")).Collection("admins").InsertOne(ctx, bson.D{
		{"telegram_user_id", telegramUserId},
	})
	if err != nil {
		log.Printf("Could not create a new user admin: %v.Reason: %v\n", telegramUserId, err)
		return apperrors.NewInternal()
	}

	return nil
}

func contains(s []int64, e int64) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
