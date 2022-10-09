package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
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

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func DecodePassword(hashPassword, password string) error {
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
			"code":   400,
			"status": "Bad Request",
			"data":   nil,
			"errors": fiber.Map{
				"message": err.Error(),
			},
		})
	}

	if err := validation.ValidatorStruct(input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":   400,
			"status": "Bad Request",
			"data":   nil,
			"errors": err,
		})
	}

	db := database.DB.DB
	if err := db.Where("nama_pengguna = ?", input.NamaPengguna).First(&user).Error; err != nil {
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
					"message": err.Error(),
				},
			})
		}
	}

	if err := DecodePassword(user.KataSandi, input.KataSandi); err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"code":   401,
			"status": "Unauthorized",
			"data":   nil,
			"errors": fiber.Map{
				"message": "Salah Decode Passwod",
			},
		})
	}

	accessToken, err := config.GenerateAccessToken(time.Now().Add(time.Minute*1), &config.JWTClaim{
		NamaPengguna: user.NamaPengguna,
		Role:         user.Role,
	}, config.JWT_ACCESS_TOKEN_SECRET)

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

	refreshToken, err := config.GenerateAccessToken(time.Now().Add(time.Hour*24*30), &config.JWTClaim{
		NamaPengguna: user.NamaPengguna,
		Role:         user.Role,
	}, config.JWT_REFRESH_TOKEN_SECRET)

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

	user.RefreshToken = refreshToken
	if err := db.Where("nama_pengguna = ?", input.NamaPengguna).Updates(&user).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":   500,
			"status": "Internal Server Error",
			"data":   nil,
			"errors": fiber.Map{
				"message": err.Error(),
			},
		})
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Path:     "/",
		Value:    refreshToken,
		Expires:  time.Now().Add(time.Hour * 24 * 30),
		HTTPOnly: true,
		SameSite: "lax",
	})

	return ctx.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"code":   202,
		"status": "Accepted",
		"data": fiber.Map{
			"access_token": accessToken,
			"user": fiber.Map{
				"nama_pengguna": user.NamaPengguna,
				"role":          user.Role,
			},
		},
		"errors": nil,
	})
}

func SignUp(ctx *fiber.Ctx) error {
	user := new(models.User)

	db := database.DB.DB

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

	if err := validation.ValidatorStruct(user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":   400,
			"status": "Bad Request",
			"data":   nil,
			"errors": err,
		})
	}

	count := db.Limit(1).Find(&user, "nama_pengguna = ?", user.NamaPengguna)
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

	if count.RowsAffected > 0 {
		return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
			"code":   409,
			"status": "Conflict",
			"data":   nil,
			"errors": fiber.Map{
				"message": "Nama Pengguna Already Exist",
			},
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

	if err := db.Create(&user).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":   500,
			"status": "Internal Server Error",
			"data":   nil,
			"errors": fiber.Map{
				"message": err.Error(),
			},
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"code":   201,
		"status": "Created",
		"data":   user,
		"errors": nil,
	})

}

func SignOut(ctx *fiber.Ctx) error {
	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Path:     "/",
		Expires:  time.Now().Add(-(time.Hour * 2)),
		HTTPOnly: true,
		SameSite: "lax",
	})

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":   200,
		"status": "OK",
		"data": fiber.Map{
			"message": "Logout Success",
		},
		"errors": nil,
	})
}

func GetToken(ctx *fiber.Ctx) error {
	refreshToken := ctx.Cookies("refresh_token")

	if refreshToken == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"code":   401,
			"status": "Unauthorized",
			"data":   nil,
			"errors": fiber.Map{
				"message": "Token Empty, Please login",
			},
		})
	}

	db := database.DB.DB
	user := new(models.User)

	if err := db.Where("refresh_token = ?", refreshToken).First(&user).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"code":   401,
				"status": "Unauthorized",
				"data":   nil,
				"errors": fiber.Map{
					"message": "Token not Valid and not found in Database",
				},
			})
		default:
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":   500,
				"status": "Internal Server Error",
				"data":   nil,
				"errors": fiber.Map{
					"message": "Terjadi gangguan di Server",
				},
			})
		}
	}

	claims := &config.JWTClaim{}

	token, err := jwt.ParseWithClaims(refreshToken, claims, func(t *jwt.Token) (interface{}, error) {
		return config.JWT_REFRESH_TOKEN_SECRET, nil
	})

	if err != nil {
		v, _ := err.(*jwt.ValidationError)
		switch v.Errors {
		case jwt.ValidationErrorSignatureInvalid:
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"code":   401,
				"status": "Unauthorized",
				"data":   nil,
				"errors": fiber.Map{
					"message": "Token not valid",
				},
			})
		case jwt.ValidationErrorExpired:
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"code":   401,
				"status": "Unaautorized",
				"data":   nil,
				"errors": fiber.Map{
					"message": "Token is expired",
				},
			})
		default:
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":   500,
				"status": "Internal Server Error",
				"data":   nil,
				"errors": fiber.Map{
					"message": v.Errors,
				},
			})
		}
	}

	if !token.Valid {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"code":    401,
			"status":  "Unauthorized",
			"data":    nil,
			"success": false,
			"errors": fiber.Map{
				"message": "Token Not valid",
			},
		})
	}

	accessToken, err := config.GenerateAccessToken(time.Now().Add(time.Minute*1), &config.JWTClaim{
		NamaPengguna: claims.NamaPengguna,
		Role:         claims.Role,
	}, config.JWT_ACCESS_TOKEN_SECRET)

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

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":   200,
		"status": "OK",
		"data": fiber.Map{
			"access_token": accessToken,
		},
		"errors": nil,
	})
}
