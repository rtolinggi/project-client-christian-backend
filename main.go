package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/rtolinggi/sales-api/database"
	"github.com/rtolinggi/sales-api/routes"
)

func main() {

	// Inital Database to connect
	database.ConnectDB()

	// initial App to running web server
	app := fiber.New()

	app.Use(cors.New())

	// inital Route
	routes.SetUserRoute(app)
	routes.SetAuthRoute(app)

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "404 Not Found",
		})
	})

	log.Fatal(app.Listen(":8080"))
}
