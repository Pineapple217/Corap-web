package handlers

import (
	"Corap-web/database"
	"Corap-web/models"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

// UserList returns a list of users
func UserList(c *fiber.Ctx) error {
	users := database.Get()

	return c.JSON(fiber.Map{
		"success": true,
		"users":   users,
	})
}

// UserCreate registers a user
func UserCreate(c *fiber.Ctx) error {
	user := &models.User{
		// Note: when writing to external database,
		// we can simply use - Name: c.FormValue("user")
		Name: utils.CopyString(c.FormValue("user")),
	}
	database.Insert(user)

	return c.JSON(fiber.Map{
		"success": true,
		"user":    user,
	})
}

// NotFound returns custom 404 page
func NotFound(c *fiber.Ctx) error {
	return c.Status(404).SendFile("./static/private/404.html")
}

func Home(c *fiber.Ctx) error {
	// Render index within layouts/main
	return c.Render("index", fiber.Map{
		"Title": "Hello, World!",
	}, "layouts/main")
}

func Devices(c *fiber.Ctx) error {
	return c.Render("devices", fiber.Map{
		"Devices": database.GetDivices(),
	}, "layouts/main")
}

func Jobs(c *fiber.Ctx) error {
	return c.Render("jobs", fiber.Map{
		"Jobs": database.GetSchedulerJobs(),
	}, "layouts/main")
}
