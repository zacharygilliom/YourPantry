package recipe

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

type Ingredient struct {
	Ingredient string
	Weight     float64
}

type Recipe struct {
	Title       string
	Ingredients []Ingredient
}

type Recipes struct {
	Recipes []Recipe
}

func GetRecipes(collection *mongo.Collection, userID interface{}) {
	var r AutoGenerated
	//ingredients := database.BuildStringFromIngredients(collection, userID)
	resp, err := http.Get("https://api.edamam.com/api/recipes/v2?type=public&q=chicken&app_id=d18be80e&app_key=3d1fe9b7d97a890ac8e4d47c3d54fa88")
	if err != nil {
		log.Fatal(err)
	} else if err == nil {
		fmt.Println("Request Sent Successfully")
	}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	var rs Recipes
	for _, hit := range r.Hits {
		var rec Recipe
		rec.Title = hit.Recipe.Label
		for _, f := range hit.Recipe.Ingredients {
			var Ing Ingredient
			Ing.Ingredient = f.Text
			Ing.Weight = f.Weight
			rec.Ingredients = append(rec.Ingredients, Ing)
		}
		rs.Recipes = append(rs.Recipes, rec)
	}
}
