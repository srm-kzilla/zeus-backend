package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Speaker struct {
	Name       string `json:"name" bson:"name"`
	Email      string `json:"email" bson:"email"`
	GithubLink string `json:"githubLink" bson: "githubLink"`
	LinkedIn   string `json:"linkedIn" bson:"linkedIn"`
}

type User struct {
	Name        string `json:"name" bson:"name"`
	Email       string `json:"email" bson:"email"`
	PhoneNumber int    `json:"phoneNumber" bson:"phoneNumber"`
	Feedback    string `json:"feedback" bson:"feedback"`
}

type Event struct {
	Id          primitive.ObjectID `json:"id" bson:"_id"`
	Title       string             `json:"title" bson:"title"`
	Slug        string             `json:"slug" bson: "slug"`
	Description string             `json:"description" bson: "description"`
	StartDate   time.Time          `json:"startDate" bson:"startDate"`
	Speakers    []Speaker          `json:"speakers,omitempty" bson:"speakers,omitempty"`
	Users       []User             `json:"users,omitempty" bson:"users,omitempty"`
	EventCover  string             `json:"eventCover,omitempty" bson:"eventCover,omitempty"` // s3-url for event cover Image
	IsCompleted bool               `json:"isCompleted,omitempty" bson:"isCompleted,omitempty"`
}
