package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/zacharygilliom/internal/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	userkey = "user"
)

type Connection struct {
	pUser       *mongo.Collection
	pIngredient *mongo.Collection
}

type userPOST struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type newUserPOST struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

type userIngredient struct {
	Ingredient string `json:"ingredient"`
}

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

	//Create instance of our database connection and run our engine
	db := new(Connection)
	db.pUser = pantryUser
	db.pIngredient = pantryIngredient

	r := engine(db)
	r.Run()
}

func engine(db *Connection) *gin.Engine {
	r := gin.Default()
	//set new cookie store and new session
	store := cookie.NewStore([]byte("secret"))
	store.Options(sessions.Options{MaxAge: 60 * 60 * 24,
		Path: "/"})
	r.Use(sessions.Sessions("mysession", store))
	r.Use(cors.Default())
	r.Use(gin.Logger())

	//endpoints to login or create account
	r.POST("/login", db.loginUser)
	r.POST("/sign-up", db.signUpUser)

	private := r.Group("/user")
	private.Use(AuthRequired())
	{
		private.POST("/ingredients/add", db.addIngredient)
		private.POST("/ingredients/remove", db.removeIngredient)
		private.GET("/ingredients/list", db.listIngredients)
	}
	return r
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")
		if user == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		c.Next()
	}
}

func (db *Connection) signUpUser(c *gin.Context) {
	session := sessions.Default(c)
	newUser := newUserPOST{}
	err := c.BindJSON(&newUser)
	if err != nil {
		c.AbortWithError(400, err)
	}
	collection := db.pUser
	var userID interface{}
	userID = database.InsertDataToUsers(collection, newUser.Email, newUser.Password, newUser.Firstname, newUser.Lastname)
	session.Set("user", userID.(primitive.ObjectID).Hex())
	if err := session.Save(); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	c.JSON(200, gin.H{"message": "user added and authenticated",
		"data": 1})
}

func (db *Connection) loginUser(c *gin.Context) {
	//Get session and bind POSTED JSON data to posteduser struct
	session := sessions.Default(c)
	postedUser := userPOST{}
	err := c.BindJSON(&postedUser)
	if err != nil {
		c.AbortWithError(400, err)
	}
	//send user info to getuser func to retrieve users based on email and password
	collection := db.pUser
	var users []string
	users = database.GetUser(collection, postedUser.Email, postedUser.Password)
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

func (db *Connection) addIngredient(c *gin.Context) {
	//get session and get user id variable
	session := sessions.Default(c)
	userHex := session.Get("user")
	userIng := userIngredient{}
	userID, err := primitive.ObjectIDFromHex(userHex.(string))
	if err != nil {
		log.Fatal(err)
	}
	// set POSTED data to new ingredient
	err = c.BindJSON(&userIng)
	if err != nil {
		log.Fatal(err)
	}
	//pass new ingredient to database to add it based on the user in the session
	collection := db.pIngredient
	database.InsertDataToIngredients(collection, userID, userIng.Ingredient)
	c.JSON(200, gin.H{
		"message": "Ingredient added",
	})
}

func (db *Connection) removeIngredient(c *gin.Context) {
	session := sessions.Default(c)
	userHex := session.Get("user")
	userID, err := primitive.ObjectIDFromHex(userHex.(string))
	if err != nil {
		log.Fatal(err)
	}
	collection := db.pIngredient
	ingredient := c.Query("ingredient")
	//userHex := c.Param("userID")
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

func (db *Connection) listIngredients(c *gin.Context) {
	session := sessions.Default(c)
	userHex := session.Get("user")
	//userHex, _ := c.Get("user")
	userID, err := primitive.ObjectIDFromHex(userHex.(string))
	if err != nil {
		log.Fatal(err)
	}
	collection := db.pIngredient
	//userHex := c.Param("userID")
	ingredientCollectionList := database.ListIngredients(collection, userID)
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
