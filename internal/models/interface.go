package models

import (
	"context"
	"cryptocurrency/internal/models/mongo_models"
	"cryptocurrency/internal/models/response_models"
	"cryptocurrency/internal/models/response_models/slack_bot"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type ObserverService interface {
	CreatePriceObserver(ctx context.Context, observer *mongo_models.Observer) (observerIsActive bool, err error)
	DeletePriceObserver(ctx context.Context, observerId primitive.ObjectID) error
	ChangeStatusPriceObserver(ctx context.Context, observerId primitive.ObjectID) (status string, err error)

	CreatePercentageObserver(ctx context.Context, percentageObserver *mongo_models.PercentageObserver) (status string, err error)
	DeletePercentageObserver(ctx context.Context, observerId primitive.ObjectID) error
}

type ObserverRepository interface {
	CreatePriceObserver(ctx context.Context, observer *mongo_models.Observer) (userTypeSubscription string, observerIsActive bool, err error)
	DeletePriceObserver(ctx context.Context, observerId primitive.ObjectID) error
	ChangeStatusPriceObserver(ctx context.Context, observerId primitive.ObjectID) (status string, err error)

	MakeCronTask(ctx context.Context, observer *mongo_models.Observer, userTypeSubscription string) error
	RemoveCronTask(ctx context.Context, observerId primitive.ObjectID) error

	CreatePercentageObserver(ctx context.Context, percentageObserver *mongo_models.PercentageObserver) (status string, err error)
	DeletePercentageObserver(ctx context.Context, observerId primitive.ObjectID) error
}

type InfoService interface {
	GetAllUsersLanguages(ctx context.Context) (map[int64]string, error)
	GetAllUsersAdmins(ctx context.Context) ([]int64, error)

	GetBasicCurrencyInfo(ctx context.Context, currencyID string) (response_models.BasicCurrencyInfoFullVersion, error)
	GetBasicCurrencyInfoShortVersion(ctx context.Context, currencyID string) (response_models.BasicCurrencyInfoShortVersion, error)

	GetTrendingCurrencies(ctx context.Context, source string) ([]response_models.TrendingCurrency, error)

	TryFindCurrencyByNameOrSlug(ctx context.Context, name string, slug string) (mongo_models.TryFindCurrency, error)
	GetSupportedVSCurrencies(ctx context.Context) ([]string, error)

	GetIndexPriceByPage(ctx context.Context, page int, currency string) ([]response_models.CurrencyInfoIndexVersion, error)

	GetPriceForSymbolCurrenciesInPortfolio([]*mongo_models.Portfolio) ([]*mongo_models.Portfolio, error)

	GetExchangeRateFromBestchange(ctx context.Context, from string, fromType string, to string, toType string, limitCurrency string, limitValue float64) (map[string][]response_models.ExchangeRateVariant, error)
}

type InfoRepository interface {
	GetAllUsersLanguages(ctx context.Context) (map[int64]string, error)
	GetAllUsersAdmins(ctx context.Context) ([]int64, error)

	TryFindCurrencyByName(ctx context.Context, name string) ([]mongo_models.TryFindCurrency, error)
	TryFindCurrencyBySlug(ctx context.Context, slug string) ([]mongo_models.TryFindCurrency, error)

	GetSupportedVSCurrencies(ctx context.Context) ([]string, error)

	GetSymbolAndNameByIDFromCoinGecko(ctx context.Context, id string) (symbol, name string, err error)
}

type LogService interface {
	LogRequestFromBot(ctx context.Context, path string, telegramUserId int64, data interface{}) error
	LogRequestFromSite(ctx context.Context, path string, userID string, data interface{}) error
}

type LogRepository interface {
	LogRequestFromBot(ctx context.Context, path string, telegramUserId int64, data interface{}) error
	LogRequestFromSite(ctx context.Context, path string, userID string, data interface{}) error
}

type AccountService interface {
	CreateAccount(ctx context.Context, telegramUserId int64, language string) error

	GetUserPriceObservers(ctx context.Context, telegramUserId int64) (observers []*mongo_models.Observer, err error)
	GetUserPercentageObservers(ctx context.Context, telegramUserId int64) (observers []*mongo_models.PercentageObserver, err error)

	ViewUserSubscription(ctx context.Context, telegramUserId int64) (string, error)
	ViewCountOfUserActiveObserver(ctx context.Context, telegramUserId int64) (int, int, int, error)

	ViewUserPortfolio(ctx context.Context, telegramUserId int64) ([]*mongo_models.Portfolio, error)
	AddToUserPortfolio(ctx context.Context, data *mongo_models.Portfolio) error
	UpdateElementUserPortfolio(ctx context.Context, id primitive.ObjectID, value float64, price float64) error
	DeleteElementUserPortfolio(ctx context.Context, id primitive.ObjectID) error

	CheckPromoCode(ctx context.Context, promoCode string, telegramUserId int64) (string, error)

	ExtendUserSubscription(ctx context.Context, telegramUserId int64, hours int) error
}

type AccountRepository interface {
	CreateAccount(ctx context.Context, telegramUserId int64, language string) error

	GetUserPriceObservers(ctx context.Context, telegramUserId int64) (observers []*mongo_models.Observer, err error)
	GetUserPercentageObservers(ctx context.Context, telegramUserId int64) (observers []*mongo_models.PercentageObserver, err error)

	ViewUserSubscription(ctx context.Context, telegramUserId int64) (string, error)
	ViewCountOfUserActiveObserver(ctx context.Context, telegramUserId int64) (int, int, int, error)

	ViewUserPortfolio(ctx context.Context, telegramUserId int64) ([]*mongo_models.Portfolio, error)
	AddToUserPortfolio(ctx context.Context, data *mongo_models.Portfolio) error
	UpdateElementPriceUserPortfolio(ctx context.Context, id primitive.ObjectID, price float64) error
	UpdateElementValueUserPortfolio(ctx context.Context, id primitive.ObjectID, value float64) error
	DeleteElementUserPortfolio(ctx context.Context, id primitive.ObjectID) error

	CheckPromoCode(ctx context.Context, promoCode string, telegramUserId int64) (string, error)

	ChangeUserSubscription(ctx context.Context, telegramUserId int64, status string) error
	AddTimeToUserSubscription(ctx context.Context, telegramUserId int64, hoursToAdd time.Duration) error

	GetUsersByFilter(ctx context.Context, filter string) ([]mongo_models.User, error)
}

type TokenService interface {
	ValidateIDToken(tokenString string) (string, error)
}

type ManagerService interface {
	CreatePromoCode(ctx context.Context, promoCode *mongo_models.PromoCode) error
	ViewPromoCodes(ctx context.Context) ([]*slack_bot.PromoCodesView, error)

	AddHoursToUserSubscription(ctx context.Context, telegramUserId int64, hours int) error

	ViewOnlineUsersStats(ctx context.Context) ([]*slack_bot.OnlineUserStats, error)
	ViewCountRequests(ctx context.Context) ([]*slack_bot.LogRequestStats, error)
	ViewUserActions(ctx context.Context, telegramUserId int64) ([]*mongo_models.LogRequest, error)
	ViewUsersInfo(ctx context.Context) ([]*mongo_models.User, error)
	ViewUserInfo(ctx context.Context, telegramUserId int64) (*slack_bot.UserInfo, error)
	ViewActiveObservers(ctx context.Context) ([]*mongo_models.Observer, error)

	SendMessage(ctx context.Context, filter string, message string) error

	AddUserToAdminGroup(ctx context.Context, telegramUserId int64) error
}

type ManagerRepository interface {
	CreatePromoCode(ctx context.Context, promoCode *mongo_models.PromoCode) error
	ViewPromoCodes(ctx context.Context) ([]*slack_bot.PromoCodesView, error)

	ViewOnlineUsersStats(ctx context.Context) ([]*slack_bot.OnlineUserStats, error)
	ViewCountRequests(ctx context.Context) ([]*slack_bot.LogRequestStats, error)
	ViewUserActions(ctx context.Context, telegramUserId int64) ([]*mongo_models.LogRequest, error)
	ViewUsersInfo(ctx context.Context) ([]*mongo_models.User, error)
	ViewUserInfo(ctx context.Context, telegramUserId int64) (*slack_bot.UserInfo, error)
	ViewActiveObservers(ctx context.Context) ([]*mongo_models.Observer, error)

	AddUserToAdminGroup(ctx context.Context, telegramUserId int64) error
}
