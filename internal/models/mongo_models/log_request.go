package mongo_models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type LogRequest struct {
	Id             primitive.ObjectID `json:"id" bson:"_id"`
	Path           string             `json:"path" bson:"path"`
	TelegramUserId int64              `json:"telegramUserId" bson:"telegram_user_id"`
	Data           interface{}        `json:"data" bson:"data"`
	CreatedAt      time.Time          `json:"createdAt" bson:"created_at"`
}
