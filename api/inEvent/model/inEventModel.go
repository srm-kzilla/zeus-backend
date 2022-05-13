package inEventModel

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InEventData struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id"`
	UserID       string             `validate:"required" json:"userId" bson:"userId"`
	EventSlug    string             `validate:"required" json:"eventSlug" bson:"eventSlug"`
	FoodReceived bool               `validate:"required" json:"foodReceived" bson:"foodReceived"`
}
