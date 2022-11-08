package repository

import (
	"context"
	"cryptocurrency/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
	"time"
)

type logRepository struct {
	DB *mongo.Client
}

func NewLogRepository(db *mongo.Client) models.LogRepository {
	return &logRepository{
		DB: db,
	}
}

func (l logRepository) LogRequestFromBot(ctx context.Context, path string, telegramUserId int64, data interface{}) error {
	_, err := l.DB.Database(os.Getenv("MONGO_DATABASE")).Collection("log_requests").InsertOne(ctx, bson.D{
		{"path", path},
		{"telegram_user_id", telegramUserId},
		{"data", data},
		{"created_at", time.Now()},
	})
	if err != nil {
		panic(err)
	}

	return nil
}

func (l logRepository) LogRequestFromSite(ctx context.Context, path string, userID string, data interface{}) error {
	_, err := l.DB.Database(os.Getenv("MONGO_DATABASE")).Collection("log_requests").InsertOne(ctx, bson.D{
		{"path", path},
		{"user_id", userID},
		{"data", data},
		{"created_at", time.Now()},
	})
	if err != nil {
		panic(err)
	}

	return nil
}
