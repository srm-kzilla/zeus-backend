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
	Name        string             `validate:"required" json:"name" bson:"name"`
	Email       string             `validate:"required,email" json:"email" bson:"email"`
	PhoneNumber int                `validate:"required,min=1000000000,max=9999999999" json:"phoneNumber" bson:"phoneNumber"`
	RegNumber   string             `validate:"required" json:"regNumber" bson:"regNumber"`
	EventSlugs  []string           `json:"events" bson:"events"`
}

type RsvpUsers struct {
	UserId       string `json:"userId" bson:"userId"`
	FoodReceived bool   `json:"foodReceived" bson:"foodReceived"`
	CheckedIn    bool   `json:"checkedIn" bson:"checkedIn"`
}

type RegisterUserReq struct {
	User      User   `validate:"required" "json:"user" bson:"user"`
	EventSlug string `validate:"required" json:"eventSlug" bson:"eventSlug"`
}

type RsvpUserReq struct {
	UserId    string `validate:"required" json:"userId" bson:"userId"`
	EventSlug string `validate:"required" json:"eventSlug" bson:"eventSlug"`
}
