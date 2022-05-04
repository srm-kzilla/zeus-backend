package userModel

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// type User struct {
// 	ID          primitive.ObjectID `json:"_id" bson:"_id"`
// 	Name        string `validate:"required" json:"name" bson:"name"`
// 	Email       string `validate:"required,email" json:"email" bson:"email"`
// 	PhoneNumber int    `validate:"required,min=10" json:"phoneNumber" bson:"phoneNumber"`
// 	Feedback    string `json:"feedback" bson:"feedback"`
// 	EventSlugs   []string `validate:"required" json:"eventSlugs" bson:"eventSlugs"`
// }

type User struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	Name        string `validate:"required" json:"name" bson:"name"`
	Email       string `validate:"required,email" json:"email" bson:"email"`
	PhoneNumber int    `validate:"required,min=10" json:"phoneNumber" bson:"phoneNumber"`
	Feedback    string `json:"feedback" bson:"feedback"`
	EventSlugs   []string `json:"events" bson:"events"`
}

type RegisterUser struct {
	User	User	`json:"user" bson:"user"`
	EventSlug	string	`json:"eventSlug" bson:"eventSlug"`
}