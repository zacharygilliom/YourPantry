package recipe

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/zacharygilliom/internal/database"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetRecipes(collection *mongo.Collection, userID interface{}) {
	ingredients := database.BuildStringFromIngredients(collection, userID)
	resp, err := http.Get("https://api.edamam.com/search/?q=" + ingredients + "app_id=d18be80e&app_key=3d1fe9b7d97a890ac8e4d47c3d54fa88&from=0&to=1")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	//fmt.Println(resp.Body)
	if err := json.NewDecoder(resp.Body); err != nil {
		log.Fatal(err)
	}
	/*
		if err := json.NewDecoder(resp.Body).Decode(); err != nil {
			log.Fatal(err)
		}

		for _, rec := range r.List {
			fmt.Println(rec.Title)
		}
	*/
}
