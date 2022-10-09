package controllers

import (
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rtolinggi/sales-api/database"
	"github.com/rtolinggi/sales-api/models"
	"github.com/rtolinggi/sales-api/validation"
	"gorm.io/gorm"
)

type UserUpdate struct {
	ID           uint   `json:"id"`
	NamaPengguna string `json:"nama_pengguna"`
	KataSandi    string `json:"kata_sandi" `
	Role         string `json:"role"`
	UpdatedAt    time.Time
}

// Get All User
func GetUsers(ctx *fiber.Ctx) error {

	db := database.DB.DB
	user := new([]models.User)

	if err := db.Order("created_at desc").Find(&user).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":   500,
			"status": "Internal Server Error",
			"data":   nil,
			"errors": fiber.Map{
				"message": err.Error(),
			},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":   200,
		"status": "OK",
		"data":   user,
		"errors": nil,
	})
}

// Get ID User
func GetUserId(ctx *fiber.Ctx) error {

	db := database.DB.DB
	user := new([]models.User)
	id := ctx.Params("id")

	if err := db.First(&user, id).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"code":   404,
				"status": "Not Found",
				"data":   nil,
				"errors": fiber.Map{
					"message": err.Error(),
				},
			})
		default:
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":   500,
				"status": "Internal Server Error",
				"data":   nil,
				"errors": fiber.Map{
					"message": strings.Split(err.Error(), "asdad"),
				},
			})
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":   200,
		"status": "OK",
		"data":   user,
		"errors": nil,
	})
}

// Update User By ID
func UpdateUserId(ctx *fiber.Ctx) error {

	db := database.DB.DB
	user := new(models.User)

	id := ctx.Params("id")

	if err := ctx.BodyParser(user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":   400,
			"status": "Bad Request",
			"data":   nil,
			"errors": fiber.Map{
				"message": err.Error(),
			},
		})
	}

	idVar, err := strconv.Atoi(id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":   400,
			"status": "Bad Request",
			"data":   nil,
			"errors": "ID Must type Number",
		})
	}

	user.ID = uint(idVar)
	if err := validation.ValidatorStruct(user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":   400,
			"status": "Bad Request",
			"data":   nil,
			"errors": err,
		})
	}

	hash, err := HashPassword(user.KataSandi)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":   500,
			"status": "Internal Server Error",
			"data":   nil,
			"errors": fiber.Map{
				"message": err.Error(),
			},
		})
	}

	user.KataSandi = hash
	count := db.Where("id = ?", id).Updates(&user)
	if count.RowsAffected == 0 {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"code":   404,
			"status": "Status Not Found",
			"data":   nil,
			"errors": "ID User Not Found",
		})
	}

	if count.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":   500,
			"status": "Internal Server Error",
			"data":   nil,
			"errors": fiber.Map{
				"message": count.Error.Error(),
			},
		})
	}

	data := []UserUpdate{
		{
			ID:           user.ID,
			NamaPengguna: user.NamaPengguna,
			KataSandi:    user.KataSandi,
			Role:         user.Role,
			UpdatedAt:    user.UpdatedAt,
		},
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":   200,
		"status": "Status OK",
		"data":   data,
		"errors": nil,
	})
}

// Delete User By ID
func DeleteUserId(ctx *fiber.Ctx) error {

	db := database.DB.DB
	user := new(models.User)

	id := ctx.Params("id")
	idVar, err := strconv.Atoi(id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":   400,
			"status": "Bad Request",
			"data":   nil,
			"errors": "ID Must type Number",
		})
	}

	user.ID = uint(idVar)
	if err := db.Unscoped().Delete(&user).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":   400,
			"status": "Bad Request",
			"data":   nil,
			"errors": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":   200,
		"status": "Status OK",
		"data":   user,
		"errors": nil,
	})

}
