package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDatabase() (*mongo.Client, context.Context) {
	client, err := mongo.NewClient(options.Client().ApplyURI(
		"mongodb+srv://admin:Branstark1@production.tobvq.mongodb.net/pantry?retryWrites=true&w=majority",
	))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	return client, ctx
}

func CreateDatabase(databaseName string, client *mongo.Client) *mongo.Database {
	database := client.Database(databaseName)
	return database
}

func CreateCollection(collectionName string, database *mongo.Database) *mongo.Collection {
	collection := database.Collection(collectionName)
	return collection
}

func InsertOneData(collection *mongo.Collection, data []bson.D, ctx context.Context) {
	_, err := collection.InsertOne(ctx, data)
	if err != nil {
		log.Fatal(err)
	}
}

func DeleteData(collection *mongo.Collection, ctx context.Context) {
	collection.Drop(ctx)
}
