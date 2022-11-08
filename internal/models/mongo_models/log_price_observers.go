package mongo_models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type LogPriceObservers struct {
	Id         primitive.ObjectID `json:"id"         bson:"_id"`
	ObserverId primitive.ObjectID `json:"observerId" bson:"observer_id"`
	ObservedAt time.Time          `json:"observedAt" bson:"observed_at"`
}
