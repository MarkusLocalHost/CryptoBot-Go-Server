package main

import (
	"context"
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
)

type dataSource struct {
	MongoDBClient *mongo.Client
	BadgerClient  *badger.DB
}

func initDS() (*dataSource, error) {
	// Initialize MongoDB Connection
	ctx := context.Background()
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err = mongoClient.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")

	// Initialize BadgerDB connection
	opts := badger.DefaultOptions("./../../tmp/badger")
	opts.Logger = nil

	badgerDB, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}

	return &dataSource{
		MongoDBClient: mongoClient,
		BadgerClient:  badgerDB,
	}, nil
}

func (d *dataSource) closeDS() error {
	if err := d.MongoDBClient.Disconnect(context.Background()); err != nil {
		return fmt.Errorf("error closing MongoDB: %w", err)
	}

	err := d.BadgerClient.Close()
	if err != nil {
		return fmt.Errorf("error closing BadgerDB: %w", err)
	}

	return nil
}
