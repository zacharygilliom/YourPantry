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
	Id string
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
	db, err := database.Init(ctx)
	if err != nil {
		log.Fatal(err)
	}

	//Create instance of our database connection and run our engine
	var conn *controllers.Connection
	conn.DB = db
	r := engine(conn.DB)
	r.Run()
}

func engine(db *database.DB) *gin.Engine {
	r := gin.Default()
	//set new cookie store and new session
	r.Use(cors.Default())
	r.Use(gin.Logger())

	authMiddleware, err := jwt.Init()
	if err != nil {
		log.Fatal(err)
	}

	//endpoints to login or create account
	r.POST("/login", authMiddleware.LoginHandler)
	r.POST("/sign-up", db.SignUpUser)

	private := r.Group("/user")
	{
		private.POST("/ingredients/add", db.AddIngredient)
		private.POST("/ingredients/remove", db.RemoveIngredient)
		private.GET("/ingredients/list", db.ListIngredients)
	}
	return r
}
