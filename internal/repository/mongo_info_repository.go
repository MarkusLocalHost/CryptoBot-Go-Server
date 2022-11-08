package repository

import (
	"context"
	"cryptocurrency/internal/models"
	"cryptocurrency/internal/models/mongo_models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
)

type mongoInfoRepository struct {
	DB *mongo.Client
}

func NewInfoRepository(db *mongo.Client) models.InfoRepository {
	return &mongoInfoRepository{
		DB: db,
	}
}

func (m mongoInfoRepository) GetAllUsersLanguages(ctx context.Context) (map[int64]string, error) {
	coll := m.DB.Database(os.Getenv("MONGO_DATABASE")).Collection("users")
	filter := bson.D{{}}

	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		panic(err)
	}

	var results []mongo_models.User
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	usersLanguages := make(map[int64]string)
	for _, result := range results {
		usersLanguages[result.TelegramUserID] = result.Language
	}

	return usersLanguages, nil
}

func (m mongoInfoRepository) GetAllUsersAdmins(ctx context.Context) ([]int64, error) {
	coll := m.DB.Database(os.Getenv("MONGO_DATABASE")).Collection("admins")
	filter := bson.D{{}}

	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		panic(err)
	}

	var results []int64
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	return results, nil
}

func (m mongoInfoRepository) TryFindCurrencyByName(ctx context.Context, name string) ([]mongo_models.TryFindCurrency, error) {
	coll := m.DB.Database(os.Getenv("MONGO_DATABASE")).Collection("currencies")
	filter := bson.D{{"slug", primitive.Regex{Pattern: name, Options: ""}}}

	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		panic(err)
	}

	var results []mongo_models.TryFindCurrency
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	return results, nil
}

func (m mongoInfoRepository) TryFindCurrencyBySlug(ctx context.Context, slug string) ([]mongo_models.TryFindCurrency, error) {
	coll := m.DB.Database(os.Getenv("MONGO_DATABASE")).Collection("currencies")
	filter := bson.D{{"symbol", primitive.Regex{Pattern: slug, Options: ""}}}

	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		panic(err)
	}

	var results []mongo_models.TryFindCurrency
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	return results, nil
}

func (m mongoInfoRepository) GetSupportedVSCurrencies(ctx context.Context) ([]string, error) {
	coll := m.DB.Database(os.Getenv("MONGO_DATABASE")).Collection("supported_vs_currencies")
	filter := bson.D{{"is_primary", true}}

	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		panic(err)
	}

	type result struct {
		CurrencyName string `bson:"currency_name"`
	}
	var resultsFromBD []result
	if err = cursor.All(context.TODO(), &resultsFromBD); err != nil {
		panic(err)
	}

	var results []string
	for _, currency := range resultsFromBD {
		results = append(results, currency.CurrencyName)
	}

	return results, nil
}

func (m mongoInfoRepository) GetSymbolAndNameByIDFromCoinGecko(ctx context.Context, id string) (symbol, name string, err error) {
	var result struct {
		Id     string `bson:"id"`
		Symbol string `bson:"symbol"`
		Name   string `bson:"name"`
	}

	coll := m.DB.Database(os.Getenv("MONGO_DATABASE")).Collection("coingecko_currencies")
	err = coll.FindOne(ctx, bson.D{
		{"id", id},
	}).Decode(&result)
	if err != nil {
		return "", "", err
	}

	return result.Symbol, result.Name, nil
}
