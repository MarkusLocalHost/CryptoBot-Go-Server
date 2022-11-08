package mongo_models

import (
	"time"
)

type User struct {
	TelegramUserID int64     `json:"telegramUserID"  bson:"telegram_user_id"`
	Status         string    `json:"status"          bson:"status"`
	Language       string    `json:"language"        bson:"language"`
	CreatedAt      time.Time `json:"createdAt"       bson:"created_at"`
}
