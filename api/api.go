package api

import (
	"github.com/gofiber/fiber/v2"
	authController "github.com/srm-kzilla/events/api/auth/controller"
	eventController "github.com/srm-kzilla/events/api/events/controller"
	inEventController "github.com/srm-kzilla/events/api/inEvent/controller"
	userController "github.com/srm-kzilla/events/api/users/controller"
)

func handleRoot(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"message": "Choose any route to continue"})
}

func SetupApp(app *fiber.App) {
	api := app.Group("/api")
	api.Get("/", handleRoot)
	//api.Post("/admin/register", authController.RegisterAdmin)
	api.Post("/admin/login", authController.LoginAdmin)
	api.Post("/admin/refresh", authController.RefreshAdmin)
	api.Get("/event", eventController.GetEventById)
	api.Get("/event/:slug", eventController.GetEventBySlug)
	api.Get("/events", eventController.GetAllEvents)
	api.Post("/register", userController.RegisterForEvent)
	api.Get("/rsvp", userController.RsvpForEvent)
	api.Get("/event/:slug/registrations", eventController.GetEventRegistrationsCount)
	protected := api.Use(authController.AuthenticateAdmin)
	protected.Get("/users", eventController.GetEventUsers)
	protected.Get("/user/:userId", userController.GetUserById)
	protected.Post("/event", eventController.CreateEvent)
	protected.Put("/event", eventController.UpdateEvent)
	protected.Post("/event/close", eventController.CloseEvent)
	protected.Post("/event/registration/close", eventController.CloseRegistrations)
	protected.Post("/inevent", inEventController.InEventHandler)
	protected.Get("/inevent/data", inEventController.GetInEventData)
	protected.Post("/upload", eventController.UploadEventCover)
	protected.Post("/event/speaker", eventController.AddSpeaker)
	protected.Put("/event/speaker", eventController.UpdateSpeaker)
}
