package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/rtolinggi/sales-api/controllers"
	"github.com/rtolinggi/sales-api/middleware"
)

func SetUserRoute(app *fiber.App) {

	api := app.Group("/api/user", logger.New(), middleware.ProtectedRoute)

	api.Get("/", controllers.GetUsers)

}
