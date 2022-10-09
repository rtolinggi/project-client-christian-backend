package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/rtolinggi/sales-api/controllers"
	// "github.com/rtolinggi/sales-api/middleware"
)

func SetUserRoute(app *fiber.App) {

	api := app.Group("/api/user", logger.New())

	api.Get("/", controllers.GetUsers)
	api.Get("/:id", controllers.GetUserId)
	api.Put("/:id", controllers.UpdateUserId)
	api.Delete("/:id", controllers.DeleteUserId)

}
