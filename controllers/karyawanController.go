package controllers

import (
	"fmt"

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

func UploadAvatar(c *fiber.Ctx) error {

	file, err := c.FormFile("avatar")

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Error to Upload File",
			"data":    err.Error(),
		})
	}

	typeFile := file.Header.Get("Content-type")
	validate := typeFile == "image/jpeg" || typeFile == "image/png"
	fmt.Println(validate)

	if !validate {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Fily Tepe not Allowed, Please insert type Image/png",
		})
	}

	c.SaveFile(file, fmt.Sprintf("./public/images/%s", file.Filename))
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Image Saved",
	})
}

func UpdateKaryawan() {

}

func DeleteKaryawan() {

}
