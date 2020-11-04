package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/zacharygilliom/internal/cli"
	"github.com/zacharygilliom/internal/database"
)

func main() {
	fmt.Println("Connection Started...")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := database.CreateConnection(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

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
