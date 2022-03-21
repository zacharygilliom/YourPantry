package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserPOST struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type NewUserPOST struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

type User struct {
	ID        primitive.ObjectID `bson:"_id, omitempty"`
	Firstname string             `bson:"firstname, omitempty"`
	Lastname  string             `bson:"lastname, omitempty"`
	Email     string             `bson:"email"`
	Password  string             `bson:"password"`
}
