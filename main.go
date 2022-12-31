package main

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html"
	"github.com/joho/godotenv"
	api "github.com/srm-kzilla/events/api"
)

var startTime time.Time

func rootFunction(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"message": "Welcome to the Zeus API"})
}

func healthCheck(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"message": "OK", "uptime": time.Since(startTime).String()})
}

func setupRoutes(app *fiber.App) {
	app.Get("/", rootFunction)
	app.Get("/health", healthCheck)
	api.SetupApp(app)
}

func init() {
	startTime = time.Now()
	err := godotenv.Load()
	if err != nil {
		log.Panicln(err)
	}
}

func main() {
	log.Println("Server Starting!!!")

	engine := html.New("./templates", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Creating a logger middleware
	app.Use(logger.New())

	//setting up cors
	app.Use(cors.New())

	//setting up a rate limiter for max 100 requests/min per user
	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 60 * time.Second,
	}))
	// setting up api routes
	setupRoutes(app)

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
