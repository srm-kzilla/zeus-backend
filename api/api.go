package api

import "github.com/gofiber/fiber/v2"

func handleRoot(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"message": "Choose any route to continue"})
}
func SetupApp(app *fiber.App) {
	api := app.Group("/api")
	api.Get("/", handleRoot)
}
