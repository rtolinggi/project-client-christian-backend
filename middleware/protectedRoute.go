package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rtolinggi/sales-api/config"
)

func ProtectedRoute(ctx *fiber.Ctx) error {
	refreshToken := ctx.Cookies("refresh-token")
	claims := &config.JWTClaim{}

	token, err := jwt.ParseWithClaims(refreshToken, claims, func(t *jwt.Token) (interface{}, error) {
		return config.JWT_REFRESH_TOKEN_SECRET, nil
	})

	if err != nil {
		v, _ := err.(*jwt.ValidationError)
		fmt.Println(v)
		switch v.Errors {
		case jwt.ValidationErrorSignatureInvalid:
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Secret tidak valid",
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
