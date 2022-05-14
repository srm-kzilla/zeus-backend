package inEventController

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	inEventModel "github.com/srm-kzilla/events/api/inEvent/model"
	userModel "github.com/srm-kzilla/events/api/users/model"
	"github.com/srm-kzilla/events/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func grantAttendance(c *fiber.Ctx, inEventDataCollection *mongo.Collection, userData *userModel.User, slug string) error {
	count, e := inEventDataCollection.CountDocuments(context.Background(), bson.M{"userId": userData.ID.Hex(), "eventSlug": slug})
	if e != nil {
		fmt.Println("Error", e)
		c.Status(500).JSON(fiber.Map{
			"message": "Error while checking pre existing attendance",
			"error":   e.Error(),
		})
		return e
	}
	fmt.Println(count)
	if count != int64(0) {
		fmt.Println("Attendance already granted")
		c.Status(409).JSON(fiber.Map{
			"message": "Attendance already granted",
		})
		return nil
	}
	doc := bson.D{{"userId", userData.ID.Hex()}, {"eventSlug", slug}, {"foodReceived", false}}
	result, err := inEventDataCollection.InsertOne(context.Background(), doc)
	if err != nil {
		fmt.Println("Error", err)
		c.Status(500).JSON(fiber.Map{
			"message": "Error while granting attendance",
			"error":   e.Error(),
		})
		return err
	}
	fmt.Println("Result", result)
	c.Status(200).JSON(fiber.Map{
		"message": "Attendance granted",
		"result":  result,
	})
	return nil
}

func handOverFood(c *fiber.Ctx, inEventDataCollection *mongo.Collection, userData *userModel.User, slug string) error {
	count, e := inEventDataCollection.CountDocuments(context.Background(), bson.M{"userId": userData.ID.Hex(), "eventSlug": slug, "foodReceived": true})
	if e != nil {
		fmt.Println("Error", e)
		c.Status(500).JSON(fiber.Map{
			"message": "Error while checking pre existing attendance",
			"error":   e.Error(),
		})
		return e
	}
	fmt.Println(count)
	if count != int64(0) {
		fmt.Println("Attendance already granted")
		c.Status(409).JSON(fiber.Map{
			"message": "Food already handed over",
		})
		return nil
	}
	result, err := inEventDataCollection.UpdateOne(context.Background(), bson.M{"userId": userData.ID.Hex(), "eventSlug": slug}, bson.D{{"$set", bson.D{{"foodReceived", true}}}})
	if err != nil {
		fmt.Println("Error", err)
		c.Status(500).JSON(fiber.Map{
			"message": "Error while logging handover",
			"error":   e.Error(),
		})
		return err
	}
	fmt.Println("Result", result)
	c.Status(200).JSON(fiber.Map{
		"message": "Hand over logged",
		"result":  result,
	})
	return nil
}

func InEventHandler(c *fiber.Ctx) error {
	attendanceQuery := new(inEventModel.AttendanceQuery)
	if err := c.QueryParser(attendanceQuery); err != nil {
		fmt.Println("Error", err)
		c.Status(400).JSON(fiber.Map{
			"message": "Query Fields missing",
			"error":   err.Error(),
		})
		return err
	}
	userCollection, e := database.GetCollection("zeus_Events", "Users")
	if e != nil {
		fmt.Println("Error", e)
		c.Status(500).JSON(fiber.Map{
			"message": "Error while getting user details",
			"error":   e.Error(),
		})
		return e
	}
	email := attendanceQuery.Email
	slug := attendanceQuery.Slug
	fmt.Println("Email", email)
	fmt.Println("Slug", slug)
	var userData userModel.User
	err := userCollection.FindOne(context.Background(), bson.M{"email": email, "events": bson.M{"$in": []string{slug}}}).Decode(&userData)
	if err != nil {
		fmt.Println("Error", e)
		c.Status(500).JSON(fiber.Map{
			"message": "Error while fetching details of " + email,
			"error":   e.Error(),
		})
		return err
	}
	if &userData == nil {
		fmt.Println("No such user exists")
		c.Status(404).JSON(fiber.Map{
			"message": "User not found for this combination of email and slug",
		})
		return nil
	}
	inEventDataCollection, e := database.GetCollection("zeus_Events", "InEventData")
	if e != nil {
		fmt.Println("Error", e)
		c.Status(500).JSON(fiber.Map{
			"message": "Error in accessing in events data",
			"error":   e.Error(),
		})
		return e
	}
	switch action := c.Params("action"); action {
	case "attendance":
		return grantAttendance(c, inEventDataCollection, &userData, slug)
	case "food":
		return handOverFood(c, inEventDataCollection, &userData, slug)
	default:
		c.Status(404).JSON(fiber.Map{
			"message": "Invalid action chosen",
		})
		return nil
	}
}
