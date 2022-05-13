package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	api "github.com/srm-kzilla/events/api"
	S3 "github.com/srm-kzilla/events/utils/services/s3"
)

func rootFunction(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"message": "Welcome to the Zeus API"})
}

func setupRoutes(app *fiber.App) {
	app.Get("/", rootFunction)
	api.SetupApp(app)
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Panicln(err)
	}
}

func main() {
	log.Println("Server Starting!!!")
	app := fiber.New()

	// Creating a logger middleware
	app.Use(logger.New())

	//setting up cors
	app.Use(cors.New())

	// setting up api routes
	setupRoutes(app)

	fmt.Println(S3.ListBuckets())

	//Setting up Port Value
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf(`
	################################################
	🛡️  Server listening on port: %s 🛡️
	################################################
  `, port)

	app.Listen(":" + port)
}
