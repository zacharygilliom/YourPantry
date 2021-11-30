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
	r.GET("/:ingredient", addIngredient(c*gin.Context, pantryIngredient))
	r.Run()
}

func addIngredients(c *gin.Context, collection *mongo.Collection) {
	database.InsertDataToIngredients(collection, userID, "chicken")
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
