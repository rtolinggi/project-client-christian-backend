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
			"code":   400,
			"status": "Bad Request",
			"data":   nil,
			"errors": fiber.Map{
				"message": err.Error(),
			},
		})
	}

	errors := validation.ValidatorStruct(*Karyawan)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":   400,
			"status": "Bad Request",
			"data":   nil,
			"errors": errors,
		})
	}

	db := database.DB.DB
	if err := db.Create(&Karyawan).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":   500,
			"status": "Internal Server Error",
			"data":   nil,
			"errors": fiber.Map{
				"message": err.Error(),
			},
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"code":   201,
		"status": "Created",
		"data":   Karyawan,
		"errors": nil,
	})
}

func UploadAvatar(c *fiber.Ctx) error {

	file, err := c.FormFile("avatar")

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":   500,
			"status": "Internal Server Error",
			"data":   nil,
			"errors": fiber.Map{
				"message": err.Error(),
			},
		})
	}

	typeFile := file.Header.Get("Content-type")
	validate := typeFile == "image/jpeg" || typeFile == "image/png"
	fmt.Println(validate)

	if !validate {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":   400,
			"status": "Bad Request",
			"data":   nil,
			"errors": fiber.Map{
				"message": "Fily Tepe not Allowed, Please insert type Image/png",
			},
		})
	}

	c.SaveFile(file, fmt.Sprintf("./public/images/%s", file.Filename))
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":   200,
		"status": "OK",
		"data": fiber.Map{
			"message": file.Filename,
		},
		"errors": nil,
	})
}

func UpdateKaryawan() {

}

func DeleteKaryawan() {

}
