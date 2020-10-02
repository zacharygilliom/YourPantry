package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateConnection() (*mongo.Client, context.Context, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://admin:Branstark1@production.tobvq.mongodb.net/pantry?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	return client, ctx, err
}

func CreateDatabase(databaseName string, client *mongo.Client) *mongo.Database {
	database := client.Database(databaseName)
	return database
}

func CreateCollection(collectionName string, database *mongo.Database) *mongo.Collection {
	collection := database.Collection(collectionName)
	return collection
}
