package main

import (
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
	//ingredientCollection := "ingredient"
	//pantryIngredient := database.CreateCollection(ingredientCollection, pantryDatabase)

	userCollection := "user"
	pantryUser := database.CreateCollection(userCollection, pantryDatabase)

	data := bson.D{
		{"firstname", "Zach"},
		{"lastname", "Gilliom"},
		{"email", "zacharygilliom@gmail.com"},
	}
	database.InsertDataToCollection(pantryUser, ctx, data)
	//pantryResult, err := pantryIngredient.InsertOne(ctx, bson.D{
	//{"name", "flour"},
	//})
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(pantryResult)

}

func getUserInput() string {
	var answer string
	return answer
}
