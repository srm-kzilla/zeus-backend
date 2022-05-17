package api

import (
	"github.com/gofiber/fiber/v2"
	eventController "github.com/srm-kzilla/events/api/events/controller"
	inEventController "github.com/srm-kzilla/events/api/inEvent/controller"
	userController "github.com/srm-kzilla/events/api/users/controller"
	authController "github.com/srm-kzilla/events/api/auth/controller"
)

func handleRoot(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"message": "Choose any route to continue"})
}

func SetupApp(app *fiber.App) {
	api := app.Group("/api")
	api.Get("/", handleRoot)
	api.Post("/admin/register", authController.RegisterAdmin)
	api.Post("/admin/login", authController.LoginAdmin)
	api.Post("/admin/refresh", authController.RefreshAdmin)
	api.Get("/event", eventController.GetEventById)
	api.Get("/event/:slug", eventController.GetEventBySlug)
	api.Get("/events", eventController.GetAllEvents)
	api.Get("/users", eventController.GetEventUsers)
	api.Post("/event", eventController.CreateEvent)
	api.Put("/event", eventController.UpdateEvent)
	api.Post("/register", userController.RegisterForEvent)
	api.Post("/rsvp", userController.RsvpForEvent)
	api.Post("/event/close", eventController.CloseEvent)
	api.Get("/inEvent/:action", inEventController.InEventHandler)
	api.Post("/event/upload/cover", eventController.UploadEventCover)
	api.Post("/event/speaker", eventController.AddSpeaker)
}
