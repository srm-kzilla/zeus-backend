package eventModel

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Speaker struct {
	ID         primitive.ObjectID `json:"_id" bson:"_id"`
	Name       string             `validate:"required" json:"name" bson:"name"`
	Email      string             `json:"email" bson:"email"`
	EventSlug  string             `validate:"required" json:"slug" bson:"slug"`
	LinkedIn   string             `json:"linkedIn" bson:"linkedIn"`
	GithubLink string             `json:"githubLink" bson:"githubLink"`
	Image      string             `json:"image" bson:"image"`
	About      string             `json:"about" bson:"about"`
}

type Timeline struct {
	Date        string `validate:"required" json:"date" bson:"date"`
	Title       string `validate:"required" json:"title" bson:"title"`
	Description string `json:"description" bson:"description"`
}

type Prizes struct {
	Amount      string `validate:"required" json:"amount" bson:"amount"`
	Description string `validate:"required" json:"description" bson:"description"`
	Asset       string `validate:"required" json:"asset" bson:"asset"`
	Sponsor     string `validate:"required" json:"sponsor" bson:"sponsor"`
}

type Event struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id"`
	Title        string             `validate:"required" json:"title" bson:"title"`
	Slug         string             `validate:"required" json:"slug" bson:"slug"`
	Description  string             `validate:"required" json:"description" bson:"description"`
	Tagline      string             `validate:"required" json:"tagline" bson:"tagline"`
	Timeline     []Timeline         `validate:"required" json:"timeline" bson:"timeline"`
	Prizes       []Prizes           `validate:"required" json:"prizes" bson:"prizes"`
	Icons        []string           `json:"icons" bson:"icons"`
	StartDate    string             `validate:"required" json:"startDate" bson:"startDate"`
	EventCover   string             `validate:"required" json:"eventCover" bson:"eventCover"`
	IsCompleted  bool               `default:"true" json:"isCompleted" bson:"isCompleted"`
	IsRegClosed  bool               `default:"true" json:"isRegClosed" bson:"isRegClosed"`
	RSVPUsers    []string           `json:"rsvpUsers" bson:"rsvpUsers"`
	MaxRsvp      int                `validate:"required" json:"maxRsvp" bson:"maxRsvp"`
	SocialHandle []SocialHandle     `validate:"required" json:"socialHandle" bson:"socialHandle"`
}

type EventWithSpeakers struct {
	Event
	Speakers []Speaker `json:"speakers" bson:"speakers"`
}

type SocialHandle struct {
	MediaType   string `validate:"required" json:"mediaType" bson:"mediaType"`
	MediaHandle string `validate:"required" json:"mediaHandle" bson:"mediaHandle"`
}
