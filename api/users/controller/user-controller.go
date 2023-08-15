package userController

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	eventModel "github.com/srm-kzilla/events/api/events/model"
	userModel "github.com/srm-kzilla/events/api/users/model"
	"github.com/srm-kzilla/events/database"
	"github.com/srm-kzilla/events/utils/constants"
	"github.com/srm-kzilla/events/utils/helpers"
	"github.com/srm-kzilla/events/utils/services/mailer"
	qr "github.com/srm-kzilla/events/utils/services/qrcode"
	"github.com/srm-kzilla/events/validators"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
********************************************************************************
Get User data for registration and allocate the respective Event Slug to the user.
********************************************************************************
*/
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
	if event.IsRegClosed {
		return c.Status(fiber.StatusLocked).JSON(fiber.Map{
			"error": "Registrations are closed",
		})
	}
	var check userModel.User
	usersCollection.FindOne(context.Background(), bson.M{
		"$or": []bson.M{
			{"email": user.Email},
			{"phoneNumber": user.PhoneNumber},
			{"regNumber": user.RegNumber},
		},
	}).Decode(&check)
	fmt.Println(check)
	if !(check.Email == "") && !(check.Email == user.Email && check.PhoneNumber == user.PhoneNumber && check.RegNumber == user.RegNumber) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User already exists",
		})
	}
	if check.Email == user.Email && check.PhoneNumber == user.PhoneNumber && check.RegNumber == user.RegNumber {
		if helpers.ExistsInArray(check.EventSlugs, reqBody.EventSlug) {
			c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "User already registered for this event",
			})
			return nil
		}
		check.EventSlugs = append(check.EventSlugs, reqBody.EventSlug)
		usersCollection.FindOneAndReplace(context.Background(), bson.M{"email": user.Email}, check)
		user = check
		c.Status(fiber.StatusCreated).JSON(check)

	} else {
		user.ID = primitive.NewObjectID()
		user.CreatedAt = fmt.Sprintf("%v", time.Now().Unix())
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
	}
	newUserEmbed := mailer.NewUserEmbed{
		Name: user.Name,
	}
	sesInput := mailer.SESInput{
		TemplateName:  mailer.TEMPLATES.NewUserTemplate,
		Subject:       "Registration Successful",
		Name:          user.Name,
		RecieverEmail: user.Email,
		SenderEmail:   os.Getenv("SENDER_EMAIL"),
		EmbedData:     newUserEmbed,
	}
	mailer.SendEmail(sesInput)
	c.Status(fiber.StatusCreated).JSON(user)
	return nil
}

/*
*****************************************************************
Checks in the RSVP parameter for the particular user for the event.
*****************************************************************
*/
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
	fmt.Println("Error", errr)
	if errr != nil {
		c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
		message := "Hmmm, It seems like you are trying to RSVP for an event that does not exist."
		lottieFile := constants.Animations.EventDoesNotExist
		c.Status(fiber.StatusOK)
		return c.Render("rsvpConfirmationTemplate", fiber.Map{
			"Message":    message,
			"LottieFile": lottieFile,
		})
	}
	if event.IsCompleted {
		c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
		message := "Hey there! Sorry the event is already completed."
		lottieFile := constants.Animations.EventCompleted
		c.Status(fiber.StatusOK)
		return c.Render("rsvpConfirmationTemplate", fiber.Map{
			"Message":    message,
			"LottieFile": lottieFile,
		})
	}
	var user userModel.User
	objId, _ := primitive.ObjectIDFromHex(reqBody.UserId)
	err := usersCollection.FindOne(context.Background(), bson.M{"_id": objId}).Decode(&user)
	if err != nil {
		log.Println("Error", err)
		message := "User not Found."
		lottieFile := constants.Animations.EventDoesNotExist
		c.Status(fiber.StatusBadRequest)
		return c.Render("rsvpConfirmationTemplate", fiber.Map{
			"Message":    message,
			"LottieFile": lottieFile,
		})
	}
	if !helpers.ExistsInArray(user.EventSlugs, reqBody.EventSlug) {
		message := "User not registered for this event."
		lottieFile := constants.Animations.EventDoesNotExist
		c.Status(fiber.StatusBadRequest)
		return c.Render("rsvpConfirmationTemplate", fiber.Map{
			"Message":    message,
			"LottieFile": lottieFile,
		})
	}

	if helpers.ExistsInArray(event.RSVPUsers, reqBody.UserId) {
		c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
		message := "Hey there! Don't be so anxious. Your seat has been reserved."
		lottieFile := constants.Animations.AlreadyRsvpd
		c.Status(fiber.StatusOK)
		return c.Render("rsvpConfirmationTemplate", fiber.Map{
			"Message":    message,
			"LottieFile": lottieFile,
		})
	}
	if len(event.RSVPUsers) >= event.MaxRsvp {
		c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
		message := "We're booked to capacity! We hope to see you in our next event."
		lottieFile := constants.Animations.FullyBooked
		c.Status(fiber.StatusOK)
		return c.Render("rsvpConfirmationTemplate", fiber.Map{
			"Message":    message,
			"LottieFile": lottieFile,
		})
	}
	event.RSVPUsers = append(event.RSVPUsers, reqBody.UserId)
	rsvpEmbed := mailer.RsvpEmbed{
		QrLink: qr.GenerateQRCode(user.ID.Hex()),
	}
	eventsCollection.FindOneAndReplace(context.Background(), bson.M{"slug": reqBody.EventSlug}, event)
	sesInput := mailer.SESInput{
		TemplateName:  mailer.TEMPLATES.RsvpTemplate,
		Subject:       "RSVP Successful",
		Name:          user.Name,
		RecieverEmail: user.Email,
		SenderEmail:   os.Getenv("SENDER_EMAIL"),
		EmbedData:     rsvpEmbed,
	}
	mailer.SendEmail(sesInput)

	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	message := fmt.Sprintf("Your seat has been reserved for %s, join us at %s on %s to find your spark!", event.Title, strings.Split(event.Timeline[0].Date, " ")[1], event.StartDate)
	lottieFile := constants.Animations.RsvpSuccess
	c.Status(fiber.StatusOK)
	return c.Render("rsvpConfirmationTemplate", fiber.Map{
		"Message":    message,
		"LottieFile": lottieFile,
	})

}

/*
*******************************************************************
Get a particular User's data from the Collection using user ObjectID.
*******************************************************************
*/
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
