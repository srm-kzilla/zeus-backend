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
	protected := api.Use(authController.AuthenticateAdmin)
	protected.Get("/event", eventController.GetEventById)
	protected.Get("/event/:slug", eventController.GetEventBySlug)
	protected.Get("/events", eventController.GetAllEvents)
	protected.Get("/users", eventController.GetEventUsers)
	protected.Post("/event", eventController.CreateEvent)
	protected.Put("/event", eventController.UpdateEvent)
	protected.Post("/register", userController.RegisterForEvent)
	protected.Post("/rsvp", userController.RsvpForEvent)
	protected.Post("/event/close", eventController.CloseEvent)
	protected.Get("/inEvent/:action", inEventController.InEventHandler)
	protected.Post("/event/upload/cover", eventController.UploadEventCover)
	protected.Post("/event/speaker", eventController.AddSpeaker)
}
