package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Speaker struct {
	ID 		   primitive.ObjectID `json:"_id" bson:"_id"`
	Name       string `json:"name" bson:"name"`
	Email      string `json:"email" bson:"email"`
	GithubLink string `json:"githubLink" bson:"githubLink"`
	LinkedIn   string `json:"linkedIn" bson:"linkedIn"`
	EventSlug  string `validate:"required" json:"slug" bson:"slug"`
}

type User struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	Name        string `validate:"required" json:"name" bson:"name"`
	Email       string `validate:"required,email" json:"email" bson:"email"`
	PhoneNumber int    `validate:"required,min=10" json:"phoneNumber" bson:"phoneNumber"`
	Feedback    string `json:"feedback" bson:"feedback"`
	EventSlug   string `validate:"required" json:"slug" bson:"slug"`
}

type Event struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	Title       string             `validate:"required" json:"title" bson:"title"`
	Slug        string             `validate:"required" json:"slug" bson:"slug"`
	Description string             `validate:"required" json:"description" bson:"description"`
	Icons	    []string           `json:"icons" bson:"icons"`
	StartDate   string             `json:"startDate" bson:"startDate"`
	EventCover  string `validate:"required" json:"eventCover" bson:"eventCover"` // s3-url for event cover Image
	IsCompleted bool   `json:"isCompleted" bson:"isCompleted"`
}
