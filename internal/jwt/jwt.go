package jwt

import (
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/zacharygilliom/internal/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var identityKey = "id"

type User struct {
	Id string
}

type login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Init() (*jwt.GinJWTMiddleware, error) {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("secret"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					identityKey: v.Id,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &User{
				Id: claims[identityKey].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			var users []string
			user = database.GetUser(collection, login.Email, login.Password)
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
			return &User{
				Id: userID,
			}, nil
			return nil, jwt.ErrFailedAuthentication
		},
	})
	return authMiddleware, err
}
