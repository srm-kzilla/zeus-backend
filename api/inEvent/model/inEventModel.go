package InEventModel

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InEventData struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id"`
	UserID       primitive.ObjectID `validate:"required" json:"userId" bson:"userId"`
	EventSlug    string             `validate:"required" json:"eventSlug" bson:"eventSlug"`
	FoodReceived bool               `validate:"required" default:"false" json:"foodReceived" bson:"foodReceived"`
}

type AttendanceQuery struct {
	UserID string `validate:"required" json:"userId" bson:"userId"`
	Slug  string `validate:"required" json:"eventSlug" bson:"eventSlug"`
	Action string `validate:"required" json:"action" bson:"action"`
}
