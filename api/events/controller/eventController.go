package eventController

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	eventModel "github.com/srm-kzilla/events/api/events/model"
	userModel "github.com/srm-kzilla/events/api/users/model"
	"github.com/srm-kzilla/events/database"
	helpers "github.com/srm-kzilla/events/utils/helpers"
	S3 "github.com/srm-kzilla/events/utils/services/s3"
	"github.com/srm-kzilla/events/validators"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/***********************************************
Get all Events present in the Events Collection.
***********************************************/
func GetAllEvents(c *fiber.Ctx) error {
	var events []bson.M
	eventsCollection, e := database.GetCollection(os.Getenv("DB_NAME"), "Events")
	if e != nil {
		fmt.Println("Error: ", e)
		c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": e.Error(),
		})
	}
	lookupStage := bson.D{{"$lookup", bson.D{{"from", "Speakers"}, {"localField", "slug"}, {"foreignField", "slug"}, {"as", "speakers"}}}}
	cursor, err := eventsCollection.Aggregate(context.Background(), mongo.Pipeline{lookupStage, bson.D{{"$sort", bson.D{{"_id", -1}}}}, bson.D{{"$project",bson.D{{"rsvpUsers",0}}}}})
	if err = cursor.All(context.Background(), &events); err != nil {
		log.Println("Error ", err)
		c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":  err.Error(),
			"events": events,
		})
	}
	c.Status(fiber.StatusOK).JSON(events)

	return nil
}

/***********************
   Create a new Event.
***********************/
func CreateEvent(c *fiber.Ctx) error {
	var event eventModel.Event

	c.BodyParser(&event)

	errors := validators.ValidateEvents(event)
	if errors != nil {
		c.Status(fiber.StatusBadGateway).JSON(errors)
		return nil
	}

	eventsCollection, e := database.GetCollection(os.Getenv("DB_NAME"), "Events")
	if e != nil {
		fmt.Println("Error: ", e)
		c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": e.Error(),
		})
	}
	eventsCollection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.M{"slug": 1},
		Options: options.Index().SetUnique(true),
	})

	event.Slug = strings.ToLower(event.Slug)
	var check eventModel.Event
	eventsCollection.FindOne(context.Background(), bson.M{"slug": event.Slug}).Decode(&check)
	if check.Slug == event.Slug {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Event already exists",
		})
		return nil
	}

	event.ID = primitive.NewObjectID()
	res, err := eventsCollection.InsertOne(context.Background(), event)
	if err != nil {
		log.Println("Error", err)
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Success":     false,
			"Inserted ID": res.InsertedID,
		})
		return err
	}

	c.Status(fiber.StatusCreated).JSON(event)

	return nil
}

/****************************************************************
Get a particular Event's data from the Collection using ObjectID.
****************************************************************/
func GetEventById(c *fiber.Ctx) error {
	var event []bson.M
	var id = c.Query("id")
	objId, _ := primitive.ObjectIDFromHex(id)
	if id == "" {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID is required",
		})
		return nil
	}

	eventsCollection, e := database.GetCollection(os.Getenv("DB_NAME"), "Events")
	if e != nil {
		fmt.Println("Error: ", e)
		c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": e.Error(),
		})
	}
	matchId := bson.D{{"$match", bson.D{{"_id", objId}}}}
	lookupStage := bson.D{{"$lookup", bson.D{{"from", "Speakers"}, {"localField", "slug"}, {"foreignField", "slug"}, {"as", "speakers"}}}}
	cur, err := eventsCollection.Aggregate(context.Background(), mongo.Pipeline{matchId, lookupStage, bson.D{{"$project", bson.D{{"rsvpUsers", 0}}}}})
	if cur.All(context.Background(), &event); err != nil {
		log.Println("Error ", err)
		c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": err.Error(),
		})
		return nil
	}

	c.Status(fiber.StatusOK).JSON(event)
	return nil

}

