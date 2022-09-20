package config

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var JWT_ACCESS_TOKEN_SECRET = []byte(Config("JWT_ACCESS_TOKEN_SECRET"))
var JWT_REFRESH_TOKEN_SECRET = []byte(Config("JWT_REFRESH_TOKEN_SECRET"))

type JWTClaim struct {
	NamaPengguna string
	Role         string
	jwt.RegisteredClaims
}

func GenerateAccessToken(time time.Time, data *JWTClaim, secret []byte) (string, error) {
	claim := &JWTClaim{
		NamaPengguna: data.NamaPengguna,
		Role:         data.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "server-token",
			ExpiresAt: jwt.NewNumericDate(time),
		},
	}

	tokenAlgoritma := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	token, err := tokenAlgoritma.SignedString(secret)

	if err != nil {
		return "", err
	}
	return token, nil

}
