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

type HandlerA struct {
	pUser       *mongo.Collection
	pIngredient *mongo.Collection
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

	Obj := new(HandlerA)
	Obj.pUser = pantryUser
	Obj.pIngredient = pantryIngredient

	//routers
	// Need to pass the userID as a paramater in the api that gets returned after the frontend calls the addUser endpoint.
	r := engine(Obj)
	r.Use(cors.Default())
	r.Use(gin.Logger())
	r.Run()
}

func engine(Obj *HandlerA) *gin.Engine {
	r := gin.New()
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))
	r.GET("/login/:email/:password", Obj.loginUser)

	private := r.Group("/user")
	private.Use(AuthRequired)
	{
		private.POST("/ingredients/add", Obj.addIngredient)
		private.POST("/ingredients/remove", Obj.removeIngredient)
		private.GET("/ingredients/list", Obj.listIngredients)
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

func (this *HandlerA) loginUser(c *gin.Context) {
	session := sessions.Default(c)
	collection := this.pUser
	email := c.Param("email")
	password := c.Param("password")
	var users []string
	users = database.GetUser(collection, email, password)
	if len(users) > 1 {
		c.JSON(200, gin.H{"message": "Multiple Users Retrieved",
			"data": users,
		})
		return
	}
	userHex, _ := primitive.ObjectIDFromHex(users[0])
	userID := userHex.Hex()
	session.Set(userkey, userID)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	c.JSON(200, gin.H{
		"message": "Successfully authenticated user",
	})
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
func (this *HandlerA) addIngredient(c *gin.Context) {
	session := sessions.Default(c)
	userHex := session.Get(userkey)
	userID, err := primitive.ObjectIDFromHex(userHex.(string))
	if err != nil {
		log.Fatal(err)
	}
	collection := this.pIngredient
	ingredient := c.Query("ingredient")
	//userHex := c.Param("userID")
	//fmt.Println(userHex)
	database.InsertDataToIngredients(collection, userID, ingredient)
	c.JSON(200, gin.H{
		"message": "Ingredient added",
	})
}

func (this *HandlerA) removeIngredient(c *gin.Context) {
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

func (this *HandlerA) listIngredients(c *gin.Context) {
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
