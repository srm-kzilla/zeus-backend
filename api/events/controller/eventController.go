package controller

import (
	"context"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	model "github.com/srm-kzilla/events/api/events/model"
	"github.com/srm-kzilla/events/database"
	"go.mongodb.org/mongo-driver/bson"
)

func GetAllEvents(c *fiber.Ctx) error {
	var events []model.Event
	// var c *fiber.Ctx

	eventsCollection, e := database.GetCollection("zeus_Events", "Events")
	if e != nil {
		fmt.Println("Error: ", e)
		c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": e.Error(),
		})
	}

	cursor, err := eventsCollection.Find(context.Background(), bson.D{})
	if err = cursor.All(context.Background(), &events); err != nil {
		log.Fatal("Error ", err)
		c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":  err.Error(),
			"events": events,
		})
	}

	c.Status(fiber.StatusOK).JSON(events)

	return nil
}

func CreateEvent(c *fiber.Ctx) error {
	var event model.Event

	c.BodyParser(&event)

	eventsCollection, e := database.GetCollection("zeus_Events", "Events")
	if e != nil {
		fmt.Println("Error: ", e)
		c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": e.Error(),
		})
	}

	res, err := eventsCollection.InsertOne(context.Background(), event)
	if err != nil {
		log.Fatal("Error", err)
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Success":     true,
			"Inserted ID": res.InsertedID,
		})

	}

	c.Status(fiber.StatusCreated).JSON(event)

	return nil
}
