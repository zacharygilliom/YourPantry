package main

import (
	"fmt"
	"log"

	"github.com/zacharygilliom/internal/database"
	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	client, ctx, err := database.CreateConnection()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	//pantryDatabase := client.Database("pantry")
	databaseName := "pantry"
	pantryDatabase := database.CreateDatabase(databaseName, client)

	//pantryIngredient := pantryDatabase.Collection("ingredient")
	collectionName := "ingredient"
	pantryIngredient := database.CreateCollection(collectionName, pantryDatabase)

	pantryResult, err := pantryIngredient.InsertOne(ctx, bson.D{
		{"name", "flour"},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pantryResult.InsertedID)

	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(databases)

}
