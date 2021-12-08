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
	// Need to pass the userID as a paramater in the api that gets returned after the frontend calls the addUser endpoint.
	r := gin.Default()
	r.POST("/ingredients/add/:ingredient", addIngredient(userID, pantryIngredient))
	r.POST("/ingredients/remove/:ingredient", removeIngredient(userID, pantryIngredient))
	r.POST("/user/add/:firstName/:lastName/:email", addUser(pantryUser))
	//r.GET("/user/list/:firstName/:lastName/:email", listUsers())
	r.GET("/ingredients/list", listIngredients(userID, pantryIngredient, ctx, client))
	r.Run()
}

func addUser(collection *mongo.Collection) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		newUser := database.User{
			Firstname: c.Param("firstName"),
			Lastname:  c.Param("lastName"),
			Email:     c.Param("email"),
		}
		newUserID := database.InsertDataToUsers(collection, newUser)
		c.JSON(200, gin.H{
			"message": "New User Created",
			"data":    newUserID,
		})
	}
	return gin.HandlerFunc(fn)
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
		ingsRemoved := database.RemoveManyFromIngredients(collection, userID, ingredient)
		data := map[int64]string{ingsRemoved: ingredient}
		if ingsRemoved > 0 {
			c.JSON(200, gin.H{
				"message": "Ingredient removed",
				"data":    data,
			})
		} else {
			c.JSON(200, gin.H{
				"message": "Ingredient is not in your pantry. Did you misspell it?",
				"data":    data,
			})
		}
	}
	return gin.HandlerFunc(fn)
}

func listIngredients(userID interface{}, collection *mongo.Collection, ctx context.Context, client *mongo.Client) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		ingredList := database.ListDocuments(collection, userID)
		resultsMap := make(map[int]string)
		for i, ing := range ingredList {
			resultsMap[i] = ing.Name
		}
		if len(resultsMap) == 0 {
			res, err := database.PingClient(ctx, client)
			if err != nil {
				c.JSON(200, gin.H{
					"message": res,
				})
			}
		}
		c.JSON(200, resultsMap)
	}
	return gin.HandlerFunc(fn)
}
