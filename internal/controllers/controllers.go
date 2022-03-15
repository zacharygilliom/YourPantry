package controllers

import (
	"log"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	database "github.com/zacharygilliom/internal/database"
	"github.com/zacharygilliom/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Connection struct {
	Conn *database.Conn
}
type login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type User struct {
	Username string
}

var identityKey = "id"

func (conn *Connection) SignUpUser(c *gin.Context) {
	newUser := models.NewUserPOST{}
	err := c.BindJSON(&newUser)
	if err != nil {
		c.AbortWithError(400, err)
	}
	var userID interface{}
	userID = conn.Conn.InsertDataToUsers(newUser.Email, newUser.Password, newUser.Firstname, newUser.Lastname)
	c.JSON(200, gin.H{"message": "user added and authenticated",
		"data": userID})
}

func (conn *Connection) LoginUser(c *gin.Context) (interface{}, error) {
	var loginVals login
	if err := c.ShouldBind(&loginVals); err != nil {
		return "", jwt.ErrMissingLoginValues
	}
	var users []string
	users = conn.Conn.GetUser(loginVals.Email, loginVals.Password)
	if len(users) > 1 {
		c.JSON(200, gin.H{"message": "Multiple Users Retrieved",
			"data": 0})
		return "", jwt.ErrMissingLoginValues
	} else if len(users) == 0 {
		c.JSON(200, gin.H{"message": "No users retrieved",
			"data": 0})
		return "", jwt.ErrMissingLoginValues
	}
	//set userID as the session "user" variable
	userHex, _ := primitive.ObjectIDFromHex(users[0])
	userID := userHex.Hex()
	authUser := User{}
	authUser.Username = userID
	if userID == "" {
		return nil, jwt.ErrFailedAuthentication
	}
	return authUser, nil
}

func (conn *Connection) AddIngredient(c *gin.Context) {
	//get session and get user id variable
	userIngredient := struct {
		Ingredient string `json:"ingredient"`
	}{}
	claims := jwt.ExtractClaims(c)
	userHex := claims[identityKey]
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
	conn.Conn.InsertDataToIngredients(userID, userIngredient.Ingredient)
	c.JSON(200, gin.H{
		"message": "Ingredient added",
	})
}

func (conn *Connection) RemoveIngredient(c *gin.Context) {
	userIngredient := struct {
		Ingredient string `json:"ingredient"`
	}{}
	claims := jwt.ExtractClaims(c)
	userHex := claims[identityKey]
	userID, err := primitive.ObjectIDFromHex(userHex.(string))
	if err != nil {
		log.Fatal(err)
	}
	err = c.BindJSON(&userIngredient)
	if err != nil {
		log.Fatal(err)
	}
	ingsRemoved := conn.Conn.RemoveManyFromIngredients(userID, userIngredient.Ingredient)
	if ingsRemoved > 0 {
		c.JSON(200, gin.H{
			"message": "Ingredient removed",
			"data":    ingsRemoved,
		})
	} else {
		c.JSON(200, gin.H{
			"message": "Ingredient is not in your pantry. Did you misspell it?",
			"data":    ingsRemoved,
		})
	}
}

func (conn *Connection) ListIngredients(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	userID := claims[identityKey]
	ingredientCollectionList := conn.Conn.ListIngredients(userID)
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
