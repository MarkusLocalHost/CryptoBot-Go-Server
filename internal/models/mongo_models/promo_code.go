package mongo_models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type PromoCode struct {
	Id                primitive.ObjectID `json:"id"                bson:"_id"`
	Title             string             `json:"title"             bson:"title"`
	Value             string             `json:"value"             bson:"value"`
	SubscriptionHours string             `json:"subscriptionHours" bson:"subscription_hours"`
	CountOfActivation int                `json:"countOfActivation" bson:"count_of_activation"`
	CreatedAt         time.Time          `json:"createdAt"         bson:"created_at"`
	ActiveBefore      time.Time          `json:"activeBefore"      bson:"active_before"`
}
