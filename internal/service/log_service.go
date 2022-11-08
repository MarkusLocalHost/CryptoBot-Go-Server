package service

import (
	"context"
	"cryptocurrency/internal/models"
	"log"
)

type logService struct {
	LogRepository models.LogRepository
}

type LSConfig struct {
	LogRepository models.LogRepository
}

func NewLogService(c *LSConfig) models.LogService {
	return &logService{
		LogRepository: c.LogRepository,
	}
}

func (l logService) LogRequestFromBot(ctx context.Context, path string, telegramUserId int64, data interface{}) error {
	err := l.LogRepository.LogRequestFromBot(ctx, path, telegramUserId, data)
	if err != nil {
		log.Fatalf("Error to log data: %s", err)
	}

	return nil
}

func (l logService) LogRequestFromSite(ctx context.Context, path string, userID string, data interface{}) error {
	err := l.LogRepository.LogRequestFromSite(ctx, path, userID, data)
	if err != nil {
		log.Fatalf("Error to log data: %s", err)
	}

	return nil
}
