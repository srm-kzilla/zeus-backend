package api

import (
	"github.com/gofiber/fiber/v2"
	eventController "github.com/srm-kzilla/events/api/events/controller"
	userController "github.com/srm-kzilla/events/api/users/controller"
)

func handleRoot(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"message": "Choose any route to continue"})
}

func SetupApp(app *fiber.App) {
	api := app.Group("/api")
	api.Get("/", handleRoot)
	api.Get("/event", eventController.GetEventById)
	api.Get("/event/:slug", eventController.GetEventBySlug)
	api.Get("/events", eventController.GetAllEvents)
	api.Get("/users", eventController.GetEventUsers)
	api.Post("/event", eventController.CreateEvent)
	api.Post("/register", userController.RegisterForEvent)
	api.Post("/rsvp", userController.RsvpForEvent)
	api.Post("/event/close", eventController.CloseEvent)
}
