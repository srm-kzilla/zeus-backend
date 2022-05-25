package userController

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	eventModel "github.com/srm-kzilla/events/api/events/model"
	userModel "github.com/srm-kzilla/events/api/users/model"
	"github.com/srm-kzilla/events/database"
	"github.com/srm-kzilla/events/utils/helpers"
	"github.com/srm-kzilla/events/utils/services/mailer"
	qr "github.com/srm-kzilla/events/utils/services/qrcode"
	"github.com/srm-kzilla/events/validators"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func RegisterForEvent(c *fiber.Ctx) error {
	var reqBody userModel.RegisterUserReq
	c.BodyParser(&reqBody)

	E := validators.ValidateRegisterUserReq(reqBody)
	if E != nil {
		c.Status(fiber.StatusBadRequest).JSON(E)
		return nil
	}

	var user userModel.User = reqBody.User
	reqBody.EventSlug = strings.ToLower(reqBody.EventSlug)

	errors := validators.ValidateUser((user))
	if errors != nil {
		c.Status(fiber.StatusBadRequest).JSON(errors)
		return nil
	}
	usersCollection, e := database.GetCollection(os.Getenv("DB_NAME"), "Users")
	if e != nil {
		fmt.Println("Error: ", e)
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   e.Error(),
			"message": "Collection Not found",
		})
	}

	usersCollection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.M{"email": 1, "phoneNumber": 1},
		Options: options.Index().SetUnique(true),
	})

	eventsCollection, e := database.GetCollection(os.Getenv("DB_NAME"), "Events")
	if e != nil {
		fmt.Println("Error: ", e)
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   e.Error(),
			"message": "Collection Not found",
		})
	}
	var event eventModel.Event
	err := eventsCollection.FindOne(context.Background(), bson.M{"slug": reqBody.EventSlug}).Decode(&event)
	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No such event/eventSlug exists",
		})
		return nil
	}
	if event.IsCompleted {
		c.Status(fiber.StatusLocked).JSON(fiber.Map{
			"error": "Event is already completed",
		})
		return nil
	}
	var check userModel.User
	usersCollection.FindOne(context.Background(), bson.M{
		"$or": []bson.M{
			{"email": user.Email},
			{"phoneNumber": user.PhoneNumber},
			{"regNumber": user.RegNumber},
		},
	}).Decode(&check)
	if check.Email == user.Email || check.PhoneNumber == user.PhoneNumber || check.RegNumber == user.RegNumber {
		// c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		// 	"error":   "User with that email already exists",
		// })
		// return nil
		if helpers.ExistsInArray(check.EventSlugs, reqBody.EventSlug) {
			c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "User already registered for this event",
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
	newUserEmbed := mailer.NewUserEmbed{
		Name: user.Name,
	}
	sesInput := mailer.SESInput{
		TemplateName:  mailer.TEMPLATES.NewUserTemplate,
		Subject:       "Registration Successfull",
		Name:          user.Name,
		RecieverEmail: user.Email,
		SenderEmail:   os.Getenv("SENDER_EMAIL"),
		EmbedData:     newUserEmbed,
	}
	mailer.SendEmail(sesInput)
	c.Status(fiber.StatusCreated).JSON(user)
	return nil
}

func RsvpForEvent(c *fiber.Ctx) error {
	var reqBody userModel.RsvpUserReq
	c.QueryParser(&reqBody)
	reqBody.EventSlug = strings.ToLower(reqBody.EventSlug)

	E := validators.ValidateRsvpUserReq(reqBody)
	if E != nil {
		c.Status(fiber.StatusBadRequest).JSON(E)
		return nil
	}

	usersCollection, e := database.GetCollection(os.Getenv("DB_NAME"), "Users")
	if e != nil {
		fmt.Println("Error: ", e)
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   e.Error(),
			"message": "Collection Not found",
		})
	}
	eventsCollection, e := database.GetCollection(os.Getenv("DB_NAME"), "Events")
	if e != nil {
		fmt.Println("Error: ", e)
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   e.Error(),
			"message": "Collection Not found",
		})
	}
	var event eventModel.Event
	errr := eventsCollection.FindOne(context.Background(), bson.M{"slug": reqBody.EventSlug}).Decode(&event)
	if errr != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No such event/eventSlug exists",
		})
		return nil
	}
	if event.IsCompleted {
		c.Status(fiber.StatusLocked).JSON(fiber.Map{
			"error": "Event is already completed",
		})
	}
	var user userModel.User
	objId, _ := primitive.ObjectIDFromHex(reqBody.UserId)
	err := usersCollection.FindOne(context.Background(), bson.M{"_id": objId}).Decode(&user)
	if err != nil {
		log.Println("Error", err)
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Success": false,
			"Error":   "User not found",
		})
		return nil
	}
	if !helpers.ExistsInArray(user.EventSlugs, reqBody.EventSlug) {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User not registered for this event",
		})
		return nil
	}

	if helpers.ExistsInArray(event.RSVPUsers, reqBody.UserId) {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User already RSVPed for this event",
		})
		return nil
	}
	event.RSVPUsers = append(event.RSVPUsers, reqBody.UserId)
	rsvpEmbed := mailer.RsvpEmbed{
		QrLink: qr.GenerateQRCode(user.ID.Hex()),
	}
	eventsCollection.FindOneAndReplace(context.Background(), bson.M{"slug": reqBody.EventSlug}, event)
	sesInput := mailer.SESInput{
		TemplateName:  mailer.TEMPLATES.RsvpTemplate,
		Subject:       "RSVP Successfull",
		Name:          user.Name,
		RecieverEmail: user.Email,
		SenderEmail:   os.Getenv("SENDER_EMAIL"),
		EmbedData:     rsvpEmbed,
	}
	mailer.SendEmail(sesInput)
	c.Status(fiber.StatusOK).JSON(fiber.Map{
		"Success": true,
		"Message": "User RSVPed for event",
	})
	return nil
}

func GetUserById(c *fiber.Ctx) error {
	userId := c.Params("userId")
	if userId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "UserId is required",
		})
	}
	objId, _ := primitive.ObjectIDFromHex(userId)
	usersCollection, e := database.GetCollection(os.Getenv("DB_NAME"), "Users")
	if e != nil {
		fmt.Println("Error: ", e)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   e.Error(),
			"message": "Error getting user collection",
		})
	}
	var user userModel.User
	err := usersCollection.FindOne(context.Background(), bson.M{"_id": objId}).Decode(&user)
	if err != nil {
		log.Println("Error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   err.Error(),
			"message": "User not found",
		})
	}
	return c.Status(fiber.StatusOK).JSON(user)
}
