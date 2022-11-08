package mongo_models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type UserSubscriptionTime struct {
	Id             primitive.ObjectID `json:"id"                bson:"_id"`
	TelegramUserId int64              `json:"telegramUserId"    bson:"telegram_user_id"`
	ActiveBefore   time.Time          `json:"activeBefore"      bson:"active_before"`
}
