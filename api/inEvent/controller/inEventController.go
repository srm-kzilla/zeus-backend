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
	"github.com/srm-kzilla/events/validators"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)


func hasAttendance(inEventDataCollection *mongo.Collection, userData UserModel.User, slug string)(bool, error) {
	count, e := inEventDataCollection.CountDocuments(context.Background(), bson.M{"userId": userData.ID, "eventSlug": slug})
	if e != nil {
		fmt.Println("Error", e)
		return false, e
	}
	fmt.Println(count)
	if count != int64(0){
		return true, nil
	}
	return false, nil
}

func grantAttendance(c *fiber.Ctx, inEventDataCollection *mongo.Collection, userData UserModel.User, slug string) error {
	// count, e := inEventDataCollection.CountDocuments(context.Background(), bson.M{"userId": userData.ID, "eventSlug": slug})
	// if e != nil {
	// 	fmt.Println("Error", e)
	// 	return c.Status(500).JSON(fiber.Map{
	// 		"message": "Error while checking pre existing attendance",
	// 		"error":   e.Error(),
	// 	})
	// }
	// fmt.Println(count)
	// if count != int64(0){
	// 	return c.Status(409).JSON(fiber.Map{
	// 		"message": "Attendance already granted",
	// 	})
	// }
	attendance, err := hasAttendance(inEventDataCollection, userData, slug)
	if err != nil {
		fmt.Println("Error", err)
		return c.Status(500).JSON(fiber.Map{
			"message": "Error while checking pre existing attendance",
			"error":   err.Error(),
		})
	}
	if attendance {
		return c.Status(409).JSON(fiber.Map{
			"message": "Attendance already granted",
		})
	}
	inEventData := InEventModel.InEventData{
		ID: primitive.NewObjectID(),
		UserID: userData.ID,
		EventSlug: slug,
		FoodReceived: false,
	}
	res, err := inEventDataCollection.InsertOne(context.Background(), inEventData)
	if err != nil {
		fmt.Println("Error", err)
		return c.Status(500).JSON(fiber.Map{
			"message": "Error while granting attendance",
			"error":   err.Error(),
		})
	}
	fmt.Println("Attendance granted to : ", res.InsertedID)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Attendance granted",
	})
}

func handOverFood(c *fiber.Ctx, inEventDataCollection *mongo.Collection, userData UserModel.User, slug string) error {
	attendance, err := hasAttendance(inEventDataCollection, userData, slug)
	if err != nil {
		fmt.Println("Error", err)
		return c.Status(500).JSON(fiber.Map{
			"message": "Error while checking pre existing attendance",
			"error":   err.Error(),
		})
	}
	if !attendance {
		return c.Status(409).JSON(fiber.Map{
			"message": "missing attendance",
		})
	}
	count, e := inEventDataCollection.CountDocuments(context.Background(), bson.M{"userId": userData.ID, "eventSlug": slug, "foodReceived": true})
	if e != nil {
		fmt.Println("Error", e)
		return c.Status(500).JSON(fiber.Map{
			"message": "Error while checking pre existing food",
			"error":   e.Error(),
		})
	}
	fmt.Println(count)
	if count != int64(0){
		return c.Status(409).JSON(fiber.Map{
			"message": "Food already handed over",
		})
	}

	result, err := inEventDataCollection.UpdateOne(context.Background(), bson.M{"userId": userData.ID, "eventSlug":slug}, bson.M{"$set": bson.M{"foodReceived": true}})
	if err != nil {
		fmt.Println("Error", err)
		return c.Status(500).JSON(fiber.Map{
			"message": "Error while logging handover",
		})
	}
	fmt.Println("Food handed over to : ", result.UpsertedID)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Food handed over",
	})
}

func InEventHandler(c *fiber.Ctx)error {
	// attendanceQuery := new(InEventModel.AttendanceQuery)
	// if err := c.QueryParser(attendanceQuery); err != nil {
	// 	fmt.Println("Error", err)
	// 	return c.Status(400).JSON(fiber.Map{
	// 		"message": "Query Fields missing",
	// 		"error":   err.Error(),
	// 	})
	// }
	var attendanceQuery InEventModel.AttendanceQuery
	c.BodyParser(&attendanceQuery)

	E := validators.ValidateAttendanceQuery(attendanceQuery)
	if E != nil {
		c.Status(fiber.StatusBadRequest).JSON(E)
		return nil
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
	inEventsCollection, err := database.GetCollection("zeus_Events", "InEvents")
	if err != nil {
		fmt.Println("Error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error while getting inEvents collection",
			"error":   err.Error(),
		})
	}
	var eventData EventModel.Event
	e := eventsCollection.FindOne(context.Background(), bson.M{"slug": attendanceQuery.Slug}).Decode(&eventData)
	if e != nil {
		fmt.Println("Error", e)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "no such event/slug exists",
		})
	}
	var userData UserModel.User
	objId, _ := primitive.ObjectIDFromHex(attendanceQuery.UserID)
	error := userCollection.FindOne(context.Background(), bson.M{"_id":objId, "events": bson.M{"$in": []string{attendanceQuery.Slug}}}).Decode(&userData)
	if error != nil {
		fmt.Println("Error", error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "user not found or user not registered for this event",
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
		grantAttendance(c, inEventsCollection, userData, eventData.Slug)
	case "food":
		// run the food code
		handOverFood(c, inEventsCollection, userData, eventData.Slug)
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid action chosen",
		})
	}
return nil
}

func GetInEventData (c *fiber.Ctx)error{
	var slug = c.Query("slug")
	slug = strings.ToLower(slug)
	if slug == ""{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "slug not provided",
		})
	}
	var eventData []InEventModel.InEventData
	inEventDataCollection, err := database.GetCollection("zeus_Events", "InEvents")
	if err != nil {
		fmt.Println("Error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error while getting inEvents collection",
			"error":   err.Error(),
		})
	}
	cursor, e := inEventDataCollection.Find(context.Background(), bson.M{"eventSlug": slug})
	if e != nil {
		fmt.Println("Error", e)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error while getting inEvents collection",
		})
	}
	if e = cursor.All(context.Background(), &eventData); e != nil {
		fmt.Println("Error ", err)
		c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":  err.Error(),
		})
	}
	var foodHandedOver []InEventModel.InEventData
	for _, v := range eventData {
		if v.FoodReceived {
			foodHandedOver = append(foodHandedOver, v)
		}
	}
	var numOfAttendees = len(eventData)
	var numOfFoodHandedOver = len(foodHandedOver)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"numOfAttendees": numOfAttendees,
		"numOfFoodHandedOver": numOfFoodHandedOver,
	})
		
}