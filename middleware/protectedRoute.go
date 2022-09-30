package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rtolinggi/sales-api/config"
)

func ProtectedRoute(ctx *fiber.Ctx) error {

	accessToken := strings.Split(ctx.Get("Authorization"), " ")
	dataCookie := ctx.Cookies("refresh_token")
	println(dataCookie)
	println("TESTING")
	if accessToken[0] == "" || accessToken[1] == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Otentikasi Gagal, Silahkan Login untuk mendapatkan Token",
		})
	}

	claims := &config.JWTClaim{}

	token, err := jwt.ParseWithClaims(accessToken[1], claims, func(t *jwt.Token) (interface{}, error) {
		return config.JWT_ACCESS_TOKEN_SECRET, nil
	})

	if err != nil {
		v, _ := err.(*jwt.ValidationError)
		switch v.Errors {
		case jwt.ValidationErrorSignatureInvalid:
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Token tidak valid",
			})
		case jwt.ValidationErrorExpired:
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Token sudah expired",
			})
		default:
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Server mengalamai gangguan",
				"error":   v,
			})
		}
	}

	if !token.Valid {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Token tidak valid",
		})
	}

	return ctx.Next()
}
