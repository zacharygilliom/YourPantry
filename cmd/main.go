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
	client, ctx, err := database.CreateConnection()
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
	AppSelection(userChoice, pantryIngredient, ctx)
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
	fmt.Println("-----------------------")
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
	case "4":
		SearchIngredients()
	}
}

func SearchIngredients() {
	ingred := "eggs,milk,cheese"
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
	//body, err := ioutil.ReadAll(resp.Body)
	//fmt.Println(resp)
	//fmt.Println(string(body))
	//if err != nil {
	//	log.Fatal(err)
	//}
	//for _, rec := range body {
	//	recString := string(rec)
	//	fmt.Println(recString)
	//}
}
