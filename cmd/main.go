package main

import (
	"fmt"
	"log"

	"github.com/zacharygilliom/internal/database"
	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	// Establish our connection to our databse
	client, ctx := database.ConnectDatabase()
	// Connect to our database
	err := client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	// Disconnects our connection after finishes running
	defer client.Disconnect(ctx)

	//err = client.Ping(ctx, readpref.Primary())
	//if err != nil {
	//	log.Fatal(err)
	//}
	// Create our Database
	var databaseName string = "pantry"
	pantryDatabase := database.CreateDatabase(databaseName, client)

	// Create a Collection
	var collectionName string = "ingredients"
	ingredientsCollection := database.CreateCollection(collectionName, pantryDatabase)

	//ingredientsResults, err := ingredientsCollection.InsertOne(ctx, bson.D{
	//{Key: "name", Value: "flour"},
	//{Key: "kind", Value: "white"},
	//})

	results := struct {
		name string
		kind string
	}{}

	fmt.Println(results)

	cursor, err := ingredientsCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	var ingredients []bson.M
	if err = cursor.All(ctx, &ingredients); err != nil {
		log.Fatal(err)
	}
	for _, ingredient := range ingredients {
		fmt.Println(ingredient)
	}

	database.DeleteData(ingredientsCollection, ctx)

	for _, ingredient := range ingredients {
		fmt.Println(ingredient)
	}

}
