package recipe

import (
	"log"
	"net/http"
)

type Recipe struct {
	Name string `json:"name"`
}

func getRecipes(ingredients string) {
	resp, err := http.Get("https://api.edamam.com/search/?q=" + ingredients + "app_id=d18be80e&app_key=3d1fe9b7d97a890ac8e4d47c3d54fa88")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

}
