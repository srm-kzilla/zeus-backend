package eventController

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	eventModel "github.com/srm-kzilla/events/api/events/model"
	userModel "github.com/srm-kzilla/events/api/users/model"
	"github.com/srm-kzilla/events/database"
	"github.com/srm-kzilla/events/validators"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Get all Events Route
func GetAllEvents(c *fiber.Ctx) error {
	var events []eventModel.Event
	eventsCollection, e := database.GetCollection("zeus_Events", "Events")
	if e != nil {
		fmt.Println("Error: ", e)
		c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": e.Error(),
		})
	}

	cursor, err := eventsCollection.Find(context.Background(), bson.D{})
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

// FIXME: Some of the data is not passsing in the database
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


	var check eventModel.Event
	eventsCollection.FindOne(context.Background(), bson.M{"slug": event.Slug}).Decode(&check)
	if check.Slug == event.Slug {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Event already exists",
		})
		return nil
	}

	event.ID = primitive.NewObjectID()
	event.Slug = strings.ToLower(event.Slug)
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
	var event eventModel.Event
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
	err := eventsCollection.FindOne(context.Background(), bson.M{"_id": objId}).Decode(&event)
	if err != nil {
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
	var event eventModel.Event
	var slug = c.Params("slug")

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
	err := eventsCollection.FindOne(context.Background(), bson.M{"slug": slug}).Decode(&event)
	if err != nil {
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
	var slug = c.Query("slug")

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
