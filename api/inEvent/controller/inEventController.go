package InEventController

import (
	"context"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	EventModel "github.com/srm-kzilla/events/api/events/model"
	InEventModel "github.com/srm-kzilla/events/api/inEvent/model"
	UserModel "github.com/srm-kzilla/events/api/users/model"
	"github.com/srm-kzilla/events/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// func grantAttendance(c *fiber.Ctx, inEventDataCollection *mongo.Collection, userData *userModel.User, slug string) error {
// 	count, e := inEventDataCollection.CountDocuments(context.Background(), bson.M{"userId": userData.ID.Hex(), "eventSlug": slug})
// 	if e != nil {
// 		fmt.Println("Error", e)
// 		c.Status(500).JSON(fiber.Map{
// 			"message": "Error while checking pre existing attendance",
// 			"error":   e.Error(),
// 		})
// 		return e
// 	}
// 	fmt.Println(count)
// 	if count != int64(0) {
// 		fmt.Println("Attendance already granted")
// 		c.Status(409).JSON(fiber.Map{
// 			"message": "Attendance already granted",
// 		})
// 		return nil
// 	}
// 	doc := bson.D{{"userId", userData.ID.Hex()}, {"eventSlug", slug}, {"foodReceived", false}}
// 	result, err := inEventDataCollection.InsertOne(context.Background(), doc)
// 	if err != nil {
// 		fmt.Println("Error", err)
// 		c.Status(500).JSON(fiber.Map{
// 			"message": "Error while granting attendance",
// 			"error":   e.Error(),
// 		})
// 		return err
// 	}
// 	fmt.Println("Result", result)
// 	c.Status(200).JSON(fiber.Map{
// 		"message": "Attendance granted",
// 		"result":  result,
// 	})
// 	return nil
// }

// func handOverFood(c *fiber.Ctx, inEventDataCollection *mongo.Collection, userData *userModel.User, slug string) error {
// 	count, e := inEventDataCollection.CountDocuments(context.Background(), bson.M{"userId": userData.ID.Hex(), "eventSlug": slug, "foodReceived": true})
// 	if e != nil {
// 		fmt.Println("Error", e)
// 		c.Status(500).JSON(fiber.Map{
// 			"message": "Error while checking pre existing attendance",
// 			"error":   e.Error(),
// 		})
// 		return e
// 	}
// 	fmt.Println(count)
// 	if count != int64(0) {
// 		fmt.Println("Attendance already granted")
// 		c.Status(409).JSON(fiber.Map{
// 			"message": "Food already handed over",
// 		})
// 		return nil
// 	}
// 	result, err := inEventDataCollection.UpdateOne(context.Background(), bson.M{"userId": userData.ID.Hex(), "eventSlug": slug}, bson.D{{"$set", bson.D{{"foodReceived", true}}}})
// 	if err != nil {
// 		fmt.Println("Error", err)
// 		c.Status(500).JSON(fiber.Map{
// 			"message": "Error while logging handover",
// 			"error":   e.Error(),
// 		})
// 		return err
// 	}
// 	fmt.Println("Result", result)
// 	c.Status(200).JSON(fiber.Map{
// 		"message": "Hand over logged",
// 		"result":  result,
// 	})
// 	return nil
// }

// func InEventHandler(c *fiber.Ctx) error {
// 	attendanceQuery := new(inEventModel.AttendanceQuery)
// 	if err := c.QueryParser(attendanceQuery); err != nil {
// 		fmt.Println("Error", err)
// 		c.Status(400).JSON(fiber.Map{
// 			"message": "Query Fields missing",
// 			"error":   err.Error(),
// 		})
// 		return err
// 	}
// 	userCollection, e := database.GetCollection("zeus_Events", "Users")
// 	if e != nil {
// 		fmt.Println("Error", e)
// 		c.Status(500).JSON(fiber.Map{
// 			"message": "Error while getting user details",
// 			"error":   e.Error(),
// 		})
// 		return e
// 	}
// 	email := attendanceQuery.Email
// 	slug := attendanceQuery.Slug
// 	if &email == nil || &slug == nil {
// 		fmt.Println("Email or slug missing")
// 		c.Status(400).JSON(fiber.Map{
// 			"message": "Email or slug missing",
// 		})
// 	}
// 	fmt.Println("Email", email)
// 	fmt.Println("Slug", slug)
// 	var userData userModel.User
// 	err := userCollection.FindOne(context.Background(), bson.M{"email": email, "events": bson.M{"$in": []string{slug}}}).Decode(&userData)
// 	if err != nil {
// 		fmt.Println("Error", e)
// 		c.Status(500).JSON(fiber.Map{
// 			"message": "Error while fetching details of " + email,
// 			"error":   e.Error(),
// 		})
// 		return err
// 	}
// 	if &userData == nil {
// 		fmt.Println("No such user exists")
// 		c.Status(404).JSON(fiber.Map{
// 			"message": "User not found for this combination of email and slug",
// 		})
// 		return nil
// 	}
// 	inEventDataCollection, e := database.GetCollection("zeus_Events", "InEventData")
// 	if e != nil {
// 		fmt.Println("Error", e)
// 		c.Status(500).JSON(fiber.Map{
// 			"message": "Error in accessing in events data",
// 			"error":   e.Error(),
// 		})
// 		return e
// 	}
// 	switch action := c.Params("action"); action {
// 	case "attendance":
// 		return grantAttendance(c, inEventDataCollection, &userData, slug)
// 	case "food":
// 		return handOverFood(c, inEventDataCollection, &userData, slug)
// 	default:
// 		c.Status(404).JSON(fiber.Map{
// 			"message": "Invalid action chosen",
// 		})
// 		return nil
// 	}
// }

func grantAttendance(c *fiber.Ctx, inEventDataCollection *mongo.Collection, userData *UserModel.User, slug string) error {
	return nil
}

func handOverFood(c *fiber.Ctx, inEventDataCollection *mongo.Collection, userData *UserModel.User, slug string) error {
	return nil
}

func InEventHandler(c *fiber.Ctx)error {
	attendanceQuery := new(InEventModel.AttendanceQuery)
	if err := c.QueryParser(attendanceQuery); err != nil {
		fmt.Println("Error", err)
		return c.Status(400).JSON(fiber.Map{
			"message": "Query Fields missing",
			"error":   err.Error(),
		})
	}
	attendanceQuery.Slug = strings.ToLower(attendanceQuery.Slug)
	eventsCollection, err := database.GetCollection("zeus_Events", "Events")
	if err != nil {
		fmt.Println("Error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error while getting events collection",
			"error":   err.Error(),
		})
	}
	userCollection, err := database.GetCollection("zeus_Events", "Users")
	if err != nil {
		fmt.Println("Error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error while getting user collection",
			"error":   err.Error(),
		})
	}
	// inEventsCollection, err := database.GetCollection("zeus_Events", "InEvents")
	// if err != nil {
	// 	fmt.Println("Error", err)
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 		"message": "Error while getting inEvents collection",
	// 		"error":   err.Error(),
	// 	})
	// }
	var eventData EventModel.Event
	e := eventsCollection.FindOne(context.Background(), bson.M{"slug": attendanceQuery.Slug}).Decode(&eventData)
	if e != nil {
		fmt.Println("Error", e)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "no such event/slug exists",
		})
	}
	var userData UserModel.User
	error := userCollection.FindOne(context.Background(), bson.M{"_id":attendanceQuery.UserID, "events": bson.M{"$in": []string{attendanceQuery.Slug}}}).Decode(&userData)
	if error != nil {
		fmt.Println("Error", error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Invalid Id, no such user found or user not registered for this event",
			"error":   error.Error(),
		})
	}
	er := eventsCollection.FindOne(context.Background(), bson.M{"slug": attendanceQuery.Slug, "rsvpUsers": bson.M{"$in": []string{userData.ID.Hex()}}}).Decode(&eventData)
	if er != nil {
		fmt.Println("Error", er)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "User not Rsvped for this event",
		})
	}
	switch action := attendanceQuery.Action; action {
	case "attendance":
		// run the attendance code
	case "food":
		// run the food code
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid action chosen",
		})
	}
return nil
}