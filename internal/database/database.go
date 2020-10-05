package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id, omitempty"`
	Firstname string             `bson:"firstname, omitempty"`
	Lastname  string             `bson:"lastname, omitempty"`
	Email     string             `bson:"email, omitempty"`
}

type IngredientList struct {
	ID         primitive.ObjectID `bson:"_id, omitempty"`
	User       primitive.ObjectID `bson:"user, omitempty"`
	Ingredient []Ingredient       `bson:"ingredient"`
}

type Ingredient struct {
	//ID   primitive.ObjectID `bson:"_id"`
	Name string `bson:"name, omitempty"`
}

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

func InsertDataToCollection(collection *mongo.Collection, ctx context.Context, data primitive.D) {
	_, err := collection.InsertOne(ctx, data)
	if err != nil {
		log.Fatal(err)
	}
}
