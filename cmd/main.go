package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/zacharygilliom/internal/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type RecipeList struct {
	List []Recipe `json:"results"`
}

type Recipe struct {
	ID                  int    `json:"id"`
	UsedIngredientCount int    `json:"usedIngredientCount"`
	Title               string `json:"title"`
}

func main() {
	fmt.Println("Connection Started...")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := database.CreateConnection(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	databaseName := "pantry"
	pantryDatabase := database.NewDatabase(databaseName, client)

	userCollection := "user"
	pantryUser := database.NewCollection(userCollection, pantryDatabase)

	ingredientCollection := "ingredient"
	pantryIngredient := database.NewCollection(ingredientCollection, pantryDatabase)

	userData := getUserInfo()
	userID := database.InsertDataToUsers(pantryUser, userData)

	userChoice := AppMenu()
	AppSelection(userChoice, pantryIngredient, userID)
}

func AppMenu() string {
	fmt.Println("Your Pantry Application")
	fmt.Println("-----------------------")
	fmt.Println("Please select an action")
	fmt.Println("1. Add Ingredient")
	fmt.Println("2. Remove Ingredient")
	fmt.Println("3. View Ingredients")
	fmt.Println("4. Search Recipes")
	fmt.Println("5. Close application")
	fmt.Println("-----------------------")
	text := getUserInput()
	return text
}

func AppSelection(choice string, collection *mongo.Collection, userID interface{}) {
	switch choice {
	case "1":
		fmt.Println("Please type the ingredient to add...")
		text := getUserInput()
		database.InsertDataToIngredients(collection, userID, text)
		userChoice := AppMenu()
		AppSelection(userChoice, collection, userID)
	case "2":
		fmt.Println("Please type the ingredient to remove...")
		text := getUserInput()
		database.RemoveManyFromIngredients(collection, userID, text)
		userChoice := AppMenu()
		AppSelection(userChoice, collection, userID)
	case "3":
		database.ListDocuments(collection, userID)
		userChoice := AppMenu()
		AppSelection(userChoice, collection, userID)
	case "4":
		SearchIngredients(collection, userID)
		userChoice := AppMenu()
		AppSelection(userChoice, collection, userID)
	case "5":
		fmt.Println("Application Closed")
	}
}

func getUserInput() string {
	var text string
	fmt.Scanf("%s", &text)
	return text
}

func getUserInfo() bson.D {
	var firstName string
	var lastName string
	var email string
	fmt.Println("Please enter the User's First Name")
	fmt.Scanf("%s", &firstName)
	fmt.Println("Please enter the User's Last Name")
	fmt.Scanf("%s", &lastName)
	fmt.Println("Please enter the User's Email")
	fmt.Scanf("%s", &email)
	data := bson.D{
		{"firstname", firstName},
		{"lastname", lastName},
		{"email", email},
	}
	return data
}

func SearchIngredients(collection *mongo.Collection, userID interface{}) {
	ingred := database.BuildIngredientString(collection, userID)
	resp, err := http.Get("https://api.spoonacular.com/recipes/complexSearch?apiKey=58bbec758ee847f7b331410b02c7252d&includeIngredients=" + ingred + "&number=10")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var r RecipeList

	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		log.Fatal(err)
	}

	fmt.Println(r)

	for _, rec := range r.List {
		fmt.Println(rec.Title)
	}
}
