package observer_func

import (
	"context"
	"cryptocurrency/internal/models/mongo_models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
)

func DeleteObserverFromBD(mongoDB *mongo.Client, observer *mongo_models.Observer) error {
	_, err := mongoDB.Database(os.Getenv("MONGO_DATABASE")).Collection("observers").UpdateOne(context.Background(), bson.D{
		{"_id", observer.Id},
	}, bson.D{
		{"$set", bson.D{{"is_active", false}}},
	})
	if err != nil {
		return err
	}

	return nil
}
