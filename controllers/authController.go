package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rtolinggi/sales-api/config"
	"github.com/rtolinggi/sales-api/database"
	"github.com/rtolinggi/sales-api/models"
	"github.com/rtolinggi/sales-api/validation"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type InputSignIn struct {
	NamaPengguna string `json:"nama_pengguna" validate:"required"`
	KataSandi    string `json:"kata_sandi" validate:"required"`
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func decodePassword(hashPassword, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password)); err != nil {
		return err
	}
	return nil
}

func SignIn(ctx *fiber.Ctx) error {
	input := new(InputSignIn)
	user := new(models.User)

	if err := ctx.BodyParser(input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"errors":  err.Error(),
			"message": "Body Parser Error",
		})
	}

	if err := validation.ValidatorStruct(input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"errors":  err,
		})
	}

	db := database.DB.DB
	if err := db.Where("nama_pengguna = ?", input.NamaPengguna).First(&user).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Nama Pengguna atau Password salah",
			})
		default:
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Terjadi gangguan di Server",
			})
		}
	}

	if err := decodePassword(user.KataSandi, input.KataSandi); err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Nama Pengguna atau Password Salah",
		})
	}

	accessToken, err := config.GenerateAccessToken(time.Now().Add(time.Minute*1), &config.JWTClaim{
		NamaPengguna: user.NamaPengguna,
		Role:         user.Role,
	}, config.JWT_ACCESS_TOKEN_SECRET)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Terjadi gangguan di Server, gagal membuat token",
			"error":   err.Error(),
		})
	}

	refreshToken, err := config.GenerateAccessToken(time.Now().Add(time.Hour*24*30), &config.JWTClaim{
		NamaPengguna: user.NamaPengguna,
		Role:         user.Role,
	}, config.JWT_REFRESH_TOKEN_SECRET)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Terjadi gangguan di Server, gagal membuat token",
			"error":   err.Error(),
		})
	}

	user.RefreshToken = refreshToken
	if err := db.Where("nama_pengguna = ?", input.NamaPengguna).Updates(&user).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "terjadi gangguan di server",
			"error":   err,
		})
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh-token",
		Path:     "/",
		Value:    refreshToken,
		Expires:  time.Now().Add(time.Hour * 24 * 30),
		HTTPOnly: true,
		SameSite: "lax",
	})

	return ctx.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"success":      true,
		"message":      "Login Berhasil",
		"access_token": accessToken,
		"user": fiber.Map{
			"nama_pengguna": user.NamaPengguna,
			"role":          user.Role,
		},
	})
}

func SignUp(ctx *fiber.Ctx) error {
	user := new(models.User)

	if err := ctx.BodyParser(user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"errors":  err.Error(),
		})
	}

	if err := validation.ValidatorStruct(*user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"errors":  err,
		})
	}

	hash, err := hashPassword(user.KataSandi)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Couldn't hash password",
			"data":    err.Error(),
		})
	}

	user.KataSandi = hash
	db := database.DB.DB
	if err := db.Create(&user).Error; err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Couldn't create user",
			"data":    err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    user,
	})
}

func SignOut(ctx *fiber.Ctx) error {
	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh-token",
		Expires:  time.Now().Add(-(time.Hour * 2)),
		HTTPOnly: true,
		SameSite: "lax",
	})
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Logout Is Success",
	})
}
