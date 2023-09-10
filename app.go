package main

import (
	"Corap-web/database"
	"Corap-web/handlers"

	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
)

var (
	port = flag.String("port", ":3000", "Port to listen on")
	prod = flag.Bool("prod", false, "Enable prefork in Production")
)

func main() {
	// Parse command-line flags
	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Connected with database
	database.Connect()

	engine := html.New("./views", ".html")

	// Create fiber app
	app := fiber.New(fiber.Config{
		Prefork: *prod, // go run app.go -prod
		Views:   engine,
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())

	app.Get("/", handlers.Home)
	app.Get("/devices", handlers.Devices)
	app.Get("/jobs", handlers.Jobs)
	// Create a /api/v1 endpoint
	v1 := app.Group("/api/v1")

	// Bind handlers
	v1.Get("/users", handlers.UserList)
	v1.Post("/users", handlers.UserCreate)

	// Setup static files
	app.Static("/", "./static/public")

	// Handle not founds
	app.Use(handlers.NotFound)

	// Listen on port 3000
	log.Fatal(app.Listen("127.0.0.1" + *port)) // go run app.go -port=:3000
}
