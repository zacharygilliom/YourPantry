package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	client, err := mongo.NewClient(options.Client().ApplyURI(
		"mongodb+srv://admin:Branstark1@production.tobvq.mongodb.net/pantry?retryWrites=true&w=majority",
	))
	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	pantry := client.Database("pantry")

	ingredientsCollection := pantry.Collection("ingredients")

	ingredientsResults, err := ingredientsCollection.InsertOne(ctx, bson.D{
		{Key: "name", Value: "flour"},
		{Key: "kind", Value: "white"},
	})

	results := struct {
		name string
		kind string
	}{}

	filter := bson.M{"kind": "white"}
	// filter := bson.M{"name": "flour"}
	err = ingredientsCollection.FindOne(context.Background(), filter).Decode(&results)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(ingredientsResults)
	fmt.Println(results)
}
