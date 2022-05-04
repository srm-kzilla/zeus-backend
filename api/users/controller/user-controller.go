package userController

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	eventModel "github.com/srm-kzilla/events/api/events/model"
	userModel "github.com/srm-kzilla/events/api/users/model"
	"github.com/srm-kzilla/events/database"
	"github.com/srm-kzilla/events/utils/constants"
	"github.com/srm-kzilla/events/utils/services/mailer"
	"github.com/srm-kzilla/events/validators"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func RegisterForEvent(c *fiber.Ctx) error {
	var reqBody userModel.RegisterUser
	c.BodyParser(&reqBody)

	var user userModel.User = reqBody.User


	errors := validators.ValidateUser((user))
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
	eventsCollection, e := database.GetCollection("zeus_Events", "Events")
	if e != nil {
		fmt.Println("Error: ", e)
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   e.Error(),
			"message": "Collection Not found ⚠️",
		})
	}
	var event eventModel.Event
	err := eventsCollection.FindOne(context.Background(), bson.M{"slug": reqBody.EventSlug}).Decode(&event)
	if err != nil {
		c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": "No such event/eventSlug exists",
		})
		return nil
	}
	var check userModel.User
	usersCollection.FindOne(context.Background(), bson.M{"email": user.Email}).Decode(&check)
	if check.Email == user.Email {
		// c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
		// 	"error":   "User with that email already exists",
		// })
		// return nil
		if constants.ExistsInArray(check.EventSlugs, reqBody.EventSlug) {
			c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
				"error":   "User already registered for this event",
			})
			return nil
		}
		check.EventSlugs = append(check.EventSlugs, reqBody.EventSlug)
		usersCollection.FindOneAndReplace(context.Background(), bson.M{"email": user.Email}, check)
		c.Status(fiber.StatusCreated).JSON(check)
		return nil
	}
	user.ID = primitive.NewObjectID()
	user.EventSlugs = append(user.EventSlugs, reqBody.EventSlug)
	res, err := usersCollection.InsertOne(context.Background(), user)
	if err != nil {
		log.Println("Error", err)
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Success":     false,
			"Inserted ID": res.InsertedID,
		})
		return err
	}
	sesInput := mailer.SESInput{
		TemplateName: "newUser.html",
		Subject: "Registration Successfull",
		Name: user.Name,
		RecieverEmail: user.Email,
		SenderEmail: os.Getenv("SENDER_EMAIL"),
	}
	mailer.SendEmail(sesInput)
	c.Status(fiber.StatusCreated).JSON(user)
	return nil
}