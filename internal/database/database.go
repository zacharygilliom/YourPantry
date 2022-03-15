package database

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/zacharygilliom/configs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id, omitempty"`
	Firstname string             `bson:"firstname, omitempty"`
	Lastname  string             `bson:"lastname, omitempty"`
	Email     string             `bson:"email"`
	Password  string             `bson:"password"`
}

type Ingredient struct {
	ID   primitive.ObjectID `bson:"_id, omitempty"`
	User primitive.ObjectID `bson:"user, omitempty"`
	Name string             `bson:"name, omitempty"`
}

type Conn struct {
	DB         *mongo.Database
	User       *mongo.Collection
	Ingredient *mongo.Collection
}

func Init(ctx context.Context) (*Conn, *mongo.Client, error) {
	user, password, databaseName := configs.GetMongoCreds()
	databaseURI := "mongodb+srv://" + user + ":" + password + "@production.tobvq.mongodb.net/" + databaseName + "?retryWrites=true&w=majority"
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(databaseURI))
	//defer client.Disconnect(ctx)
	db := client.Database("pantry")
	u := db.Collection("user")
	i := db.Collection("ingredient")
	d := Conn{db, u, i}
	return &d, client, err
}

func PingClient(ctx context.Context, client *mongo.Client) (string, error) {
	err := client.Ping(ctx, readpref.Primary())
	message := ""
	if err != nil {
		message = "client not connected"
	} else {
		message = "client connected"
	}
	return message, err
}

func (conn *Conn) GetUser(email string, password string) []string {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := conn.User.Find(ctx,
		bson.M{"email": email})
	var mongoUser User
	var emails []string
	if err != nil {
		log.Fatal(err)
	}
	for cursor.Next(ctx) {
		err := cursor.Decode(&mongoUser)
		if err != nil {
			log.Fatal(err)
		}
		idString := mongoUser.ID.Hex()
		if checkPasswordHash(password, mongoUser.Password) {
			emails = append(emails, idString)
		}
	}
	defer cancel()
	defer ctx.Done()
	return emails
}

func (conn *Conn) InsertDataToUsers(email, password, fname, lname string) interface{} {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		log.Fatal(err)
	}
	password = hashedPassword
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	defer ctx.Done()
	cursor, err := conn.User.Find(ctx, bson.M{"email": email})
	var mongoUser User
	var userCheck string
	if err != nil {
		log.Fatal(err)
	} else {
		for cursor.Next(ctx) {
			err := cursor.Decode(&mongoUser)
			if err != nil {
				log.Fatal(err)
			}
			userCheck = mongoUser.Email
		}
	}
	if userCheck == "" {
		data := bson.D{
			{"firstname", fname},
			{"lastname", lname},
			{"email", email},
			{"password", password},
		}
		result, err := conn.User.InsertOne(ctx, data)
		fmt.Println("User Added to Collection")
		if err != nil {
			log.Fatal(err)
		}
		return result.InsertedID
	} else {
		return mongoUser.ID
	}
}

func (conn *Conn) InsertDataToIngredients(userHex interface{}, data string) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	defer ctx.Done()
	ingredient := bson.D{
		{"user", userHex},
		{"name", data},
	}
	result, err := conn.Ingredient.InsertOne(ctx, ingredient)
	fmt.Printf("%v added to collection\n", data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result.InsertedID)
}

func (conn *Conn) RemoveManyFromIngredients(userHex interface{}, data string) int64 {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	defer ctx.Done()
	filter := bson.D{
		{"user", userHex},
		{"name", data},
	}
	result, err := conn.Ingredient.DeleteMany(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}
	if result.DeletedCount > 1 {
		fmt.Printf("%v instances of %v deleted\n", result.DeletedCount, data)
	} else if result.DeletedCount == 0 {
		fmt.Printf("%v instances of %v exist, no action taken\n", result.DeletedCount, data)
	} else {
		fmt.Printf("%v instance of %v deleted\n", result.DeletedCount, data)
	}
	return result.DeletedCount
}

func (conn *Conn) ListIngredients(userHex interface{}) []Ingredient {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	defer ctx.Done()
	cursor, err := conn.Ingredient.Find(ctx, bson.M{"user": userHex})
	var results []Ingredient
	if err != nil {
		log.Fatal(err)
	} else {
		for cursor.Next(ctx) {
			var result Ingredient
			err := cursor.Decode(&result)
			if err != nil {
				log.Fatal(err)
			} else {
				results = append(results, result)
			}
		}
	}
	return results
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func BuildStringFromIngredients(collection *mongo.Collection, userID interface{}) string {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	defer ctx.Done()
	cursor, err := collection.Find(ctx, bson.M{"user": userID})
	var ings strings.Builder
	if err != nil {
		log.Fatal(err)
	} else {
		for cursor.Next(ctx) {
			var result Ingredient
			err := cursor.Decode(&result)
			if err != nil {
				log.Fatal(err)
			} else {
				ings.WriteString(result.Name + ",")
			}
		}
	}
	ingsSlice := ings.String()[:len(ings.String())-1]
	return ingsSlice
}
