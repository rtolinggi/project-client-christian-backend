package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/rtolinggi/sales-api/controllers"
	"github.com/rtolinggi/sales-api/middleware"
)

func SetKaryawanRoute(app *fiber.App) {

	api := app.Group("/api/karyawan", middleware.ProtectedRoute)

	api.Post("upload", controllers.UploadAvatar)

}
