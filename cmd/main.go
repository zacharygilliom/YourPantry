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
	Email    string `json: "email"`
	Password string `json: "password"`
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

	db := new(Connection)
	db.pUser = pantryUser
	db.pIngredient = pantryIngredient

	r := engine(db)
	r.Run()
}

func engine(db *Connection) *gin.Engine {
	r := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))
	r.Use(cors.Default())
	//r.Use(gin.Logger())
	r.POST("/login", db.loginUser)

	private := r.Group("/user")
	private.Use(AuthRequired)
	{
		private.POST("/ingredients/add", db.addIngredient)
		private.POST("/ingredients/remove", db.removeIngredient)
		private.GET("/ingredients/list", db.listIngredients)
	}

	//r.POST("/:userID/ingredients/add", addIngredient(pantryIngredient))
	//r.POST("/:userID/ingredients/remove", removeIngredient(pantryIngredient))
	//r.POST("/user/add", addUser(pantryUser))
	//r.GET("/:userID/ingredients/list", listIngredients(pantryIngredient, ctx, client))
	return r
}

func AuthRequired(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userkey)
	if user == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	c.Next()
}

func (this *Connection) loginUser(c *gin.Context) {
	session := sessions.Default(c)
	postedUser := &userPOST{}
	c.Bind(&postedUser)
	collection := this.pUser
	var users []string
	users = database.GetUser(collection, postedUser.Email, postedUser.Password)
	if len(users) > 1 {
		c.JSON(200, gin.H{"message": "Multiple Users Retrieved"})
		return
	} else if len(users) == 0 {
		c.JSON(200, gin.H{"message": "No users retrieved"})
		return
	}
	userHex, _ := primitive.ObjectIDFromHex(users[0])
	userID := userHex.Hex()
	session.Set(userkey, userID)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	c.JSON(200, gin.H{"message": "authenticated user"})

}

/*
func addUser(collection *mongo.Collection) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		newUser := database.User{
			Firstname: c.Query("firstname"),
			Lastname:  c.Query("lastname"),
			Email:     c.Query("email"),
			Password:  c.Query("password"),
		}
		newUserID := database.InsertDataToUsers(collection, newUser)
		c.JSON(200, gin.H{
			"message": "New User Created",
			"data":    newUserID,
		})
	}
	return gin.HandlerFunc(fn)
}
*/
func (this *Connection) addIngredient(c *gin.Context) {
	session := sessions.Default(c)
	userHex := session.Get(userkey)
	userID, err := primitive.ObjectIDFromHex(userHex.(string))
	if err != nil {
		log.Fatal(err)
	}
	collection := this.pIngredient
	ingredient := c.Query("ingredient")
	database.InsertDataToIngredients(collection, userID, ingredient)
	c.JSON(200, gin.H{
		"message": "Ingredient added",
	})
}

func (this *Connection) removeIngredient(c *gin.Context) {
	session := sessions.Default(c)
	userHex := session.Get(userkey)
	userID, err := primitive.ObjectIDFromHex(userHex.(string))
	if err != nil {
		log.Fatal(err)
	}
	collection := this.pIngredient
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

func (this *Connection) listIngredients(c *gin.Context) {
	session := sessions.Default(c)
	userHex := session.Get(userkey)
	userID, err := primitive.ObjectIDFromHex(userHex.(string))
	if err != nil {
		log.Fatal(err)
	}
	collection := this.pIngredient
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
