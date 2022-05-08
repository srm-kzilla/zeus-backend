package eventModel

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Speaker struct {
	ID 	     primitive.ObjectID 	`json:"_id" bson:"_id"`
	Name       string				`json:"name" bson:"name"`
	Email      string 			`json:"email" bson:"email"`
	GithubLink string 			`json:"githubLink" bson:"githubLink"`
	LinkedIn   string 			`json:"linkedIn" bson:"linkedIn"`
	EventSlug  string 			`validate:"required" json:"slug" bson:"slug"`
}


type Event struct {
	ID          primitive.ObjectID 	`json:"_id" bson:"_id"`
	Title       string             	`validate:"required" json:"title" bson:"title"`
	Slug        string             	`validate:"required" json:"slug" bson:"slug"`
	Description string             	`validate:"required" json:"description" bson:"description"`
	Icons	    []string           		`json:"icons" bson:"icons"`
	StartDate   string             	`validate:"required" json:"startDate" bson:"startDate"`
	EventCover  string 			`validate:"required" json:"eventCover" bson:"eventCover"` // s3-url for event cover Image
	IsCompleted bool   			`json:"isCompleted" bson:"isCompleted"`
	RSVP_Users []string 			`json:"rsvp_users" bson:"rsvp_users"`
}
