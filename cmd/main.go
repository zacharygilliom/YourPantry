package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/zacharygilliom/internal/database"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	fmt.Println("Connection Started...")
	client, ctx, err := database.CreateConnection()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	//pantryDatabase := client.Database("pantry")
	databaseName := "pantry"
	pantryDatabase := database.NewDatabase(databaseName, client)

	//pantryIngredient := pantryDatabase.Collection("ingredient")
	ingredientCollection := "ingredient"
	pantryIngredient := database.NewCollection(ingredientCollection, pantryDatabase)

	//userCollection := "user"
	//pantryUser := database.CreateCollection(userCollection, pantryDatabase)

	//pantryResult, err := pantryIngredient.InsertOne(ctx, bson.D{
	//	{"name", "flour"},
	//})
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(pantryResult.InsertedID)
	userChoice := AppMenu()
	AppSelection(userChoice, pantryIngredient, ctx)
	ing := "Jam"

	database.InsertDataToCollection(pantryIngredient, ctx, ing)
}

func AppMenu() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Your Pantry Application")
	fmt.Println("-----------------------")
	fmt.Println("Please select an action")
	fmt.Println("1. Add Ingredient")
	fmt.Println("2. Remove Ingredient")
	fmt.Println("3. View Ingredients")
	fmt.Println("4. Search Recipes")
	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)
	return text
}

func AppSelection(choice string, collection *mongo.Collection, ctx context.Context) {
	switch choice {
	case "1":
		fmt.Println("Please type the ingredient to add...")
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		database.InsertDataToCollection(collection, ctx, text)
		userChoice := AppMenu()
		AppSelection(userChoice, collection, ctx)
	case "2":
		fmt.Println("Please type the ingredient to remove...")
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		database.RemoveManyFromCollection(collection, ctx, text)
		userChoice := AppMenu()
		AppSelection(userChoice, collection, ctx)
	case "3":
		database.ListDocuments(collection, ctx)
		userChoice := AppMenu()
		AppSelection(userChoice, collection, ctx)
	}
}
