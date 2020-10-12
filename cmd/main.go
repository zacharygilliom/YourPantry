package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/zacharygilliom/internal/database"
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
	client, err := database.CreateConnection()
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	databaseName := "pantry"
	pantryDatabase := database.NewDatabase(databaseName, client)

	ingredientCollection := "ingredient"
	pantryIngredient := database.NewCollection(ingredientCollection, pantryDatabase)

	userChoice := AppMenu()
	AppSelection(userChoice, pantryIngredient)
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
	fmt.Println("5. Close application")
	fmt.Println("-----------------------")
	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)
	return text
}

func AppSelection(choice string, collection *mongo.Collection) {
	switch choice {
	case "1":
		fmt.Println("Please type the ingredient to add...")
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		database.InsertDataToCollection(collection, text)
		userChoice := AppMenu()
		AppSelection(userChoice, collection)
	case "2":
		fmt.Println("Please type the ingredient to remove...")
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		database.RemoveManyFromCollection(collection, text)
		userChoice := AppMenu()
		AppSelection(userChoice, collection)
	case "3":
		database.ListDocuments(collection)
		userChoice := AppMenu()
		AppSelection(userChoice, collection)
	case "4":
		SearchIngredients(collection)
		userChoice := AppMenu()
		AppSelection(userChoice, collection)
	case "5":
		fmt.Println("Application Closed")
	}
}

func SearchIngredients(collection *mongo.Collection) {
	ingred := database.BuildIngredientString(collection)
	ingredString := ingred.String()
	fmt.Println(ingredString)
	resp, err := http.Get("https://api.spoonacular.com/recipes/complexSearch?apiKey=58bbec758ee847f7b331410b02c7252d&includeIngredients=" + ingredString + "&number=10")
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
