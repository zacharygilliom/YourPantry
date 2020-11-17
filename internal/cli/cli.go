package cli

import (
	"fmt"

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

func Menu() string {
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

func getUserInput() string {
	var text string
	fmt.Scanln(&text)
	return text
}

func Selection(choice string, collection *mongo.Collection, userID interface{}) {
	switch choice {
	case "1":
		fmt.Println("Please type the ingredient to add...")
		text := getUserInput()
		database.InsertDataToIngredients(collection, userID, text)
		userChoice := Menu()
		Selection(userChoice, collection, userID)
	case "2":
		fmt.Println("Please type the ingredient to remove...")
		text := getUserInput()
		database.RemoveManyFromIngredients(collection, userID, text)
		userChoice := Menu()
		Selection(userChoice, collection, userID)
	case "3":
		database.ListDocuments(collection, userID)
		userChoice := Menu()
		Selection(userChoice, collection, userID)
	case "4":
		SearchIngredients(collection, userID)
		userChoice := Menu()
		Selection(userChoice, collection, userID)
	case "5":
		fmt.Println("Application Closed")
	}
}

/*
func SearchIngredients(collection *mongo.Collection, userID interface{}) {
	ingred := database.BuildIngredientString(collection, userID)
	resp, err := http.Get("https://api.spoonacular.com/recipes/complexSearch?apiKey=58bbec758ee847f7b331410b02c7252d&findByIngredients=" + ingred + "&number=10")
	if err != nil {
		log.Fatal(err)
	}
	var r RecipeList

	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		log.Fatal(err)
	}

	fmt.Println(r)

	for _, rec := range r.List {
		fmt.Println(rec.Title)
	}
}
*/
