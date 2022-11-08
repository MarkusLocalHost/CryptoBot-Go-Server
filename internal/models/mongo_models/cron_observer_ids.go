package mongo_models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CronObserverIds struct {
	Id            primitive.ObjectID `json:"id"            bson:"_id"`
	CronId        int                `json:"cronId"        bson:"cron_id"`
	ObserverId    primitive.ObjectID `json:"observerId"    bson:"observer_id"`
	TimeToRestart string             `json:"timeToRestart" bson:"time_to_restart"`
}
