package controllers

import (
	"log"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	database "github.com/zacharygilliom/internal/database"
	"github.com/zacharygilliom/internal/models"
	"github.com/zacharygilliom/internal/recipe"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Connection is a middleman from the database connection type so that we can write the functions as methods
type Connection struct {
	Conn *database.Conn
}
type login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var identityKey = "id"

// SignUpUser the user returns the newly created userID
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

//LoginUser logs the user in and sets the jwt token as the userid for the current session
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
	if userID == "" {
		return nil, jwt.ErrFailedAuthentication
	}
	authUser := struct {
		name string
	}{
		name: userID,
	}
	return authUser, nil
}

// AddIngredient takes the jwt session toek and then adds the user generated ingredient into that users data
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

// GetUserData extracts the user's data given it exists for use on the account profile page
func (conn *Connection) GetUserData(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	userHex := claims[identityKey]
	userID, err := primitive.ObjectIDFromHex(userHex.(string))
	if err != nil {
		log.Fatal(err)
	}
	var user models.User
	user = conn.Conn.GetUserData(userID)
	c.JSON(200, gin.H{
		"firstname": user.Firstname,
		"lastname":  user.Lastname,
		"email":     user.Email,
	})
}

// RemoveIngredient gets jwt token from session and removes the chosen ingredient if it exists from the users pantry
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

// ListIngredients returns the list of user's ingredients
func (conn *Connection) ListIngredients(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	userHex := claims[identityKey]
	userID, err := primitive.ObjectIDFromHex(userHex.(string))
	if err != nil {
		log.Fatal(err)
	}
	ingredientCollectionList := conn.Conn.ListIngredients(userID)
	var ingredList []string
	for _, ing := range ingredientCollectionList {
		ingredList = append(ingredList, ing.Name)
	}
	c.JSON(200, gin.H{
		"ingredients": ingredList,
		"size":        len(ingredList),
	})
}

// SearchRecipes takes the ingredients from the users database and searches for ingredients that are compatible
func (conn *Connection) SearchRecipes(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	userHex := claims[identityKey]
	userID, err := primitive.ObjectIDFromHex(userHex.(string))
	if err != nil {
		log.Fatal(err)
	}
	recipes := recipe.GetRecipes()
	ingredientCollectionList := conn.Conn.ListIngredients(userID)
	//var ingredList []string
	/*
		for _, ing := range ingredientCollectionList {
			fmt.Println(ing.Name)
			ingredList = append(ingredList, ing.Name)
		}
		for _, recipe := range recipes.Recipes {
			for _, r := range recipe.Ingredients {
				fmt.Println(r.Ingredient)
				fmt.Println(r.Weight)
			}
		}
	*/
	c.JSON(200, gin.H{
		"ingredients": ingredientCollectionList,
		"recipes":     recipes,
	})
}
