package mongo_models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Portfolio struct {
	Id             primitive.ObjectID `json:"id"              bson:"_id"`
	Cryptocurrency string             `json:"cryptocurrency"  bson:"cryptocurrency"`
	TelegramUserID int64              `json:"telegramUserID"  bson:"telegram_user_id"`
	Value          float64            `json:"value"           bson:"value"`
	Price          float64            `json:"price"           bson:"price"`
	ActualPrice    float64            `json:"actualPrice"     bson:"-"`
	Type           string             `json:"type"            bson:"type"`
	CreatedAt      time.Time          `json:"createdAt"       bson:"created_at"`
}
