package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rtolinggi/sales-api/database"
	"github.com/rtolinggi/sales-api/models"
	"github.com/rtolinggi/sales-api/validation"
)

func GetKaryawan() {

}

func CreateKaryawan(c *fiber.Ctx) error {
	Karyawan := new(models.Karyawan)

	if err := c.BodyParser(Karyawan); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"errors":  err.Error(),
		})
	}

	errors := validation.ValidatorStruct(*Karyawan)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"errors":  errors,
		})
	}

	db := database.DB.DB
	if err := db.Create(&Karyawan).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Couldn't create Karyawan",
			"data":    err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    Karyawan,
	})
}

func UpdateKaryawan() {

}

func DeleteKaryawan() {

}
