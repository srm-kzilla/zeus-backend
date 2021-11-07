package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	api "github.com/srm-kzilla/events/api"
)

func rootFunction(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"message": "Welcome to the Zeus API"})
}

func setupRoutes(app *fiber.App) {
	app.Get("/", rootFunction)
	api.SetupApp(app)
}

func main() {
	log.Println("Server Starting!!!")
	app := fiber.New()
	setupRoutes(app)
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	app.Listen(":" + port)
}
