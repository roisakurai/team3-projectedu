package middleware

import (
	"os"

	echojwt "github.com/labstack/echo-jwt/v4"
)

func JWTMiddleware() echojwt.Config {
	return echojwt.Config{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}
}
