package recipe

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

func GetRecipes(collection *mongo.Collection, userID interface{}) {
	var rs Recipes
	//ingredients := database.BuildStringFromIngredients(collection, userID)
	resp, err := http.Get("https://api.edamam.com/api/recipes/v2?type=public&q=chicken&app_id=d18be80e&app_key=3d1fe9b7d97a890ac8e4d47c3d54fa88")
	if err != nil {
		log.Fatal(err)
	} else if err == nil {
		fmt.Println("Request Sent Successfully")
	}
	dec := json.NewDecoder(resp.Body)
	dec.Decode(&rs)
	fmt.Printf("%v %v %v %v", rs.From, rs.To, rs.Count, rs.Hits)
	//recipes, err := io.ReadAll(resp.Body)
	//fmt.Println(recipes)
	//fmt.Println(resp.Body)
	/*
		if err := json.NewDecoder(resp.Body); err != nil {
			fmt.Println(err)
			log.Fatal(err)
		} else {
			fmt.Println("New Decoder Successful")
		}
		if err := json.NewDecoder(resp.Body).Decode(&rs); err != nil {
			log.Fatal(err)
		} else {
			fmt.Println("Decoded into recipes struct successfully")
		}
		fmt.Println(rs.Hits)
	*/
	defer resp.Body.Close()

	/*
		for _, rec := range r.List {
			fmt.Println(rec.Title)
		}
	*/
}
