package recipe

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
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

func GetRecipes() Recipes {
	var r AutoGenerated
	// Below code is for API call.  to save on API calls, a single api call was made and saved into a json file for testing.
	/*
		resp, err := http.Get("https://api.edamam.com/api/recipes/v2?type=public&q=" + keyword + "&app_id=d18be80e&app_key=3d1fe9b7d97a890ac8e4d47c3d54fa88")
		if err != nil {
			log.Fatal(err)
		} else if err == nil {
			fmt.Println("Request Sent Successfully")
		}
		if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
	*/
	// Below is the pre loaded json file for testing
	jsonFile, err := os.Open("/home/zach/programming/golang/YourPantry/internal/recipe/data.json")
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &r)

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
	return rs
}
