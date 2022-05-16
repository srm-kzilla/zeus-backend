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
)

// Get all Events Route
func GetAllEvents(c *fiber.Ctx) error {
	var events []bson.M
	eventsCollection, e := database.GetCollection("zeus_Events", "Events")
	if e != nil {
		fmt.Println("Error: ", e)
		c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": e.Error(),
		})
	}
	lookupStage := bson.D{{"$lookup", bson.D{{"from","Speakers"}, {"localField","slug"}, {"foreignField","slug"},{"as","speakers"}}}}
	cursor, err := eventsCollection.Aggregate(context.Background(), mongo.Pipeline{lookupStage})
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

func CreateEvent(c *fiber.Ctx) error {
	var event eventModel.Event

	c.BodyParser(&event)

	errors := validators.ValidateEvents(event)
	if errors != nil {
		c.Status(fiber.StatusBadGateway).JSON(errors)
		return nil
	}

	eventsCollection, e := database.GetCollection("zeus_Events", "Events")
	if e != nil {
		fmt.Println("Error: ", e)
		c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": e.Error(),
		})
	}


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

	eventsCollection, e := database.GetCollection("zeus_Events", "Events")
	if e != nil {
		fmt.Println("Error: ", e)
		c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": e.Error(),
		})
	}
	matchId := bson.D{{"$match", bson.D{{"_id",objId}}}}
	lookupStage := bson.D{{"$lookup", bson.D{{"from","Speakers"}, {"localField","slug"}, {"foreignField","slug"},{"as","speakers"}}}}
	cur, err := eventsCollection.Aggregate(context.Background(), mongo.Pipeline{matchId, lookupStage})
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

	eventsCollection, e := database.GetCollection("zeus_Events", "Events")
	if e != nil {
		fmt.Println("Error: ", e)
		c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": e.Error(),
		})
	}
	matchSlug := bson.D{{"$match", bson.D{{"slug",slug}}}}
	lookupStage := bson.D{{"$lookup", bson.D{{"from","Speakers"}, {"localField","slug"}, {"foreignField","slug"},{"as","speakers"}}}}
	// err := eventsCollection.FindOne(context.Background(), bson.M{"slug": slug}).Decode(&event)
	cur, err := eventsCollection.Aggregate(context.Background(), mongo.Pipeline{matchSlug, lookupStage})
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

func GetEventUsers(c *fiber.Ctx) error {
	var users []userModel.User
	var slug = strings.ToLower(c.Query("slug"))

	usersCollection, e := database.GetCollection("zeus_Events", "Users")
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
	if len(users) == 0 {
		c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "No users found",
		})
		return nil
	}
	c.Status((fiber.StatusOK)).JSON(users)

	return nil
}

func CloseEvent(c *fiber.Ctx) error {
	var event eventModel.Event
	var slug = strings.ToLower(c.Query("slug"))
	eventsCollection, e := database.GetCollection("zeus_Events", "Events")
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
	eventsCollection.FindOneAndReplace(context.Background(), bson.M{"slug": slug}, event)
	c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"event":   event,
		"message":"Event is successfully closed",
	})
	return nil
}

func UploadEventCover(c *fiber.Ctx) error {
	var slug  = c.Query("slug")
	if slug == "" {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Slug is required",
		})
		return nil
	}
	file, err := c.FormFile("cover")
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
	var filePath string = "./"+file.Filename
	S3.UploadFile(buf, filePath, file.Size)
	c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "File uploaded successfully",
		"key": os.Getenv("S3_LINK") + file.Filename,
	})
	return nil
}

func AddSpeaker(c *fiber.Ctx) error {
	var speaker eventModel.Speaker

	c.BodyParser(&speaker)

	errors := validators.ValidateSpeaker(speaker)
	if errors != nil {
		c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": errors,
		})
		return nil
	}

	speaker.EventSlug = strings.ToLower(speaker.EventSlug)
	eventsCollection, e := database.GetCollection("zeus_Events", "Events")
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
	speaker.ID = primitive.NewObjectID()
	speakerCollection, e := database.GetCollection("zeus_Events", "Speakers")
	if e != nil {
		fmt.Println("Error: ", e)
		c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": e.Error(),
		})
		return nil
	}
	res, err := speakerCollection.InsertOne(context.Background(), speaker)
	if err != nil {
		log.Println("Error ", err)
		c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": err.Error(),
			"InsertedId": res.InsertedID,
		})
		return nil
	}

	c.Status(fiber.StatusOK).JSON(speaker)

	return nil
}

func UpdateEvent(c *fiber.Ctx) error {
	var event eventModel.Event
	var check eventModel.Event
	c.BodyParser(&event)

	errors := validators.ValidateEvents(event)
	if errors != nil {
		c.Status(fiber.StatusBadGateway).JSON(errors)
		return nil
	}

	eventsCollection, e := database.GetCollection("zeus_Events", "Events")
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
	event.RSVP_Users = check.RSVP_Users
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