package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Speaker struct {
	Name       string `json:"name" bson:"name"`
	Email      string `json:"email" bson:"email"`
	GithubLink string `json:"githubLink" bson:"githubLink"`
	LinkedIn   string `json:"linkedIn" bson:"linkedIn"`
	EventSlug  string `validate:"required" json:"slug" bson:"slug"`
}

type User struct {
	Name        string `validate:"required" json:"name" bson:"name"`
	Email       string `validate:"required,email" json:"email" bson:"email"`
	PhoneNumber int    `validate:"required,min=10" json:"phoneNumber" bson:"phoneNumber"`
	Feedback    string `json:"feedback" bson:"feedback"`
	EventSlug   string `validate:"required" json:"slug" bson:"slug"`
}

type Event struct {
	ID			primitive.ObjectID `json:"_id" bson:"_id"`
	Title       string    `validate:"required" json:"title" bson:"title"`
	Slug        string    `validate:"required" json:"slug" bson:"slug"`
	Description string    `validate:"required" json:"description" bson:"description"`
	StartDate   time.Time `json:"startDate" bson:"startDate"`
	// Speakers    []Speaker `json:"speakers" bson:"speakers"`
	// Users       []User    `json:"users" bson:"users"`
	EventCover  string `validate:"required" json:"eventCover" bson:"eventCover"` // s3-url for event cover Image
	IsCompleted bool   `json:"isCompleted" bson:"isCompleted"`
}
