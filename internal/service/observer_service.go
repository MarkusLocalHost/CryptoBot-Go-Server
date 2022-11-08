package service

import (
	"context"
	"cryptocurrency/internal/models"
	"cryptocurrency/internal/models/mongo_models"
	"cryptocurrency/internal/utils/apperrors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
)

type observerService struct {
	ObserverRepository models.ObserverRepository
	InfoRepository     models.InfoRepository
}

type OSConfig struct {
	ObserverRepository models.ObserverRepository
	InfoRepository     models.InfoRepository
}

func NewObserverService(c *OSConfig) models.ObserverService {
	return &observerService{
		ObserverRepository: c.ObserverRepository,
		InfoRepository:     c.InfoRepository,
	}
}

func (o observerService) CreatePriceObserver(ctx context.Context, observer *mongo_models.Observer) (observerIsActive bool, err error) {
	cryptoSymbol, cryptoName, err := o.InfoRepository.GetSymbolAndNameByIDFromCoinGecko(ctx, observer.CryptoID)
	if err != nil {
		log.Fatal(err)
	}
	observer.CryptoSymbol = cryptoSymbol
	observer.CryptoName = cryptoName

	userTypeSubscription, observerIsActive, err := o.ObserverRepository.CreatePriceObserver(ctx, observer)
	if err != nil {
		return false, apperrors.NewInternal()
	}

	//Initialize Cron tasks
	if observerIsActive {
		err = o.ObserverRepository.MakeCronTask(ctx, observer, userTypeSubscription)
		if err != nil {
			log.Fatalf("Failure to start cron tasks: %v\n", err)
		}
		log.Println("Cron task started")
	}

	return observerIsActive, nil
}

func (o observerService) DeletePriceObserver(ctx context.Context, observerId primitive.ObjectID) error {
	err := o.ObserverRepository.DeletePriceObserver(ctx, observerId)
	if err != nil {
		return apperrors.NewInternal()
	}

	err = o.ObserverRepository.RemoveCronTask(ctx, observerId)
	if err != nil {
		return apperrors.NewInternal()
	}

	return nil
}

func (o observerService) ChangeStatusPriceObserver(ctx context.Context, observerId primitive.ObjectID) (status string, error error) {
	status, err := o.ObserverRepository.ChangeStatusPriceObserver(ctx, observerId)
	if err != nil {
		return "", apperrors.NewInternal()
	}

	return status, nil
}

func (o observerService) CreatePercentageObserver(ctx context.Context, percentageObserver *mongo_models.PercentageObserver) (status string, err error) {
	status, err = o.ObserverRepository.CreatePercentageObserver(ctx, percentageObserver)
	if err != nil {
		return "", apperrors.NewInternal()
	}

	return status, nil
}

func (o observerService) DeletePercentageObserver(ctx context.Context, observerId primitive.ObjectID) error {
	err := o.ObserverRepository.DeletePercentageObserver(ctx, observerId)
	if err != nil {
		return apperrors.NewInternal()
	}

	return nil
}
