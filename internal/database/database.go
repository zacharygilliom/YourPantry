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

/*
type CreatedUser struct {
	Firstname string
	Lastname  string
	Email     string
}
*/

type Ingredient struct {
	ID   primitive.ObjectID `bson:"_id, omitempty"`
	User primitive.ObjectID `bson:"user, omitempty"`
	Name string             `bson:"name, omitempty"`
}

func CreateConnection(ctx context.Context) (*mongo.Client, error) {
	user, password, database := configs.GetMongoCreds()
	databaseURI := "mongodb+srv://" + user + ":" + password + "@production.tobvq.mongodb.net/" + database + "?retryWrites=true&w=majority"
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(databaseURI))
	if err != nil {
		log.Fatal(err)
	}
	return client, err
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

func NewDatabase(client *mongo.Client) *mongo.Database {
	database := client.Database("pantry")
	return database
}

func NewCollection(collectionName string, database *mongo.Database) *mongo.Collection {
	collection := database.Collection(collectionName)
	return collection
}

func GetUserInfo() User {
	var NewUser User
	/*
		fmt.Println("Please enter the User's First Name")
		fmt.Scanf("%s", &NewUser.Firstname)
		fmt.Println("Please enter the User's Last Name")
		fmt.Scanf("%s", &NewUser.Lastname)
		fmt.Println("Please enter the User's Email")
		fmt.Scanf("%s", &NewUser.Email)
	*/
	NewUser.Firstname = "Zachary"
	NewUser.Lastname = "Gilliom"
	NewUser.Email = "zacharygilliom@gmail.com"
	NewUser.Password = "123"
	return NewUser
}

func GetUser(collection *mongo.Collection, email string, password string) []string {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	defer ctx.Done()
	cursor, err := collection.Find(ctx,
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
	return emails
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func InsertDataToUsers(collection *mongo.Collection, email, password, fname, lname string) interface{} {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		log.Fatal(err)
	}
	password = hashedPassword
	fmt.Println(email, password, fname, lname)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	defer ctx.Done()
	cursor, err := collection.Find(ctx, bson.M{"email": email})
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
		result, err := collection.InsertOne(ctx, data)
		fmt.Println("User Added to Collection")
		if err != nil {
			log.Fatal(err)
		}
		return result.InsertedID
	} else {
		return mongoUser.ID
	}
}

func InsertDataToIngredients(collection *mongo.Collection, userHex interface{}, data string) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	defer ctx.Done()
	//userID, err := primitive.ObjectIDFromHex(userHex)
	//fmt.Println(userID)
	//if err != nil {
	//	log.Fatal(err)
	//}
	ingredient := bson.D{
		{"user", userHex},
		{"name", data},
	}
	result, err := collection.InsertOne(ctx, ingredient)
	fmt.Printf("%v added to collection\n", data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result.InsertedID)
}

func RemoveManyFromIngredients(collection *mongo.Collection, userHex interface{}, data string) int64 {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	defer ctx.Done()
	//userID, err := primitive.ObjectIDFromHex(userHex)
	//if err != nil {
	//	log.Fatal(err)
	//}
	filter := bson.D{
		{"user", userHex},
		{"name", data},
	}
	result, err := collection.DeleteMany(ctx, filter)
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

func ListIngredients(collection *mongo.Collection, userHex interface{}) []Ingredient {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	defer ctx.Done()
	//fmt.Println(userHex)
	//userID, err := primitive.ObjectIDFromHex(userHex)
	//if err != nil {
	//	log.Fatal(err)
	//}
	cursor, err := collection.Find(ctx, bson.M{"user": userHex})
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
