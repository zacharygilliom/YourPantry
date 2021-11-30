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

	r := gin.Default()
	r.GET("/:ingredient", addIngredient(userID, pantryIngredient))
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
