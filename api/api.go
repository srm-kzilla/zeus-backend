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
	api.Get("/event", controller.GetEventById)
	api.Get("/event/:slug", controller.GetEventBySlug)
	api.Get("/events", controller.GetAllEvents)
	api.Get("/users", controller.GetEventUsers)
	api.Post("/event", controller.CreateEvent)
	api.Post("/register", controller.RegisterForEvent)
}
