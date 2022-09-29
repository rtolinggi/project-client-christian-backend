package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/rtolinggi/sales-api/controllers"
	"github.com/rtolinggi/sales-api/middleware"
)

func SetKaryawanRoute(app *fiber.App) {

	api := app.Group("/api/karyawan", logger.New(), middleware.ProtectedRoute)

	api.Post("upload", controllers.UploadAvatar)

}
