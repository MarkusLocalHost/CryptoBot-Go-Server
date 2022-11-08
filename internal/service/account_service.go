package service

import (
	"context"
	"cryptocurrency/internal/models"
	"cryptocurrency/internal/models/mongo_models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type accountService struct {
	AccountRepository models.AccountRepository
}

type ASConfig struct {
	AccountRepository models.AccountRepository
}

func NewAccountService(c *ASConfig) models.AccountService {
	return &accountService{
		AccountRepository: c.AccountRepository,
	}
}

func (a accountService) CreateAccount(ctx context.Context, telegramUserId int64, language string) error {
	err := a.AccountRepository.CreateAccount(ctx, telegramUserId, language)
	if err != nil {
		return err
	}

	return nil
}

func (a accountService) GetUserPriceObservers(ctx context.Context, telegramUserId int64) (observer []*mongo_models.Observer, err error) {
	data, err := a.AccountRepository.GetUserPriceObservers(ctx, telegramUserId)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (a accountService) GetUserPercentageObservers(ctx context.Context, telegramUserId int64) (observers []*mongo_models.PercentageObserver, err error) {
	data, err := a.AccountRepository.GetUserPercentageObservers(ctx, telegramUserId)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (a accountService) ViewUserSubscription(ctx context.Context, telegramUserId int64) (string, error) {
	data, err := a.AccountRepository.ViewUserSubscription(ctx, telegramUserId)
	if err != nil {
		return "", err
	}

	return data, nil
}

func (a accountService) ViewCountOfUserActiveObserver(ctx context.Context, telegramUserId int64) (int, int, int, error) {
	dataTier1, dataTier2, dataPercentage, err := a.AccountRepository.ViewCountOfUserActiveObserver(ctx, telegramUserId)
	if err != nil {
		return 0, 0, 0, err
	}

	return dataTier1, dataTier2, dataPercentage, nil
}

func (a accountService) ViewUserPortfolio(ctx context.Context, telegramUserId int64) ([]*mongo_models.Portfolio, error) {
	data, err := a.AccountRepository.ViewUserPortfolio(ctx, telegramUserId)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (a accountService) AddToUserPortfolio(ctx context.Context, data *mongo_models.Portfolio) error {
	err := a.AccountRepository.AddToUserPortfolio(ctx, data)
	if err != nil {
		return err
	}

	return nil
}

func (a accountService) UpdateElementUserPortfolio(ctx context.Context, id primitive.ObjectID, value float64, price float64) error {
	if value != 0 && price == 0 {
		err := a.AccountRepository.UpdateElementValueUserPortfolio(ctx, id, value)
		if err != nil {
			return err
		}
	} else if value == 0 && price != 0 {
		err := a.AccountRepository.UpdateElementPriceUserPortfolio(ctx, id, price)
		if err != nil {
			return err
		}
	}

	return nil

}

func (a accountService) DeleteElementUserPortfolio(ctx context.Context, id primitive.ObjectID) error {
	err := a.AccountRepository.DeleteElementUserPortfolio(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (a accountService) CheckPromoCode(ctx context.Context, promoCode string, telegramUserId int64) (string, error) {
	status, err := a.AccountRepository.CheckPromoCode(ctx, promoCode, telegramUserId)
	if err != nil {
		return "", err
	}

	return status, nil
}

func (a accountService) ExtendUserSubscription(ctx context.Context, telegramUserId int64, hours int) error {
	err := a.AccountRepository.ChangeUserSubscription(ctx, telegramUserId, "premium")
	if err != nil {
		return err
	}

	err = a.AccountRepository.AddTimeToUserSubscription(ctx, telegramUserId, time.Duration(hours*1000000000*60*60))
	if err != nil {
		return err
	}

	return nil
}
