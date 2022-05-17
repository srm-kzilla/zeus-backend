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
	protected := api.Use(authController.AuthenticateAdmin)
	protected.Post("/event", eventController.CreateEvent)
	protected.Put("/event", eventController.UpdateEvent)
	protected.Post("/register", userController.RegisterForEvent)
	protected.Post("/rsvp", userController.RsvpForEvent)
	protected.Post("/event/close", eventController.CloseEvent)
	protected.Get("/inEvent/:action", inEventController.InEventHandler)
	protected.Post("/upload", eventController.UploadEventCover)
	protected.Post("/event/speaker", eventController.AddSpeaker)
}
