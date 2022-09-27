package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/rtolinggi/sales-api/controllers"
)

func SetAuthRoute(app *fiber.App) {

	api := app.Group("/api/auth", logger.New())

	api.Post("/signin", controllers.SignIn)
	api.Post("/signup", controllers.SignUp)
	api.Get("/signout", controllers.SignOut)

}
