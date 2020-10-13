package database

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//type User struct {
//	ID        primitive.ObjectID `bson:"_id, omitempty"`
//	Firstname string             `bson:"firstname, omitempty"`
//	Lastname  string             `bson:"lastname, omitempty"`
//	Email     string             `bson:"email, omitempty"`
//}
//
//type IngredientList struct {
//	ID         primitive.ObjectID `bson:"_id, omitempty"`
//	User       primitive.ObjectID `bson:"user, omitempty"`
//	Ingredient []Ingredient       `bson:"ingredient"`
//}
//
type Ingredient struct {
	ID   primitive.ObjectID `bson:"_id, omitempty"`
	Name string             `bson:"name, omitempty"`
}

func CreateConnection() (*mongo.Client, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://admin:Branstark1@production.tobvq.mongodb.net/pantry?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}
	return client, err
}

func NewDatabase(databaseName string, client *mongo.Client) *mongo.Database {
	database := client.Database(databaseName)
	return database
}

func NewCollection(collectionName string, database *mongo.Database) *mongo.Collection {
	collection := database.Collection(collectionName)
	return collection
}

func InsertDataToCollection(collection *mongo.Collection, data string) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	defer ctx.Done()
	ingredient := bson.D{
		{"name", data},
	}
	result, err := collection.InsertOne(ctx, ingredient)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result.InsertedID)
}

func RemoveManyFromCollection(collection *mongo.Collection, data string) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	defer ctx.Done()
	filter := bson.D{
		{"name", data},
	}
	result, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}
	if result.DeletedCount > 1 {
		fmt.Printf("%v instances of %v deleted", result.DeletedCount, data)
	} else {
		fmt.Printf("%v instance of %v deleted", result.DeletedCount, data)
	}
	fmt.Println("")
	fmt.Println("")
}

func ListDocuments(collection *mongo.Collection) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	defer ctx.Done()
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	} else {
		for cursor.Next(ctx) {
			var result Ingredient
			err := cursor.Decode(&result)
			if err != nil {
				log.Fatal(err)
			} else {
				fmt.Println(result.Name)
			}
		}
	}
}

func BuildIngredientString(collection *mongo.Collection) string {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	defer ctx.Done()
	cursor, err := collection.Find(ctx, bson.M{})
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
				ings.WriteString(result.Name + "&")
			}
		}
	}
	ingsSlice := ings.String()[:len(ings.String())-1]
	return ingsSlice
}