/******************************************
Get Event from collection using Event Slug.
******************************************/
func GetEventBySlug(c *fiber.Ctx) error {
	// var event eventModel.Event
	var event []bson.M
	var slug = strings.ToLower(c.Params("slug"))

	if slug == "" {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Slug is required",
		})
		return nil
	}

	eventsCollection, e := database.GetCollection(os.Getenv("DB_NAME"), "Events")
	if e != nil {
		fmt.Println("Error: ", e)
		c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": e.Error(),
		})
	}
	matchSlug := bson.D{{"$match", bson.D{{"slug", slug}}}}
	lookupStage := bson.D{{"$lookup", bson.D{{"from", "Speakers"}, {"localField", "slug"}, {"foreignField", "slug"}, {"as", "speakers"}}}}
	cur, err := eventsCollection.Aggregate(context.Background(), mongo.Pipeline{matchSlug, lookupStage, bson.D{{"$project", bson.D{{"rsvpUsers", 0}}}}})
	if cur.All(context.Background(), &event); err != nil {
		log.Println("Error ", err)
		c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": err.Error(),
		})
		return nil
	}
	c.Status(fiber.StatusOK).JSON(event)
	return nil

}

/**************************************************
Get all Users of a specific Event using Event Slug.
**************************************************/
func GetEventUsers(c *fiber.Ctx) error {
	var users []userModel.User
	var slug = strings.ToLower(c.Query("slug"))

	usersCollection, e := database.GetCollection(os.Getenv("DB_NAME"), "Users")
	if e != nil {
		fmt.Println("Error: ", e)
		c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": e.Error(),
		})
	}

	cursor, err := usersCollection.Find(context.Background(), bson.M{"events": bson.M{"$in": []string{slug}}})
	if err = cursor.All(context.Background(), &users); err != nil {
		log.Println("Error ", err)
		c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": err.Error(),
			"users": users,
		})
	}
	c.Status((fiber.StatusOK)).JSON(fiber.Map{
		"users":      users,
		"numOfUsers": len(users),
	})

	return nil
}

/****************************
Close Event using Event Slug.
****************************/
func CloseEvent(c *fiber.Ctx) error {
	var event eventModel.Event
	var slug = strings.ToLower(c.Query("slug"))
	eventsCollection, e := database.GetCollection(os.Getenv("DB_NAME"), "Events")
	if e != nil {
		fmt.Println("Error: ", e)
		c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": e.Error(),
		})
		return nil
	}
	err := eventsCollection.FindOne(context.Background(), bson.M{"slug": slug}).Decode(&event)
	if err != nil {
		log.Println("Error ", err)
		c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": err.Error(),
		})
		return nil
	}
	event.IsCompleted = true
	event.IsRegClosed = true
	eventsCollection.FindOneAndReplace(context.Background(), bson.M{"slug": slug}, event)
	c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"event":   event,
		"message": "Event is successfully closed",
	})
	return nil
}

/******************************************
Close Event Registrations using Event slug.
******************************************/
func CloseRegistrations(c *fiber.Ctx) error {
	var event eventModel.Event
	var slug = strings.ToLower(c.Query("slug"))
	eventsCollection, e := database.GetCollection(os.Getenv("DB_NAME"), "Events")
	if e != nil {
		fmt.Println("Error: ", e)
		c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": e.Error(),
		})
		return nil
	}
	err := eventsCollection.FindOne(context.Background(), bson.M{"slug": slug}).Decode(&event)
	if err != nil {
		log.Println("Error ", err)
		c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": err.Error(),
		})
		return nil
	}
	event.IsRegClosed = true
	eventsCollection.FindOneAndReplace(context.Background(), bson.M{"slug": slug}, event)
	c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"event":   event,
		"message": "Event Registrations is successfully closed",
	})
	return nil
}

/****************************************
Upload Event Cover File using Event Slug.
****************************************/
func UploadEventCover(c *fiber.Ctx) error {
	var slug = c.Query("slug")
	if slug == "" {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Slug is required",
		})
		return nil
	}
	file, err := c.FormFile("file")
	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
		return nil
	}
	fileBody, _ := file.Open()
	buf, e := ioutil.ReadAll(fileBody)
	if e != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
		return nil
	}

	file.Filename = fmt.Sprintf("%s/covers/%s.%s", slug, helpers.GenerateNanoID(10), strings.Split(file.Filename, ".")[1])
	var filePath string = "./" + file.Filename
	S3.UploadFile(buf, filePath, file.Size)
	c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "File uploaded successfully",
		"key":     os.Getenv("S3_LINK") + file.Filename,
	})
	return nil
}

