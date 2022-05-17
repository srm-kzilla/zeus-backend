package eventModel

import (
	userModel "github.com/srm-kzilla/events/api/users/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Speaker struct {
	ID 	     primitive.ObjectID 	`json:"_id" bson:"_id"`
	Name       string				`validate:"required" json:"name" bson:"name"`
	Email      string 				`json:"email" bson:"email"`
	EventSlug  string 				`validate:"required" json:"slug" bson:"slug"`
	LinkedIn   string 				`json:"linkedIn" bson:"linkedIn"`
	GithubLink string 				`json:"githubLink" bson:"githubLink"`
	Image	   string 				`json:"image" bson:"image"`
	About	   string				`json:"about" bson:"about"`
}

type Timeline struct {
	Date	string		`validate:"required" json:"date" bson:"date"`
	Title	string		`validate:"required" json:"title" bson:"title"`
	Description	string	`json:"description" bson:"description"`
}

type Prizes struct {
	Amount 		string 		`validate:"required" json:"amount" bson:"amount"`
	Description string 		`validate:"required" json:"description" bson:"description"`
	Asset		string 		`validate:"required" json:"asset" bson:"asset"`
	Sponsor		string 		`validate:"required" json:"sponsor" bson:"sponsor"`
}

type Event struct {
	ID          primitive.ObjectID 		`json:"_id" bson:"_id"`
	Title       string             		`validate:"required" json:"title" bson:"title"`
	Slug        string             		`validate:"required" json:"slug" bson:"slug"`
	Description string             		`validate:"required" json:"description" bson:"description"`
	Tagline		string					`validate:"required" json:"tagline" bson:"tagline"`
	Timeline	[]Timeline				`validate:"required" json:"timeline" bson:"timeline"`
	Prizes		[]Prizes				`validate:"required" json:"prizes" bson:"prizes"`
	Icons	    []string           		`json:"icons" bson:"icons"`
	StartDate   string             		`validate:"required" json:"startDate" bson:"startDate"`
	EventCover  string 					`validate:"required" json:"eventCover" bson:"eventCover"`
	IsCompleted bool   					`json:"isCompleted" bson:"isCompleted"`
	RSVP_Users  []userModel.RsvpUsers 	`json:"rsvp_users" bson:"rsvp_users"`
}

type EventWithSpeakers struct {
	Event
	Speakers []Speaker	`json:"speakers" bson:"speakers"`
}