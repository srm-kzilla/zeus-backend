package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/srm-kzilla/events/api/events/controller"
)

func handleRoot(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"message": "Choose any route to continue"})
}

func SetupApp(app *fiber.App) {
	api := app.Group("/api")
	api.Get("/", handleRoot)
	api.Get("/events", controller.GetAllEvents)
	api.Post("/event", controller.CreateEvent)
	api.Post("/register", controller.RegisterForEvent)
}
