package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	database "github.com/zacharygilliom/internal/database"
	"github.com/zacharygilliom/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Connection struct {
	DB *database.DB
}

func (conn *Connection) SignUpUser(c *gin.Context) {
	session := sessions.Default(c)
	newUser := models.NewUserPOST{}
	err := c.BindJSON(&newUser)
	if err != nil {
		c.AbortWithError(400, err)
	}
	var userID interface{}
	userID = conn.DB.InsertDataToUsers(newUser.Email, newUser.Password, newUser.Firstname, newUser.Lastname)
	session.Set("user", userID.(primitive.ObjectID).Hex())
	if err := session.Save(); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	c.JSON(200, gin.H{"message": "user added and authenticated",
		"data": 1})
}

func (conn *Connection) LoginUser(c *gin.Context) {
	//Get session and bind POSTED JSON data to posteduser struct
	session := sessions.Default(c)
	postedUser := models.UserPOST{}
	err := c.BindJSON(&postedUser)
	if err != nil {
		c.AbortWithError(400, err)
	}
	//send user info to getuser func to retrieve users based on email and password
	var users []string
	users = conn.DB.GetUser(postedUser.Email, postedUser.Password)
	if len(users) > 1 {
		c.JSON(200, gin.H{"message": "Multiple Users Retrieved",
			"data": 0})
		return
	} else if len(users) == 0 {
		c.JSON(200, gin.H{"message": "No users retrieved",
			"data": 0})
		return
	}
	//set userID as the session "user" variable
	userHex, _ := primitive.ObjectIDFromHex(users[0])
	userID := userHex.Hex()
	session.Set("user", userID)
	c.Set("user", userID)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	c.JSON(200, gin.H{"message": "authenticated user",
		"data": 1})
}

func (conn *Connection) AddIngredient(c *gin.Context) {
	//get session and get user id variable
	session := sessions.Default(c)
	userHex := session.Get("user")
	userIngredient := struct {
		Ingredient string `json:"ingredient"`
	}{}
	//userIng := userIngredient{}
	userID, err := primitive.ObjectIDFromHex(userHex.(string))
	if err != nil {
		log.Fatal(err)
	}
	// set POSTED data to new ingredient
	err = c.BindJSON(&userIngredient)
	if err != nil {
		log.Fatal(err)
	}
	//pass new ingredient to database to add it based on the user in the session
	//collection := conn.db.Ingredient
	conn.DB.InsertDataToIngredients(userID, userIngredient.Ingredient)
	c.JSON(200, gin.H{
		"message": "Ingredient added",
	})
}

func (conn *Connection) RemoveIngredient(c *gin.Context) {
	session := sessions.Default(c)
	userHex := session.Get("user")
	userID, err := primitive.ObjectIDFromHex(userHex.(string))
	if err != nil {
		log.Fatal(err)
	}
	//collection := conn.db.Ingredient
	ingredient := c.Query("ingredient")
	ingsRemoved := conn.DB.RemoveManyFromIngredients(userID, ingredient)
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

func (conn *Connection) ListIngredients(c *gin.Context) {
	session := sessions.Default(c)
	userHex := session.Get("user")
	//userHex, _ := c.Get("user")
	userID, err := primitive.ObjectIDFromHex(userHex.(string))
	if err != nil {
		log.Fatal(err)
	}
	//collection := conn.db.Ingredient
	//userHex := c.Param("userID")
	ingredientCollectionList := conn.DB.ListIngredients(userID)
	resultsMap := make(map[string][]string)
	var ingredList []string
	for _, ing := range ingredientCollectionList {
		ingredList = append(ingredList, ing.Name)
	}
	resultsMap["ingredients"] = ingredList
	/*
		if len(resultsMap) == 0 {
			res, err := database.PingClient(ctx, client)
			if err != nil {
				c.JSON(200, gin.H{
					"message": res,
				})
			}
		}
	*/
	c.JSON(200, resultsMap)
}
