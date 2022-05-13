package inEventController

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/srm-kzilla/events/database"
)

func AdmitUser(c *fiber.Ctx) error {
	userCollection, e := database.GetCollection("zeus_Events", "Users")
	if e != nil {
		fmt.Println("Error", e)

	}
}
