package observer_func

import (
	"context"
	"cryptocurrency/internal/models/mongo_models"
	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
)

func UpdateCronTaskAfterSignal(mongoDB *mongo.Client, observer *mongo_models.Observer, cronObservers *cron.Cron) error {
	var deletedDocument bson.M
	err := mongoDB.Database(os.Getenv("MONGO_DATABASE")).Collection("cron_observer_ids").FindOneAndDelete(context.Background(),
		bson.D{{"observer_id", observer.Id}}).Decode(&deletedDocument)
	if err != nil {
		return err
	}

	cronId := int(deletedDocument["cron_id"].(int32))
	cronObservers.Remove(cron.EntryID(cronId))

	return nil
}
