package controller

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	model "github.com/srm-kzilla/events/api/events/model"
	"github.com/srm-kzilla/events/api/events/services/mailer"
	"github.com/srm-kzilla/events/database"
	"github.com/srm-kzilla/events/validators"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Get all Events Route
func GetAllEvents(c *fiber.Ctx) error {
	var events []model.Event
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
	var event model.Event

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

// TODO: Event Regsitration Route handler
func RegisterForEvent(c *fiber.Ctx) error {
	var user model.User
	c.BodyParser(&user)

	errors := validators.ValidateUser(user)
	if errors != nil {
		c.Status(fiber.StatusBadGateway).JSON(errors)
		return nil
	}

	usersCollection, e := database.GetCollection("zeus_Events", "Users")
	if e != nil {
		fmt.Println("Error: ", e)
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   e.Error(),
			"message": "Collection Not found ⚠️",
		})
	}
	user.ID = primitive.NewObjectID()
	res, err := usersCollection.InsertOne(context.Background(), user)
	if err != nil {
		log.Println("Error", err)
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Success":     false,
			"Inserted ID": res.InsertedID,
		})
		return err
	}

	senderEmail := os.Getenv("SENDER_EMAIL")

	sesInput := model.SESInput{
		TemplateName:  "newUser.html",
		Subject:       "Registration Successfull",
		Name:          user.Name,
		RecieverEmail: user.Email,
		SenderEmail:   senderEmail,
	}

	mailer.SendEmail(sesInput)

	c.Status(fiber.StatusCreated).JSON(user)

	return nil
}

func GetEventById(c *fiber.Ctx) error {
	var event model.Event
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
	var event model.Event
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
	var users []model.User
	var slug = c.Query("slug")

	usersCollection, e := database.GetCollection("zeus_Events", "Users")
	if e != nil {
		fmt.Println("Error: ", e)
		c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": e.Error(),
		})
	}

	cursor, err := usersCollection.Find(context.Background(), bson.M{"slug": slug})
	if err = cursor.All(context.Background(), &users); err != nil {
		log.Println("Error ", err)
		c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": err.Error(),
			"users": users,
		})
	}
	c.Status((fiber.StatusOK)).JSON(users)

	return nil
}
