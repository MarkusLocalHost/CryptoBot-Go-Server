package mongo_models

import "go.mongodb.org/mongo-driver/bson/primitive"

type CronSubscriptionIds struct {
	Id             primitive.ObjectID `json:"id"            bson:"_id"`
	CronId         int                `json:"cronId"        bson:"cron_id"`
	TelegramUserID primitive.ObjectID `json:"observerId"    bson:"telegram_user_id"`
	TimeToRestart  string             `json:"timeToRestart" bson:"time_to_restart"`
}
