package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zacharygilliom/internal/database"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	//Connect to database
	fmt.Println("Connection Started...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	client, err := database.CreateConnection(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	fmt.Println("Connection Established")

	//Initialize our database and collections
	pantryDatabase := database.NewDatabase(client)
	pantryUser := database.NewCollection("user", pantryDatabase)
	pantryIngredient := database.NewCollection("ingredient", pantryDatabase)

	//Get user info and either add user or return valid user
	userData := database.GetUserInfo()
	userID := database.InsertDataToUsers(pantryUser, userData)

	//routers
	r := gin.Default()
	r.GET("/ingredients/add/:ingredient", addIngredient(userID, pantryIngredient))
	r.GET("/ingredients/remove/:ingredient", removeIngredient(userID, pantryIngredient))
	r.GET("/ingredients/list/", listIngredients(userID, pantryIngredient))
	r.Run()
}

func addIngredient(userID interface{}, collection *mongo.Collection) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		ingredient := c.Param("ingredient")
		database.InsertDataToIngredients(collection, userID, ingredient)
		c.JSON(200, gin.H{
			"message": "Ingredient added",
		})
	}
	return gin.HandlerFunc(fn)
}

func removeIngredient(userID interface{}, collection *mongo.Collection) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		ingredient := c.Param("ingredient")
		database.RemoveManyFromIngredients(collection, userID, ingredient)
		c.JSON(200, gin.H{
			"message": "Ingredient removed",
		})
	}
	return gin.HandlerFunc(fn)
}

func listIngredients(userID interface{}, collection *mongo.Collection) gin.HandlerFunc {
	//Need to add error catch in case resultsMap returns an empty string
	fn := func(c *gin.Context) {
		ingredList := database.ListDocuments(collection, userID)
		resultsMap := make(map[int]string)
		for i, ing := range ingredList {
			resultsMap[i] = ing.Name
		}
		c.JSON(200, resultsMap)
	}
	return gin.HandlerFunc(fn)
}

/*
func main() {
	fmt.Println("Connection Started...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	client, err := database.CreateConnection(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	fmt.Println("Connection Established")

	databaseName := "pantry"
	pantryDatabase := database.NewDatabase(databaseName, client)

	userCollection := "user"
	pantryUser := database.NewCollection(userCollection, pantryDatabase)

	ingredientCollection := "ingredient"
	pantryIngredient := database.NewCollection(ingredientCollection, pantryDatabase)

	userData := database.GetUserInfo()
	userID := database.InsertDataToUsers(pantryUser, userData)

	userChoice := cli.Menu()
	cli.Selection(userChoice, pantryIngredient, userID)
}
*/