/****************************************
Add Speakers to a Event using Event Slug.
****************************************/
func AddSpeaker(c *fiber.Ctx) error {
	var speaker eventModel.Speaker
	var check eventModel.Speaker
	c.BodyParser(&speaker)

	errors := validators.ValidateSpeaker(speaker)
	if errors != nil {
		c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": errors,
		})
		return nil
	}

	speaker.EventSlug = strings.ToLower(speaker.EventSlug)
	eventsCollection, e := database.GetCollection(os.Getenv("DB_NAME"), "Events")
	if e != nil {
		fmt.Println("Error: ", e)
		c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": e.Error(),
		})
		return nil
	}
	var event eventModel.Event
	err := eventsCollection.FindOne(context.Background(), bson.M{"slug": speaker.EventSlug}).Decode(&event)
	if err != nil {
		log.Println("Error ", err)
		c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": "no such event/eventSlug exists",
		})
		return nil
	}
	speakerCollection, e := database.GetCollection(os.Getenv("DB_NAME"), "Speakers")
	if e != nil {
		fmt.Println("Error: ", e)
		c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": e.Error(),
		})
		return nil
	}
	speakerCollection.FindOne(context.Background(), bson.M{"email": speaker.Email}).Decode(&check)
	if check.Email == speaker.Email {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Speaker with this email already exists",
		})
	}
	speaker.ID = primitive.NewObjectID()
	res, err := speakerCollection.InsertOne(context.Background(), speaker)
	if err != nil {
		log.Println("Error ", err)
		c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":      err.Error(),
			"InsertedId": res.InsertedID,
		})
		return nil
	}

	c.Status(fiber.StatusOK).JSON(speaker)

	return nil
}

/*************************************
Update Event Details using Event Slug.
*************************************/
func UpdateEvent(c *fiber.Ctx) error {
	var event eventModel.Event
	var check eventModel.Event
	c.BodyParser(&event)

	errors := validators.ValidateEvents(event)
	if errors != nil {
		c.Status(fiber.StatusBadGateway).JSON(errors)
		return nil
	}
	if event.ID == primitive.NilObjectID {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": "ObjectID is required",
		})
	}
	eventsCollection, e := database.GetCollection(os.Getenv("DB_NAME"), "Events")
	if e != nil {
		fmt.Println("Error: ", e)
		c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": e.Error(),
		})
	}
	event.Slug = strings.ToLower(event.Slug)
	err := eventsCollection.FindOne(context.Background(), bson.M{"slug": event.Slug}).Decode(&check)
	if err != nil {
		log.Println("Error ", err)
		c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": "no such event/eventSlug exists",
		})
		return nil
	}
	event.RSVPUsers = check.RSVPUsers
	errr := eventsCollection.FindOneAndReplace(context.Background(), bson.M{"slug": event.Slug}, event).Decode(&check)
	if errr != nil {
		fmt.Println("Error: ", errr)
		c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": errr.Error(),
		})
		return nil
	}
	c.Status(fiber.StatusOK).JSON(event)
	return nil
}

/***************************************
Update Speaker Details using Event Slug.
***************************************/
func UpdateSpeaker(c *fiber.Ctx) error {
	var speaker eventModel.Speaker
	var check eventModel.Speaker

	c.BodyParser(&speaker)

	errors := validators.ValidateSpeaker(speaker)
	if errors != nil {
		return c.Status(fiber.StatusBadGateway).JSON(errors)
	}
	if speaker.ID == primitive.NilObjectID {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": "ObjectID is required",
		})
	}
	speaker.EventSlug = strings.ToLower(speaker.EventSlug)

	speakerCollection, e := database.GetCollection(os.Getenv("DB_NAME"), "Speakers")
	if e != nil {
		fmt.Println("Error: ", e)
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": e.Error(),
		})
	}
	eventsCollection, e := database.GetCollection(os.Getenv("DB_NAME"), "Events")
	if e != nil {
		fmt.Println("Error: ", e)
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": e.Error(),
		})
	}
	var event eventModel.Event
	err := eventsCollection.FindOne(context.Background(), bson.M{"slug": speaker.EventSlug}).Decode(&event)
	if err != nil {
		log.Println("Error ", err)
		c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": "no such event/eventSlug exists",
		})
		return nil
	}
	errr := speakerCollection.FindOne(context.Background(), bson.M{"email": speaker.Email}).Decode(&check)
	if errr != nil {
		log.Println("Error ", errr)
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": "no such speaker exists",
		})
	}
	errrr := speakerCollection.FindOneAndReplace(context.Background(), bson.M{"email": speaker.Email}, speaker).Decode(&check)
	if errrr != nil {
		fmt.Println("Error: ", errrr)
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": errr.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(speaker)

}
