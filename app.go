package main

import (
	"Corap-web/database"
	"Corap-web/handlers"

	"flag"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/gofiber/template/html/v2"
)

var (
	port   = flag.String("port", ":3000", "Port to listen on")
	listen = flag.String("listen", "127.0.0.1", "Where to listen, 0.0.0.0 is needed for docker")
	prod   = flag.Bool("prod", false, "Enable prefork in Production")
)

func main() {
	flag.Parse()

	log.SetLevel(log.LevelInfo)

	database.Connect()

	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Prefork: *prod, // go run app.go -prod
		Views:   engine,
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cache.New(cache.Config{
		KeyGenerator: func(c *fiber.Ctx) string {
			return utils.CopyString(c.OriginalURL())
		},
	}))

	app.Get("/", handlers.Home)

	app.Get("/devices/table", handlers.DevicesTable)
	app.Get("/devices/:deveui", handlers.Device)
	app.Get("/devices", handlers.Devices)

	app.Get("/jobs", handlers.Jobs)

	app.Get("/scrape", handlers.Scrape)

	app.Get("/devices/:deveui/plots/:plot_type", handlers.DevicePlots)

	// Setup static files
	app.Static("/", "./static/public")

	// Handle not founds
	app.Use(handlers.NotFound)

	// Listen on port 3000
	log.Fatal(app.Listen(*listen + *port)) // go run app.go -port=:3000
}
