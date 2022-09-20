package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rtolinggi/sales-api/database"
	"github.com/rtolinggi/sales-api/models"
)

// Get All User
func GetUsers(c *fiber.Ctx) error {
	db := database.DB.DB
	user := new([]models.User)

	if err := db.Find(&user).Error; err != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": false,
			"message": err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    user,
	})
}
