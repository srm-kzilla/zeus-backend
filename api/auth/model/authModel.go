package authModel

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID 	     	primitive.ObjectID 	`json:"_id" bson:"_id"`
	Email 		string 				`validate:"required" json:"email" bson:"email"`
	Password 	string				`validate:"required" json:"password" bson:"password"`
	CreatedAt	time.Time			`json:"createdAt" bson:"createdAt"`
}