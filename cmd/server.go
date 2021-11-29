package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run()
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
