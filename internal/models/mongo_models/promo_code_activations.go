package mongo_models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type PromoCodeActivations struct {
	Id             primitive.ObjectID `json:"id"                bson:"_id"`
	PromoCodeId    primitive.ObjectID `json:"promoCodeId"       bson:"promo_code_id"`
	TelegramUserId int64              `json:"telegramUserId"    bson:"telegram_user_id"`
	ActivatedAt    time.Time          `json:"activatedAt"       bson:"activated_at"`
}
