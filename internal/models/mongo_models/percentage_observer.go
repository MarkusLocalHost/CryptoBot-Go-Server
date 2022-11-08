package mongo_models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type PercentageObserver struct {
	Id                 primitive.ObjectID `json:"id"                   bson:"_id"`
	TelegramUserID     int64              `json:"telegram_user_id"     bson:"telegram_user_id"`
	Observe20Minutes   bool               `json:"observe_20_minutes"   bson:"observe_20_minutes"`
	Observe60Minutes   bool               `json:"observe_60_minutes"   bson:"observe_60_minutes"`
	FirstFilterType    string             `json:"first_filter_type"    bson:"first_filter_type"`
	FirstFilterAmount  float64            `json:"first_filter_amount"  bson:"first_filter_amount"`
	SecondFilterType   string             `json:"second_filter_type"   bson:"second_filter_type"`
	SecondFilterAmount float64            `json:"second_filter_amount" bson:"second_filter_amount "`
	CreatedAt          time.Time          `json:"createdAt"            bson:"created_at"`
}
