package main

import (
	"context"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/zacharygilliom/internal/controllers"
	"github.com/zacharygilliom/internal/database"
	"github.com/zacharygilliom/internal/jwt"
)

const (
	userkey = "user"
)

type User struct {
	Username string
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
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	//Initialize our database and collections
	var dat *database.Conn = new(database.Conn)
	dat, client, err := database.Init(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	var conn *controllers.Connection = new(controllers.Connection)
	conn.Conn = dat
	//Create instance of our database connection and run our engine
	r := engine(conn)
	r.Run()
}

func engine(conn *controllers.Connection) *gin.Engine {
	r := gin.Default()
	//set new cookie store and new session
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:8000"}
	config.AllowHeaders = []string{"Content-Type, Origin, Authorization, Access-Control-Allow-Headers"}
	config.AllowCredentials = true
	r.Use(cors.New(config))
	r.Use(gin.Logger())

	authMiddleware, err := jwt.Init(conn)
	if err != nil {
		log.Fatal(err)
	}

	r.POST("/sign-up", conn.SignUpUser)

	//endpoints to login account and authenticate the user
	r.POST("/login", authMiddleware.LoginHandler)
	r.POST("/logout", authMiddleware.LogoutHandler)
	private := r.Group("/user")
	private.Use(authMiddleware.MiddlewareFunc())
	{
		private.POST("/ingredients/add", conn.AddIngredient)
		private.POST("/ingredients/remove", conn.RemoveIngredient)
		private.GET("/ingredients/list", conn.ListIngredients)
		private.GET("/recipes/search", conn.SearchRecipes)
		private.GET("/data", conn.GetUserData)
	}
	return r
}
